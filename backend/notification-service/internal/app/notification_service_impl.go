package app

import (
	"context"
	"fmt"
	"github.com/Acad600-TPA/WEB-WE-251/notification-service/internal/adapters/outbound/cache"
	"github.com/Acad600-TPA/WEB-WE-251/notification-service/internal/app/domain"
	"github.com/Acad600-TPA/WEB-WE-251/notification-service/ports/in"
	"github.com/Acad600-TPA/WEB-WE-251/notification-service/ports/out"
	"google.golang.org/grpc/codes"
	"time"
)

type NotificationServiceImpl struct {
	repo  out.NotificationRepository
	cache *cache.MemcachedConnection
}

func NewNotificationService(repo out.NotificationRepository, cacheConn *cache.MemcachedConnection) in.NotificationService {
	return &NotificationServiceImpl{repo: repo, cache: cacheConn}
}

func (s *NotificationServiceImpl) GetNotificationsByUserID(ctx context.Context, userID int32) (int32, string, []*domain.Notification, error) {
	cacheKey := fmt.Sprintf("notifications:user:%d", userID)
	var notifs []*domain.Notification
	err := s.cache.Get(cacheKey, &notifs)
	if err == nil && notifs != nil {
		return int32(codes.OK), "Notifications fetched from cache", notifs, nil
	}
	notifs, err = s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch notifications", nil, err
	}
	_ = s.cache.Set(cacheKey, notifs, 5*time.Minute)
	return int32(codes.OK), "Notifications fetched successfully", notifs, nil
}

func (s *NotificationServiceImpl) MarkAsRead(ctx context.Context, notificationID int32) (int32, string, error) {
	notification, err := s.repo.FindByID(ctx, notificationID)
	if err != nil {
		return int32(codes.NotFound), "Notification not found", err
	}
	notification.Read = true
	if err := s.repo.Update(ctx, notification); err != nil {
		return int32(codes.Internal), "Failed to mark notification as read", err
	}
	cacheKey := fmt.Sprintf("notifications:user:%d", notification.UserID)
	_ = s.cache.Delete(cacheKey)
	return int32(codes.OK), "Notification marked as read", nil
}

func (s *NotificationServiceImpl) MarkAsUnread(ctx context.Context, notificationID int32) (int32, string, error) {
	notification, err := s.repo.FindByID(ctx, notificationID)
	if err != nil {
		return int32(codes.NotFound), "Notification not found", err
	}
	notification.Read = false
	if err := s.repo.Update(ctx, notification); err != nil {
		return int32(codes.Internal), "Failed to mark notification as unread", err
	}
	cacheKey := fmt.Sprintf("notifications:user:%d", notification.UserID)
	_ = s.cache.Delete(cacheKey)
	return int32(codes.OK), "Notification marked as unread", nil
}

func (s *NotificationServiceImpl) DeleteNotification(ctx context.Context, notificationID int32) (int32, string, error) {
	notification, err := s.repo.FindByID(ctx, notificationID)
	if err == nil {
		cacheKey := fmt.Sprintf("notifications:user:%d", notification.UserID)
		_ = s.cache.Delete(cacheKey)
	}
	if err := s.repo.Delete(ctx, notificationID); err != nil {
		return int32(codes.Internal), "Failed to delete notification", err
	}
	return int32(codes.OK), "Notification deleted successfully", nil
}

func (s *NotificationServiceImpl) DeleteAllNotificationsByUserID(ctx context.Context, userID int32) (int32, string, error) {
	if err := s.repo.DeleteAllByUserID(ctx, userID); err != nil {
		return int32(codes.Internal), "Failed to delete all notifications for user", err
	}
	cacheKey := fmt.Sprintf("notifications:user:%d", userID)
	_ = s.cache.Delete(cacheKey)
	return int32(codes.OK), "All notifications deleted for user", nil
}

func (s *NotificationServiceImpl) CreateNotification(ctx context.Context, notif *domain.Notification) (int32, string, error) {
	if notif.Type != "mention" && notif.Type != "message" && notif.Type != "new-follower" {
		return int32(codes.InvalidArgument), fmt.Sprintf("invalid notification type: %s", notif.Type), fmt.Errorf("invalid notification type: %s", notif.Type)
	}
	if err := s.repo.Create(ctx, notif); err != nil {
		return int32(codes.Internal), "Failed to create notification", err
	}
	cacheKey := fmt.Sprintf("notifications:user:%d", notif.UserID)
	s.cache.Delete(cacheKey) // Invalidate cache
	return int32(codes.OK), "Notification created successfully", nil
}
