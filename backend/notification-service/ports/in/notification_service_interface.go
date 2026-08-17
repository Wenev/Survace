package in

import (
	"context"
	"github.com/Wenev/Survace/notification-service/internal/app/domain"
)

type NotificationService interface {
	GetNotificationsByUserID(ctx context.Context, userID int32) (int32, string, []*domain.Notification, error)
	MarkAsRead(ctx context.Context, notificationID int32) (int32, string, error)
	MarkAsUnread(ctx context.Context, notificationID int32) (int32, string, error)
	DeleteNotification(ctx context.Context, notificationID int32) (int32, string, error)
	DeleteAllNotificationsByUserID(ctx context.Context, userID int32) (int32, string, error)
	CreateNotification(ctx context.Context, notif *domain.Notification) (int32, string, error)
}
