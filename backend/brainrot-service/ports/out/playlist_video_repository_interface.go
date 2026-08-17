package out

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
)

type PlaylistVideoRepository interface {
	Create(ctx context.Context, playlistVideo *domain.PlaylistVideo) error
	FindById(ctx context.Context, id int32) (*domain.PlaylistVideo, error)
	FindByPlaylistId(ctx context.Context, playlistId int32) ([]*domain.PlaylistVideo, error)
	FindByVideoId(ctx context.Context, videoId int32) ([]*domain.PlaylistVideo, error)
	Update(ctx context.Context, playlistVideo *domain.PlaylistVideo) error
	UpdateOrder(ctx context.Context, playlistID int32, videoID int32, newOrder int32) error
	Delete(ctx context.Context, id int32) error
	DeleteByPlaylistId(ctx context.Context, playlistId int32) error
	DeleteByVideoId(ctx context.Context, videoId int32) error
}
