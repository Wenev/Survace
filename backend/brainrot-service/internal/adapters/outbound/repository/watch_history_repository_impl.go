package repository

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
	"gorm.io/gorm"
)

type WatchHistoryRepositoryImpl struct {
	db *gorm.DB
}

func NewWatchHistoryRepository(db *gorm.DB) *WatchHistoryRepositoryImpl {
	return &WatchHistoryRepositoryImpl{
		db: db,
	}
}

func (r *WatchHistoryRepositoryImpl) Create(ctx context.Context, watchHistory *domain.WatchHistory) error {
	return r.db.WithContext(ctx).Create(watchHistory).Error
}

func (r *WatchHistoryRepositoryImpl) FindById(ctx context.Context, id int32) (*domain.WatchHistory, error) {
	var watchHistory domain.WatchHistory
	result := r.db.WithContext(ctx).Preload("Video").First(&watchHistory, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &watchHistory, nil
}

func (r *WatchHistoryRepositoryImpl) FindByUserId(ctx context.Context, userId int32) ([]*domain.WatchHistory, error) {
	var watchHistories []*domain.WatchHistory
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).
		Preload("Video").
		Order("created_at DESC").
		Find(&watchHistories).Error
	if err != nil {
		return nil, err
	}
	return watchHistories, nil
}

func (r *WatchHistoryRepositoryImpl) Update(ctx context.Context, watchHistory *domain.WatchHistory) error {
	return r.db.WithContext(ctx).Save(watchHistory).Error
}

func (r *WatchHistoryRepositoryImpl) DeleteByUserId(ctx context.Context, userId int32) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Delete(&domain.WatchHistory{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WatchHistoryRepositoryImpl) DeleteByVideoId(ctx context.Context, videoId int32) error {
	return r.db.WithContext(ctx).Where("video_id = ?", videoId).Delete(&domain.WatchHistory{}).Error
}

func (r *WatchHistoryRepositoryImpl) CountByVideoId(ctx context.Context, videoId int32) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.WatchHistory{}).Where("video_id = ?", videoId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
