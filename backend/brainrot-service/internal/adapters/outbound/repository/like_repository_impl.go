package repository

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
	"gorm.io/gorm"
	//"log"
)

type LikeRepositoryImpl struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepositoryImpl {
	return &LikeRepositoryImpl{
		db: db,
	}
}

func (r *LikeRepositoryImpl) Create(ctx context.Context, like *domain.Like) error {
	err := r.db.WithContext(ctx).Create(like).Error
	return err
}

func (r *LikeRepositoryImpl) FindById(ctx context.Context, id int32) (*domain.Like, error) {
	var like domain.Like
	result := r.db.WithContext(ctx).First(&like, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &like, nil
}

func (r *LikeRepositoryImpl) FindByUserId(ctx context.Context, userId int32) ([]*domain.Like, error) {
	var likes []*domain.Like
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&likes).Error
	if err != nil {
		return nil, err
	}
	return likes, nil
}

func (r *LikeRepositoryImpl) FindByVideoId(ctx context.Context, videoId int32) ([]*domain.Like, error) {
	var likes []*domain.Like
	err := r.db.WithContext(ctx).Where("video_id = ?", videoId).Find(&likes).Error
	if err != nil {
		return nil, err
	}
	return likes, nil
}

func (r *LikeRepositoryImpl) Update(ctx context.Context, like *domain.Like) error {
	err := r.db.WithContext(ctx).Save(&like).Error
	return err
}

func (r *LikeRepositoryImpl) Delete(ctx context.Context, id int32) error {
	result := r.db.WithContext(ctx).Delete(&domain.Like{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *LikeRepositoryImpl) DeleteByUserId(ctx context.Context, userId int32) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Delete(&domain.Like{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *LikeRepositoryImpl) DeleteByVideoId(ctx context.Context, videoId int32) error {
	return r.db.WithContext(ctx).Where("video_id = ?", videoId).Delete(&domain.Like{}).Error
}

func (r *LikeRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Like{}).Count(&count).Error
	return count, err
}

func (r *LikeRepositoryImpl) CountByVideoId(ctx context.Context, videoId int32) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Like{}).Where("video_id = ?", videoId).Count(&count).Error
	return count, err
}
