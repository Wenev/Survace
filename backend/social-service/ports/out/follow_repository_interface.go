package out

import (
	"context"
	"github.com/Wenev/Survace/backend/social-service/internal/app/domain"
)

type FollowRepository interface {
	Create(ctx context.Context, follow *domain.Follow) error
	FindById(ctx context.Context, id int32) (*domain.Follow, error)
	FindByFollowerId(ctx context.Context, followerId int32) ([]*domain.Follow, error)
	FindByFolloweeId(ctx context.Context, followeeId int32) ([]*domain.Follow, error)
	Update(ctx context.Context, follow *domain.Follow) error
	Delete(ctx context.Context, id int32) error
	DeleteUserFollows(ctx context.Context, userId int32) error
}
