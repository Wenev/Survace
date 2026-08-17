package repository

import (
	"context"
	"github.com/Wenev/Survace/brainrot-service/internal/app/domain"
	"gorm.io/gorm"
)

type PlaylistRepositoryImpl struct {
	db *gorm.DB
}

func NewPlaylistRepository(db *gorm.DB) *PlaylistRepositoryImpl {
	return &PlaylistRepositoryImpl{
		db: db,
	}
}

func (r *PlaylistRepositoryImpl) Create(ctx context.Context, playlist *domain.Playlist) error {
	return r.db.WithContext(ctx).Create(playlist).Error
}

func (r *PlaylistRepositoryImpl) FindById(ctx context.Context, id int32) (*domain.Playlist, error) {
	var playlist domain.Playlist
	result := r.db.WithContext(ctx).Preload("Videos").First(&playlist, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &playlist, nil
}

func (r *PlaylistRepositoryImpl) FindByUserId(ctx context.Context, userId int32) ([]*domain.Playlist, error) {
	var playlists []*domain.Playlist
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Preload("Videos").Find(&playlists).Error
	if err != nil {
		return nil, err
	}
	return playlists, nil
}

func (r *PlaylistRepositoryImpl) Update(ctx context.Context, playlist *domain.Playlist) error {
	return r.db.WithContext(ctx).Save(playlist).Error
}

func (r *PlaylistRepositoryImpl) Delete(ctx context.Context, id int32) error {
	result := r.db.WithContext(ctx).Delete(&domain.Playlist{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PlaylistRepositoryImpl) DeleteByUserId(ctx context.Context, userId int32) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Delete(&domain.Playlist{}).Error
	if err != nil {
		return err
	}
	return nil
}
