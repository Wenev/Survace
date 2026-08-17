package out

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/notification-service/internal/app/domain"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
	FindByID(ctx context.Context, id int32) (*domain.Notification, error)
	FindByUserID(ctx context.Context, userID int32) ([]*domain.Notification, error)
	Update(ctx context.Context, notification *domain.Notification) error
	Delete(ctx context.Context, id int32) error
	DeleteAllByUserID(ctx context.Context, userID int32) error
}
