package out

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/stream-service/internal/app/domain"
)

type StreamChatRepository interface {
	Create(ctx context.Context, message *domain.StreamChat) error
	FindByCallerId(ctx context.Context, callerId string) ([]*domain.StreamChat, error)
}
