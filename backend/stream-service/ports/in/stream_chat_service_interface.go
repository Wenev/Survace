package in

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/stream-service/internal/app/domain"
)

type StreamChatService interface {
	SendMessage(ctx context.Context, callId string, senderId int32, content string) (int32, string, *domain.StreamChat, error)
	GetMessagesByCallID(ctx context.Context, callId string) (int32, string, []*domain.StreamChat, error)
}
