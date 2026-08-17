package in

import (
	"context"
	"github.com/Wenev/Survace/backend/social-service/internal/app/domain"
)

type ChatService interface {
	SendMessage(ctx context.Context, senderId int32, receiverId int32, content string) (int32, string, *domain.Message,
		error)
	GetChatWithUser(ctx context.Context, userId int32, otherUserId int32) (int32, string, []*domain.Message, error)
	UnsendMessage(ctx context.Context, messageId int32, senderId int32) (int32, string, error)
	ListUserMessage(ctx context.Context, userId int32) (int32, string, []int32, error)
}
