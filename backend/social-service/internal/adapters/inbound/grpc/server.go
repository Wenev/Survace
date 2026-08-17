package grpc

import (
	"context"
	"fmt"
	"github.com/Acad600-TPA/WEB-WE-251/backend/social-service/ports/in"
	"github.com/Acad600-TPA/WEB-WE-251/proto/gen/controller"
	"github.com/Acad600-TPA/WEB-WE-251/proto/gen/dto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"sync"
)

type GrpcServer struct {
	grpcPort           int
	server             *grpc.Server
	chatService        in.ChatService
	followService      in.FollowService
	notificationClient *GalactusClient
	controller.SocialServiceServer
}

func NewGrpcServer(grpcPort int, chatService in.ChatService, followService in.FollowService, notificationClient *GalactusClient) *GrpcServer {
	return &GrpcServer{
		grpcPort:           grpcPort,
		server:             grpc.NewServer(),
		chatService:        chatService,
		followService:      followService,
		notificationClient: notificationClient,
	}
}

func (g *GrpcServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", g.grpcPort))
	if err != nil {
		fmt.Print("HELP")
		return err
	}
	controller.RegisterSocialServiceServer(g.server, g)
	err = g.server.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

func (g *GrpcServer) Stop() {
	g.server.GracefulStop()
}

// CalculateRoomId generates a deterministic room ID from two user IDs
func calculateRoomId(userId1, userId2 int32) int64 {
	var smallerId, largerId int32
	if userId1 < userId2 {
		smallerId = userId1
		largerId = userId2
	} else {
		smallerId = userId2
		largerId = userId1
	}
	return int64(smallerId)*100000 + int64(largerId)
}

func (g *GrpcServer) SendMessage(ctx context.Context, in *dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
	code, message, msg, err := g.chatService.SendMessage(ctx, in.SenderId, in.ReceiverId, in.Content)
	if err != nil {
		return &dto.SendMessageResponse{
			Code:    code,
			Message: message,
		}, nil
	}

	// Calculate room ID internally (not exposed to client)
	_ = calculateRoomId(in.SenderId, in.ReceiverId)

	// Send notification to the receiver after successful message send
	if g.notificationClient != nil {
		_, err = g.notificationClient.NotificationClient.SendNotification(ctx, &dto.SendNotificationRequest{
			UserId:  in.ReceiverId,
			Type:    "message",
			Content: fmt.Sprintf("You have a new message from user %d", in.SenderId),
		})
	}

	return &dto.SendMessageResponse{
		Code:    code,
		Message: message,
		Data: &dto.ChatMessage{
			SenderId:   msg.SenderID,
			ReceiverId: msg.ReceiverID,
			Content:    msg.Content,
			CreatedAt:  timestamppb.New(msg.CreatedAt),
		},
	}, nil
}

func (g *GrpcServer) ListMessages(ctx context.Context, in *dto.ListMessagesRequest) (*dto.ListMessagesResponse, error) {
	code, message, msgs, err := g.chatService.GetChatWithUser(ctx, in.UserId, in.OtherUserId)
	if err != nil {
		return &dto.ListMessagesResponse{
			Code:    code,
			Message: message,
		}, nil
	}

	// Calculate room ID internally (not exposed to client)
	_ = calculateRoomId(in.UserId, in.OtherUserId)

	var data []*dto.ChatMessage
	for _, m := range msgs {
		data = append(data, &dto.ChatMessage{
			SenderId:   m.SenderID,
			ReceiverId: m.ReceiverID,
			Content:    m.Content,
			CreatedAt:  timestamppb.New(m.CreatedAt),
		})
	}
	return &dto.ListMessagesResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}, nil
}

func (g *GrpcServer) ListFriends(ctx context.Context, in *dto.ListFriendsRequest) (*dto.ListFriendsResponse, error) {
	code, message, friends, err := g.followService.GetFriends(in.UserId)
	if err != nil {
		return &dto.ListFriendsResponse{
			Code:    code,
			Message: message,
		}, nil
	}
	var data []*dto.Friend
	for _, f := range friends {
		data = append(data, &dto.Friend{
			UserId:    f.FollowerID,
			FriendId:  f.FolloweeID,
			CreatedAt: timestamppb.New(f.CreatedAt),
		})
	}
	return &dto.ListFriendsResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}, nil
}

func (g *GrpcServer) Follow(ctx context.Context, in *dto.FollowRequest) (*dto.FollowResponse, error) {
	code, message, err := g.followService.Follow(in.UserId, in.TargetId)
	if err != nil {
		return &dto.FollowResponse{
			Code:    code,
			Message: message,
			Success: false,
		}, nil
	}
	return &dto.FollowResponse{
		Code:    code,
		Message: message,
		Success: true,
	}, nil
}

func (g *GrpcServer) Unfollow(ctx context.Context, in *dto.UnfollowRequest) (*dto.UnfollowResponse, error) {
	code, message, err := g.followService.Unfollow(in.UserId, in.TargetId)
	if err != nil {
		return &dto.UnfollowResponse{
			Code:    code,
			Message: message,
			Success: false,
		}, nil
	}
	return &dto.UnfollowResponse{
		Code:    code,
		Message: message,
		Success: true,
	}, nil
}

