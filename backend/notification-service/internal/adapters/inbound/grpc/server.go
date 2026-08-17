package grpc

import (
	"context"
	"fmt"
	"github.com/Wenev/Survace/notification-service/internal/app/domain"
	out "github.com/Wenev/Survace/notification-service/ports/in"
	"net"

	"github.com/Wenev/Survace/proto/gen/controller"
	"github.com/Wenev/Survace/proto/gen/dto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	grpcPort int
	server   *grpc.Server
	controller.NotificationServiceServer
	notifService out.NotificationService
}

func NewGrpcServer(grpcPort int, notifService out.NotificationService) *GrpcServer {
	return &GrpcServer{
		grpcPort:     grpcPort,
		server:       grpc.NewServer(grpc.MaxRecvMsgSize(100 * 1024 * 1024)),
		notifService: notifService,
	}
}

func (g *GrpcServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", g.grpcPort))
	if err != nil {
		return err
	}
	controller.RegisterNotificationServiceServer(g.server, g)
	return g.server.Serve(listener)
}

func (g *GrpcServer) Stop() {
	g.server.GracefulStop()
}

func (g *GrpcServer) GetNotificationsByUserID(ctx context.Context, req *dto.GetNotificationsByUserIDRequest) (*dto.GetNotificationsByUserIDResponse, error) {
	code, msg, notifications, err := g.notifService.GetNotificationsByUserID(ctx, req.UserId)
	if err != nil {
		return &dto.GetNotificationsByUserIDResponse{
			Code:          code,
			Message:       msg,
			Notifications: nil,
		}, err
	}
	var notifDtos []*dto.Notification
	for _, n := range notifications {
		notifDtos = append(notifDtos, &dto.Notification{
			Id:        n.ID,
			UserId:    n.UserID,
			Type:      n.Type,
			Content:   n.Content,
			Read:      n.Read,
			CreatedAt: timestamppb.New(n.CreatedAt),
		})
	}
	return &dto.GetNotificationsByUserIDResponse{
		Code:          code,
		Message:       msg,
		Notifications: notifDtos,
	}, nil
}

func (g *GrpcServer) MarkAsRead(ctx context.Context, req *dto.MarkNotificationRequest) (*dto.MarkNotificationResponse, error) {
	code, msg, err := g.notifService.MarkAsRead(ctx, req.NotificationId)
	return &dto.MarkNotificationResponse{
		Code:    code,
		Message: msg,
	}, err
}

func (g *GrpcServer) MarkAsUnread(ctx context.Context, req *dto.MarkNotificationRequest) (*dto.MarkNotificationResponse, error) {
	code, msg, err := g.notifService.MarkAsUnread(ctx, req.NotificationId)
	return &dto.MarkNotificationResponse{
		Code:    code,
		Message: msg,
	}, err
}

func (g *GrpcServer) DeleteNotification(ctx context.Context, req *dto.DeleteNotificationRequest) (*dto.DeleteNotificationResponse, error) {
	code, msg, err := g.notifService.DeleteNotification(ctx, req.NotificationId)
	return &dto.DeleteNotificationResponse{
		Code:    code,
		Message: msg,
	}, err
}

func (g *GrpcServer) DeleteAllNotificationsByUserID(ctx context.Context, req *dto.DeleteAllNotificationsByUserIDRequest) (*dto.DeleteAllNotificationsByUserIDResponse, error) {
	code, msg, err := g.notifService.DeleteAllNotificationsByUserID(ctx, req.UserId)
	return &dto.DeleteAllNotificationsByUserIDResponse{
		Code:    code,
		Message: msg,
	}, err
}

func (g *GrpcServer) SendNotification(ctx context.Context, req *dto.SendNotificationRequest) (*dto.
	SendNotificationResponse, error) {
	notif := &domain.Notification{
		UserID:  req.UserId,
		Type:    req.Type,
		Content: req.Content,
		Read:    false,
	}
	code, msg, err := g.notifService.CreateNotification(ctx, notif)
	if err != nil {
		return &dto.SendNotificationResponse{
			Code:         code,
			Message:      msg,
			Notification: nil,
		}, err
	}
	notifDto := &dto.Notification{
		Id:        notif.ID,
		UserId:    notif.UserID,
		Type:      notif.Type,
		Content:   notif.Content,
		Read:      notif.Read,
		CreatedAt: timestamppb.New(notif.CreatedAt),
	}
	return &dto.SendNotificationResponse{
		Code:         code,
		Message:      msg,
		Notification: notifDto,
	}, nil
}
