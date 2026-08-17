package grpc

import (
	"context"
	"fmt"
	"github.com/Wenev/Survace/proto/gen/controller"
	"github.com/Wenev/Survace/proto/gen/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GalactusClient struct {
	Client             controller.GalactusControllerClient
	NotificationClient controller.NotificationServiceClient
}

func NewGalactusClient(address string, notificationAddress string) (*GalactusClient, error) {
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Galactus: %v", err)
	}
	client := controller.NewGalactusControllerClient(conn)

	notifConn, err := grpc.Dial(
		notificationAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NotificationService: %v", err)
	}
	notificationClient := controller.NewNotificationServiceClient(notifConn)

	return &GalactusClient{Client: client, NotificationClient: notificationClient}, nil
}

func (g *GalactusClient) FindByUserId(ctx context.Context, userId int32) (*dto.FindByUserIdResponse, error) {
	resp, err := g.Client.FindByUserId(ctx, &dto.FindByUserIdRequest{Id: userId})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GalactusClient) IsValid(ctx context.Context, userId int32) (*dto.IsValidResponse, error) {
	resp, err := g.Client.IsValid(ctx, &dto.IsValidRequest{UserId: userId})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GalactusClient) GetNotificationsByUserID(ctx context.Context, userId int32) (*dto.GetNotificationsByUserIDResponse, error) {
	resp, err := g.NotificationClient.GetNotificationsByUserID(ctx, &dto.GetNotificationsByUserIDRequest{UserId: userId})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GalactusClient) MarkAsRead(ctx context.Context, notificationId int32) (*dto.MarkNotificationResponse, error) {
	resp, err := g.NotificationClient.MarkAsRead(ctx, &dto.MarkNotificationRequest{NotificationId: notificationId})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GalactusClient) MarkAsUnread(ctx context.Context, notificationId int32) (*dto.MarkNotificationResponse, error) {
	resp, err := g.NotificationClient.MarkAsUnread(ctx, &dto.MarkNotificationRequest{NotificationId: notificationId})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GalactusClient) DeleteNotification(ctx context.Context, notificationId int32) (*dto.DeleteNotificationResponse, error) {
	resp, err := g.NotificationClient.DeleteNotification(ctx, &dto.DeleteNotificationRequest{NotificationId: notificationId})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GalactusClient) DeleteAllNotificationsByUserID(ctx context.Context, userId int32) (*dto.DeleteAllNotificationsByUserIDResponse, error) {
	resp, err := g.NotificationClient.DeleteAllNotificationsByUserID(ctx, &dto.DeleteAllNotificationsByUserIDRequest{UserId: userId})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GalactusClient) SendNotification(ctx context.Context, userId int32, notifType, content string) (*dto.SendNotificationResponse, error) {
	resp, err := g.NotificationClient.SendNotification(ctx, &dto.SendNotificationRequest{
		UserId:  userId,
		Type:    notifType,
		Content: content,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
