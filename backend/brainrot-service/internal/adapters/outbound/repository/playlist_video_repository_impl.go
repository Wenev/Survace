package repository

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
	"gorm.io/gorm"
)

type PlaylistVideoRepositoryImpl struct {
	db *gorm.DB
}

func NewPlaylistVideoRepository(db *gorm.DB) *PlaylistVideoRepositoryImpl {
	return &PlaylistVideoRepositoryImpl{
		db: db,
	}
}

func (r *PlaylistVideoRepositoryImpl) Create(ctx context.Context, playlistVideo *domain.PlaylistVideo) error {
	return r.db.WithContext(ctx).Create(playlistVideo).Error
}

func (r *PlaylistVideoRepositoryImpl) FindById(ctx context.Context, id int32) (*domain.PlaylistVideo, error) {
	var playlistVideo domain.PlaylistVideo
	result := r.db.WithContext(ctx).
		Preload("Video").
		Preload("Playlist").
		First(&playlistVideo, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &playlistVideo, nil
}

func (r *PlaylistVideoRepositoryImpl) FindByPlaylistId(ctx context.Context, playlistId int32) ([]*domain.PlaylistVideo, error) {
	var playlistVideos []*domain.PlaylistVideo
	err := r.db.WithContext(ctx).
		Where("playlist_id = ?", playlistId).
		Preload("Video").
		Order("order_index ASC").
		Find(&playlistVideos).Error

	if err != nil {
		return nil, err
	}
	return playlistVideos, nil
}

func (r *PlaylistVideoRepositoryImpl) FindByVideoId(ctx context.Context, videoId int32) ([]*domain.PlaylistVideo, error) {
	var playlistVideos []*domain.PlaylistVideo
	err := r.db.WithContext(ctx).
		Where("video_id = ?", videoId).
		Preload("Playlist").
		Find(&playlistVideos).Error

	if err != nil {
		return nil, err
	}
	return playlistVideos, nil
}

func (r *PlaylistVideoRepositoryImpl) Update(ctx context.Context, playlistVideo *domain.PlaylistVideo) error {
	return r.db.WithContext(ctx).Save(playlistVideo).Error
}

func (r *PlaylistVideoRepositoryImpl) UpdateOrder(ctx context.Context, playlistID int32, videoID int32, newOrder int32) error {
	var playlistVideo domain.PlaylistVideo

	result := r.db.WithContext(ctx).
		Where("playlist_id = ? AND video_id = ?", playlistID, videoID).
		First(&playlistVideo)

	if result.Error != nil {
		return result.Error
	}

	playlistVideo.Order = newOrder

	return r.db.WithContext(ctx).Save(&playlistVideo).Error
}

func (r *PlaylistVideoRepositoryImpl) Delete(ctx context.Context, id int32) error {
	result := r.db.WithContext(ctx).Delete(&domain.PlaylistVideo{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PlaylistVideoRepositoryImpl) DeleteByPlaylistId(ctx context.Context, playlistId int32) error {
	err := r.db.WithContext(ctx).
		Where("playlist_id = ?", playlistId).
		Delete(&domain.PlaylistVideo{}).Error

	if err != nil {
		return err
	}
	return nil
}

func (r *PlaylistVideoRepositoryImpl) DeleteByVideoId(ctx context.Context, videoId int32) error {
	err := r.db.WithContext(ctx).
		Where("video_id = ?", videoId).
		Delete(&domain.PlaylistVideo{}).Error

	return err
}
