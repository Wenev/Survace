package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"sync"

	"github.com/Wenev/Survace/proto/gen/controller"
	"github.com/Wenev/Survace/proto/gen/dto"
	"github.com/Wenev/Survace/stream-service/ports/in"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	grpcPort          int
	server            *grpc.Server
	service           in.StreamService
	streamChatService in.StreamChatService
	controller.UnimplementedStreamServiceServer
}

func NewGrpcServer(grpcPort int, service in.StreamService, streamChatService in.StreamChatService) *GrpcServer {
	return &GrpcServer{
		grpcPort:          grpcPort,
		server:            grpc.NewServer(),
		service:           service,
		streamChatService: streamChatService,
	}
}

func (g *GrpcServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", g.grpcPort))
	if err != nil {
		return err
	}
	controller.RegisterStreamServiceServer(g.server, g)
	return g.server.Serve(listener)
}

func (g *GrpcServer) Stop() {
	g.server.GracefulStop()
}

func (g *GrpcServer) GetStreamToken(ctx context.Context, req *dto.GetStreamTokenRequest) (*dto.GetStreamTokenResponse, error) {
	fmt.Print("Hii....")
	code, msg, token, callId, err := g.service.GetStreamToken(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &dto.GetStreamTokenResponse{
		Code:    code,
		Message: msg,
		Token:   token,
		CallId:  callId,
	}, nil
}

func (g *GrpcServer) GetLiveStreamId(ctx context.Context, req *dto.GetLiveStreamIdRequest) (*dto.GetLiveStreamIdResponse, error) {
	code, msg, callId, isLive, err := g.service.GetLiveStreamId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &dto.GetLiveStreamIdResponse{
		Code:     code,
		Message:  msg,
		StreamId: callId,
		IsLive:   isLive,
	}, nil
}

func (g *GrpcServer) IsUserLive(ctx context.Context, req *dto.IsUserLiveRequest) (*dto.IsUserLiveResponse, error) {
	isLive, err := g.service.IsUserLive(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &dto.IsUserLiveResponse{
		IsLive: isLive,
	}, nil
}

func (g *GrpcServer) GoLive(ctx context.Context, req *dto.GoLiveRequest) (*dto.GoLiveResponse, error) {
	code, msg, err := g.service.GoLive(ctx, req.UserId, req.CallId)
	if err != nil {
		return nil, err
	}
	return &dto.GoLiveResponse{
		Code:    code,
		Message: msg,
	}, nil
}

func (g *GrpcServer) StopLive(ctx context.Context, req *dto.StopLiveRequest) (*dto.StopLiveResponse, error) {
	code, msg, err := g.service.StopLive(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &dto.StopLiveResponse{
		Code:    code,
		Message: msg,
	}, nil
}

func (g *GrpcServer) ListLiveStreams(ctx context.Context, req *dto.ListLiveStreamsRequest) (*dto.
ListLiveStreamsResponse, error) {
	streams, err := g.service.ListLiveStreams(ctx)
	if err != nil {
		return &dto.ListLiveStreamsResponse{
			Code:    13,
			Message: "Failed to fetch live streams",
			Streams: nil,
		}, err
	}
	var respStreams []*dto.LiveStreamInfo
	for _, s := range streams {
		respStreams = append(respStreams, &dto.LiveStreamInfo{
			UserId: s.UserId,
			CallId: s.CallId,
		})
	}
	return &dto.ListLiveStreamsResponse{
		Code:    0,
		Message: "success",
		Streams: respStreams,
	}, nil
}

func (g *GrpcServer) StreamChat(stream controller.StreamService_StreamChatServer) error {
	ctx := stream.Context()
	fmt.Print("HFDHJSKJDSHFJ")
	firstMsg, err := stream.Recv()
	if err != nil {
		return err
	}
	callId := fmt.Sprintf("%d", firstMsg.CallId)
	senderId := firstMsg.SenderId

	room := getOrCreateChatRoom(callId)

	client := &streamClient{
		id: senderId,
		ch: make(chan *dto.StreamChatMessage, 10),
	}
	room.register <- client
	defer func() { room.unregister <- client }()

	_, _, savedMsg, err := g.streamChatService.SendMessage(ctx, callId, senderId, firstMsg.Content)
	if err != nil {
		_ = stream.Send(&dto.StreamChatMessage{
			SenderId:  senderId,
			CallId:    firstMsg.CallId,
			Content:   "[ERROR] Failed to save message",
			CreatedAt: firstMsg.CreatedAt,
		})
	} else {
		room.broadcast <- &dto.StreamChatMessage{
			SenderId:  savedMsg.SenderID,
			CallId:    firstMsg.CallId,
			Content:   savedMsg.Content,
			CreatedAt: timestamppb.New(savedMsg.CreatedAt),
		}
	}
	done := make(chan struct{})
	go func() {
		for {
			select {
			case msg, ok := <-client.ch:
				if !ok {
					return
				}
				if err := stream.Send(msg); err != nil {
					return
				}
			case <-done:
				return
			}
		}
	}()

	for {
		msg, err := stream.Recv()
		if err != nil {
			close(done)
			return err
		}

		_, _, savedMsg, err := g.streamChatService.SendMessage(ctx, callId, senderId, msg.Content)
		if err != nil {
			_ = stream.Send(&dto.StreamChatMessage{
				SenderId:  senderId,
				CallId:    msg.CallId,
				Content:   "[ERROR] Failed to save message",
				CreatedAt: msg.CreatedAt,
			})
		} else {
			room.broadcast <- &dto.StreamChatMessage{
				SenderId:  savedMsg.SenderID,
				CallId:    msg.CallId,
				Content:   savedMsg.Content,
				CreatedAt: timestamppb.New(savedMsg.CreatedAt),
			}
		}
	}
}

type streamClient struct {
	id int32
	ch chan *dto.StreamChatMessage
}

type streamRoom struct {
	clients    map[int32]*streamClient
	register   chan *streamClient
	unregister chan *streamClient
	broadcast  chan *dto.StreamChatMessage
}

var (
	streamRooms   = make(map[string]*streamRoom)
	streamRoomsMu sync.Mutex
)

func getOrCreateChatRoom(callId string) *streamRoom {
	streamRoomsMu.Lock()
	defer streamRoomsMu.Unlock()

	if room, exists := streamRooms[callId]; exists {
		return room
	}

	room := &streamRoom{
		clients:    make(map[int32]*streamClient),
		register:   make(chan *streamClient),
		unregister: make(chan *streamClient),
		broadcast:  make(chan *dto.StreamChatMessage, 32),
	}

	go func() {
		for {
			select {
			case client := <-room.register:
				room.clients[client.id] = client
			case client := <-room.unregister:
				if _, ok := room.clients[client.id]; ok {
					delete(room.clients, client.id)
					close(client.ch)
				}
			case msg := <-room.broadcast:
				for _, client := range room.clients {
					select {
					case client.ch <- msg:
					default:
						delete(room.clients, client.id)
						close(client.ch)
					}
				}
			}
		}
	}()

	streamRooms[callId] = room
	return room
}
