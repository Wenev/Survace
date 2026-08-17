package repository

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/notification-service/internal/app/domain"
	"gorm.io/gorm"
)

type NotificationRepositoryImpl struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepositoryImpl {
	return &NotificationRepositoryImpl{db: db}
}

func (r *NotificationRepositoryImpl) Create(ctx context.Context, notification *domain.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *NotificationRepositoryImpl) FindByID(ctx context.Context, id int32) (*domain.Notification, error) {
	var notif domain.Notification
	err := r.db.WithContext(ctx).First(&notif, id).Error
	if err != nil {
		return nil, err
	}
	return &notif, nil
}

func (r *NotificationRepositoryImpl) FindByUserID(ctx context.Context, userID int32) ([]*domain.Notification, error) {
	var notifs []*domain.Notification
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&notifs).Error
	if err != nil {
		return nil, err
	}
	return notifs, nil
}

func (r *NotificationRepositoryImpl) Update(ctx context.Context, notification *domain.Notification) error {
	return r.db.WithContext(ctx).Save(notification).Error
}

func (r *NotificationRepositoryImpl) Delete(ctx context.Context, id int32) error {
	return r.db.WithContext(ctx).Delete(&domain.Notification{}, id).Error
}

func (r *NotificationRepositoryImpl) DeleteAllByUserID(ctx context.Context, userID int32) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.Notification{}).Error
}

func (r *NotificationRepositoryImpl) SetRead(ctx context.Context, notificationID int32, read bool) error {
	return r.db.WithContext(ctx).Model(&domain.Notification{}).Where("id = ?", notificationID).Update("read", read).Error
}