type chatClient struct {
	id     int32
	roomId int64
	ch     chan *dto.ChatEvent
}

type chatRoom struct {
	clients    map[int32]*chatClient
	register   chan *chatClient
	unregister chan *chatClient
	broadcast  chan *dto.ChatEvent
	mu         sync.RWMutex
}

func newChatRoom() *chatRoom {
	room := &chatRoom{
		clients:    make(map[int32]*chatClient),
		register:   make(chan *chatClient),
		unregister: make(chan *chatClient),
		broadcast:  make(chan *dto.ChatEvent),
	}
	go room.run()
	return room
}

func (r *chatRoom) run() {
	for {
		select {
		case c := <-r.register:
			r.mu.Lock()
			r.clients[c.id] = c
			r.mu.Unlock()
		case c := <-r.unregister:
			r.mu.Lock()
			delete(r.clients, c.id)
			close(c.ch)
			r.mu.Unlock()
		case event := <-r.broadcast:
			var receiverId, senderId int32
			switch e := event.Event.(type) {
			case *dto.ChatEvent_Message:
				receiverId = e.Message.ReceiverId
				senderId = e.Message.SenderId
			case *dto.ChatEvent_Typing:
				receiverId = e.Typing.ReceiverId
				senderId = e.Typing.SenderId
			case *dto.ChatEvent_Unsend:
				receiverId = e.Unsend.ReceiverId
				senderId = e.Unsend.SenderId
			}
			r.mu.RLock()
			if c, ok := r.clients[receiverId]; ok {
				c.ch <- event
			}
			if senderId != receiverId {
				if c, ok := r.clients[senderId]; ok {
					c.ch <- event
				}
			}
			r.mu.RUnlock()
		}
	}
}

// --- Room Manager ---

type chatRoomManager struct {
	rooms map[int64]*chatRoom
	mu    sync.RWMutex
}

func newChatRoomManager() *chatRoomManager {
	return &chatRoomManager{
		rooms: make(map[int64]*chatRoom),
	}
}

func (m *chatRoomManager) getRoom(roomId int64) *chatRoom {
	m.mu.RLock()
	room, ok := m.rooms[roomId]
	m.mu.RUnlock()
	if !ok {
		m.mu.Lock()
		room = newChatRoom()
		m.rooms[roomId] = room
		m.mu.Unlock()
	}
	return room
}

var globalRoomManager = newChatRoomManager()

// --- gRPC Methods ---

func (g *GrpcServer) SendChatEvent(ctx context.Context, event *dto.ChatEvent) (*dto.SendMessageResponse, error) {
	resp := &dto.SendMessageResponse{
		Code:    0,
		Message: "ok",
	}

	var senderId, receiverId int32
	var roomId int64

	switch e := event.Event.(type) {
	case *dto.ChatEvent_Message:
		senderId = e.Message.SenderId
		receiverId = e.Message.ReceiverId

		roomId = calculateRoomId(senderId, receiverId)

		code, message, savedMsg, svcErr := g.chatService.SendMessage(ctx, senderId, receiverId, e.Message.Content)
		resp.Code = code
		resp.Message = message
		if svcErr != nil {
			return resp, nil
		}

		// Create response without roomId for client
		clientMsg := &dto.ChatMessage{
			SenderId:   savedMsg.SenderID,
			ReceiverId: savedMsg.ReceiverID,
			Content:    savedMsg.Content,
			CreatedAt:  timestamppb.New(savedMsg.CreatedAt),
		}
		resp.Data = clientMsg

		//// Create internal message with roomId for server-side routing
		//internalMsg := *clientMsg
		//// roomId is only used for internal routing, kept in memory but not sent to client
		//
		//// Broadcast the message to the room
		globalRoomManager.getRoom(roomId).broadcast <- &dto.ChatEvent{
			Event: &dto.ChatEvent_Message{
				Message: clientMsg,
			},
		}

	case *dto.ChatEvent_Typing:
		senderId = e.Typing.SenderId
		receiverId = e.Typing.ReceiverId

		// Calculate roomId deterministically from user IDs
		roomId = calculateRoomId(senderId, receiverId)

		// Create typing notification without exposing roomId to client
		typingEvent := &dto.TypingNotification{
			SenderId:   senderId,
			ReceiverId: receiverId,
			IsTyping:   e.Typing.IsTyping,
		}

		globalRoomManager.getRoom(roomId).broadcast <- &dto.ChatEvent{
			Event: &dto.ChatEvent_Typing{
				Typing: typingEvent,
			},
		}

	case *dto.ChatEvent_Unsend:
		senderId = e.Unsend.SenderId
		receiverId = e.Unsend.ReceiverId

		// Calculate roomId deterministically from user IDs
		roomId = calculateRoomId(senderId, receiverId)

		// Unsend the message
		code, message, err := g.chatService.UnsendMessage(ctx, e.Unsend.MessageId, senderId)
		resp.Code = code
		resp.Message = message

		if err == nil && code == 0 {
			// Create unsend notification without exposing roomId to client
			unsendEvent := &dto.UnsendNotification{
				MessageId:  e.Unsend.MessageId,
				SenderId:   senderId,
				ReceiverId: receiverId,
			}

			// Broadcast the unsend event
			globalRoomManager.getRoom(roomId).broadcast <- &dto.ChatEvent{
				Event: &dto.ChatEvent_Unsend{
					Unsend: unsendEvent,
				},
			}
		}
	}

	return resp, nil
}

