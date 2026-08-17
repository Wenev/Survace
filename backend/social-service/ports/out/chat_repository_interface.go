package out

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/backend/social-service/internal/app/domain"
)

type ChatRepository interface {
	Create(ctx context.Context, message *domain.Message) error
	FindById(ctx context.Context, id int32) (*domain.Message, error)
	FindBySenderId(ctx context.Context, senderId int32) ([]*domain.Message, error)
	FindByReceiverId(ctx context.Context, receiverId int32) ([]*domain.Message, error)
	Update(ctx context.Context, message *domain.Message) error
	Delete(ctx context.Context, id int32) error
	DeleteUserMessages(ctx context.Context, userId int32) error
}
