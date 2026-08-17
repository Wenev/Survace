package out

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
)

type LikeRepository interface {
	Create(ctx context.Context, like *domain.Like) error
	FindById(ctx context.Context, id int32) (*domain.Like, error)
	FindByUserId(ctx context.Context, userId int32) ([]*domain.Like, error)
	FindByVideoId(ctx context.Context, videoId int32) ([]*domain.Like, error)
	Update(ctx context.Context, like *domain.Like) error
	Delete(ctx context.Context, id int32) error
	DeleteByUserId(ctx context.Context, userId int32) error
	DeleteByVideoId(ctx context.Context, videoId int32) error
	CountByVideoId(ctx context.Context, videoId int32) (int64, error)
}
