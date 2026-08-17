package out

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	FindById(ctx context.Context, id int32) (*domain.Comment, error)
	FindByUserId(ctx context.Context, userId int32) ([]*domain.Comment, error)
	FindByVideoId(ctx context.Context, videoId int32) ([]*domain.Comment, error)
	DeleteUserComment(ctx context.Context, userId int32) error
	Delete(ctx context.Context, userId int32) error
	CountByVideoId(ctx context.Context, videoId int32) (int64, error)
	DeleteByVideoId(ctx context.Context, videoId int32) error
}
