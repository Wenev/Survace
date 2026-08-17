package repository

import (
	"context"
	"fmt"
	"github.com/Wenev/Survace/brainrot-service/internal/app/domain"
	"gorm.io/gorm"
)

type LikeCommentRepositoryImpl struct {
	db *gorm.DB
}

func NewLikeCommentRepository(db *gorm.DB) *LikeCommentRepositoryImpl {
	return &LikeCommentRepositoryImpl{
		db: db,
	}
}

func (r *LikeCommentRepositoryImpl) Create(ctx context.Context, like *domain.LikeComment) error {
	err := r.db.WithContext(ctx).Create(like).Error
	return err
}

func (r *LikeCommentRepositoryImpl) FindById(ctx context.Context, id int32) (*domain.LikeComment, error) {
	var like domain.LikeComment
	result := r.db.WithContext(ctx).First(&like, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &like, nil
}

func (r *LikeCommentRepositoryImpl) FindByUserId(ctx context.Context, userId int32) ([]*domain.LikeComment, error) {
	var likes []*domain.LikeComment
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&likes).Error
	if err != nil {
		return nil, err
	}
	return likes, nil
}

func (r *LikeCommentRepositoryImpl) FindByCommentId(ctx context.Context, commentId int32) ([]*domain.LikeComment,
	error) {
	var likes []*domain.LikeComment
	err := r.db.WithContext(ctx).Where("comment_id = ?", commentId).Find(&likes).Error
	if err != nil {
		return nil, err
	}
	return likes, nil
}

func (r *LikeCommentRepositoryImpl) Update(ctx context.Context, like *domain.LikeComment) error {
	err := r.db.WithContext(ctx).Save(&like).Error
	return err
}

func (r *LikeCommentRepositoryImpl) Delete(ctx context.Context, id int32) error {
	result := r.db.WithContext(ctx).Delete(&domain.LikeComment{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *LikeCommentRepositoryImpl) DeleteByUserId(ctx context.Context, userId int32) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Delete(&domain.LikeComment{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *LikeCommentRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.LikeComment{}).Count(&count).Error
	return count, err
}

func (r *LikeCommentRepositoryImpl) CountByCommentId(ctx context.Context, commentId int32) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.LikeComment{}).Where("comment_id = ?", commentId).Count(&count).Error
	fmt.Println("CountByCommentId for commentId", commentId, "=", count)
	return count, err
}
