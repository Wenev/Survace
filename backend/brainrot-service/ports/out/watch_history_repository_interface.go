package out

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
)

type WatchHistoryRepository interface {
	Create(ctx context.Context, watchHistory *domain.WatchHistory) error
	FindById(ctx context.Context, id int32) (*domain.WatchHistory, error)
	FindByUserId(ctx context.Context, userId int32) ([]*domain.WatchHistory, error)
	Update(ctx context.Context, watchHistory *domain.WatchHistory) error
	DeleteByUserId(ctx context.Context, userId int32) error
	DeleteByVideoId(ctx context.Context, videoId int32) error
	CountByVideoId(ctx context.Context, videoId int32) (int64, error)
}
