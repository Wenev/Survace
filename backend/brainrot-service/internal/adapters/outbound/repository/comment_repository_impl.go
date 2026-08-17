package repository

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
	"gorm.io/gorm"
)

type CommentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepositoryImpl {
	return &CommentRepositoryImpl{
		db: db,
	}
}

func (r *CommentRepositoryImpl) Create(ctx context.Context, comment *domain.Comment) error {
	err := r.db.WithContext(ctx).Create(comment).Error
	return err
}

func (r *CommentRepositoryImpl) FindById(ctx context.Context, id int32) (*domain.Comment, error) {
	var comment domain.Comment
	result := r.db.WithContext(ctx).First(&comment, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &comment, nil
}

func (r *CommentRepositoryImpl) FindByUserId(ctx context.Context, userId int32) ([]*domain.Comment, error) {
	var comments []*domain.Comment
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepositoryImpl) DeleteUserComment(ctx context.Context, userId int32) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Delete(&domain.Comment{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentRepositoryImpl) Delete(ctx context.Context, id int32) error {
	result := r.db.WithContext(ctx).Delete(&domain.Comment{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CommentRepositoryImpl) FindByVideoId(ctx context.Context, videoId int32) ([]*domain.Comment, error) {
	var comments []*domain.Comment
	err := r.db.WithContext(ctx).Where("video_id = ?", videoId).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Comment{}).Count(&count).Error
	return count, err
}

func (r *CommentRepositoryImpl) CountByVideoId(ctx context.Context, videoId int32) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Comment{}).Where("video_id = ?", videoId).Count(&count).Error
	return count, err
}

func (r *CommentRepositoryImpl) DeleteByVideoId(ctx context.Context, videoId int32) error {
	return r.db.WithContext(ctx).Where("video_id = ?", videoId).Delete(&domain.Comment{}).Error
}
