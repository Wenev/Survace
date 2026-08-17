package out

import (
	"context"
	"github.com/Wenev/Survace/brainrot-service/internal/app/domain"
)

type LikeCommentRepository interface {
	Create(ctx context.Context, like *domain.LikeComment) error
	FindById(ctx context.Context, id int32) (*domain.LikeComment, error)
	FindByUserId(ctx context.Context, userId int32) ([]*domain.LikeComment, error)
	FindByCommentId(ctx context.Context, videoId int32) ([]*domain.LikeComment, error)
	Update(ctx context.Context, like *domain.LikeComment) error
	Delete(ctx context.Context, id int32) error
	DeleteByUserId(ctx context.Context, userId int32) error
	CountByCommentId(ctx context.Context, commentId int32) (int64, error)
}
