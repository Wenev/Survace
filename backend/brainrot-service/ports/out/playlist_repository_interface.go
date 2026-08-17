package out

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
)

type PlaylistRepository interface {
	Create(ctx context.Context, playlist *domain.Playlist) error
	FindById(ctx context.Context, id int32) (*domain.Playlist, error)
	FindByUserId(ctx context.Context, userId int32) ([]*domain.Playlist, error)
	Update(ctx context.Context, playlist *domain.Playlist) error
	Delete(ctx context.Context, id int32) error
	DeleteByUserId(ctx context.Context, userId int32) error
}
