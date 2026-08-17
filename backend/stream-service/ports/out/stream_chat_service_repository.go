package out

import (
	"context"
	"github.com/Wenev/Survace/stream-service/internal/app/domain"
)

type StreamChatRepository interface {
	Create(ctx context.Context, message *domain.StreamChat) error
	FindByCallerId(ctx context.Context, callerId string) ([]*domain.StreamChat, error)
}