func (g *GrpcServer) ListFollowing(ctx context.Context, in *dto.ListFollowingRequest) (*dto.ListFollowingResponse, error) {
	code, message, follows, err := g.followService.GetFollowing(in.UserId)
	if err != nil {
		return &dto.ListFollowingResponse{
			Code:    code,
			Message: message,
		}, nil
	}
	var data []*dto.Follow
	for _, f := range follows {
		data = append(data, &dto.Follow{
			Id:         f.ID,
			FollowerId: f.FollowerID,
			FolloweeId: f.FolloweeID,
			CreatedAt:  timestamppb.New(f.CreatedAt),
		})
	}
	return &dto.ListFollowingResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}, nil
}

func (g *GrpcServer) ListFollowers(ctx context.Context, in *dto.ListFollowingRequest) (*dto.ListFollowingResponse, error) {
	code, message, follows, err := g.followService.GetFollowers(in.UserId)
	if err != nil {
		return &dto.ListFollowingResponse{
			Code:    code,
			Message: message,
		}, nil
	}
	var data []*dto.Follow
	for _, f := range follows {
		data = append(data, &dto.Follow{
			Id:         f.ID,
			FollowerId: f.FollowerID,
			FolloweeId: f.FolloweeID,
			CreatedAt:  timestamppb.New(f.CreatedAt),
		})
	}
	return &dto.ListFollowingResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}, nil
}

// SubscribeChatEvents implements server streaming for realtime chat
func (g *GrpcServer) SubscribeChatEvents(req *dto.ListMessagesRequest, stream controller.SocialService_SubscribeChatEventsServer) error {
	// Input validation
	if req == nil {
		return fmt.Errorf("invalid request: request cannot be nil")
	}

	clientId := req.UserId
	otherUserId := req.OtherUserId

	// Validate user IDs
	if clientId <= 0 {
		return fmt.Errorf("invalid user ID: %d", clientId)
	}

	if otherUserId <= 0 {
		return fmt.Errorf("invalid recipient user ID: %d", otherUserId)
	}

	// For private chat between two users, create a deterministic room ID
	// by sorting the user IDs and concatenating them
	var smallerId, largerId int32
	if clientId < otherUserId {
		smallerId = clientId
		largerId = otherUserId
	} else {
		smallerId = otherUserId
		largerId = clientId
	}

	// Create a deterministic room ID from the two user IDs
	roomId := int64(smallerId)*100000 + int64(largerId)

	// Set up panic recovery to prevent server crashes from affecting other connections
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in SubscribeChatEvents: %v\n", r)
		}
	}()

	// Get or create the chat room
	room := globalRoomManager.getRoom(roomId)
	if room == nil {
		return fmt.Errorf("failed to create chat room for users %d and %d", clientId, otherUserId)
	}

	// Create client
	myClient := &chatClient{
		id:     clientId,
		roomId: roomId,
		ch:     make(chan *dto.ChatEvent, 10),
	}

	// Register the client with the room
	room.register <- myClient

	// Ensure client is unregistered when done
	defer func() {
		select {
		case room.unregister <- myClient:
			// Successfully unregistered
		default:
			// Room's unregister channel might be full or closed
			fmt.Printf("Warning: Could not unregister client %d from room %d\n", clientId, roomId)
		}
	}()

	ctx := stream.Context()

	// Connection established log
	fmt.Printf("Chat connection established: User %d connected to room %d with user %d\n",
		clientId, roomId, otherUserId)

	for {
		select {
		case event, ok := <-myClient.ch:
			if !ok {
				// Channel was closed
				fmt.Printf("Client channel closed for user %d in room %d\n", clientId, roomId)
				return nil
			}

			// Send the event to the client
			if err := stream.Send(event); err != nil {
				fmt.Printf("Error sending event to user %d: %v\n", clientId, err)
				return fmt.Errorf("failed to send chat event: %v", err)
			}

		case <-ctx.Done():
			// Context canceled (client disconnected or timeout)
			err := ctx.Err()
			fmt.Printf("Client context done for user %d in room %d: %v\n", clientId, roomId, err)
			return err
		}
	}
}

func (g *GrpcServer) ListUserMessage(ctx context.Context, in *dto.ListUserMessageRequest) (*dto.ListUserMessageResponse, error) {
	code, message, partners, err := g.chatService.ListUserMessage(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &dto.ListUserMessageResponse{
		Code:    code,
		Message: message,
		UserId:  partners,
	}, nil
}
