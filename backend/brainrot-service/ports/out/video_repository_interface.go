package out

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
)

type VideoRepository interface {
	UploadMinIO(ctx context.Context, objectName string, data []byte, contentType string) (string, error)
	Create(ctx context.Context, video *domain.Video) error
	FindById(ctx context.Context, id int32) (*domain.Video, error)
	FindByUserId(ctx context.Context, userId int32) ([]*domain.Video, error)
	Update(ctx context.Context, video *domain.Video) error
	Delete(ctx context.Context, id int32) error
	DeleteByUserId(ctx context.Context, userId int32) error
	FindNewAndPopular(ctx context.Context, limit int, offset int) ([]*domain.Video, error)
	FindWatchedVideoIDs(ctx context.Context, userId int32) ([]int32, error)
	FindRandomVideos(ctx context.Context, limit int, offset int) ([]*domain.Video, error)
	FindByUserIdAndAll(ctx context.Context, userId int32) ([]*domain.Video, error)
}
