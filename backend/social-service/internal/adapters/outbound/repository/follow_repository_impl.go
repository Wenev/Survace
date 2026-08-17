package repository

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/backend/social-service/internal/app/domain"
	"gorm.io/gorm"
)

type FollowRepositoryImpl struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepositoryImpl {
	return &FollowRepositoryImpl{db: db}
}

func (r *FollowRepositoryImpl) Create(ctx context.Context, follow *domain.Follow) error {
	return r.db.WithContext(ctx).Create(follow).Error
}

func (r *FollowRepositoryImpl) FindById(ctx context.Context, id int32) (*domain.Follow, error) {
	var follow domain.Follow
	result := r.db.WithContext(ctx).First(&follow, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &follow, nil
}

func (r *FollowRepositoryImpl) FindByFollowerId(ctx context.Context, followerId int32) ([]*domain.Follow, error) {
	var follows []*domain.Follow
	err := r.db.WithContext(ctx).Where("follower_id = ?", followerId).Find(&follows).Error
	if err != nil {
		return nil, err
	}
	return follows, nil
}

func (r *FollowRepositoryImpl) FindByFolloweeId(ctx context.Context, followeeId int32) ([]*domain.Follow, error) {
	var follows []*domain.Follow
	err := r.db.WithContext(ctx).Where("followee_id = ?", followeeId).Find(&follows).Error
	if err != nil {
		return nil, err
	}
	return follows, nil
}

func (r *FollowRepositoryImpl) Update(ctx context.Context, follow *domain.Follow) error {
	return r.db.WithContext(ctx).Save(follow).Error
}

func (r *FollowRepositoryImpl) Delete(ctx context.Context, id int32) error {
	return r.db.WithContext(ctx).Delete(&domain.Follow{}, id).Error
}

func (r *FollowRepositoryImpl) DeleteUserFollows(ctx context.Context, userId int32) error {
	return r.db.WithContext(ctx).Where("follower_id = ? OR followee_id = ?", userId, userId).Delete(&domain.Follow{}).Error
}
