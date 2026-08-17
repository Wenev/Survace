package repository

import (
	"context"
	"fmt"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
	"gorm.io/gorm"
)

type LikeReplyRepositoryImpl struct {
	db *gorm.DB
}

func NewLikeReplyRepository(db *gorm.DB) *LikeReplyRepositoryImpl {
	return &LikeReplyRepositoryImpl{db: db}
}

func (r *LikeReplyRepositoryImpl) Create(ctx context.Context, like *domain.LikeReply) error {
	return r.db.WithContext(ctx).Create(like).Error
}

func (r *LikeReplyRepositoryImpl) FindByReplyIdAndUserId(ctx context.Context, replyId int32, userId int32) (*domain.LikeReply, error) {
	var like domain.LikeReply
	result := r.db.WithContext(ctx).Where("reply_id = ? AND user_id = ?", replyId, userId).First(&like)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &like, nil
}

func (r *LikeReplyRepositoryImpl) Delete(ctx context.Context, id int32) error {
	return r.db.WithContext(ctx).Delete(&domain.LikeReply{}, id).Error
}

func (r *LikeReplyRepositoryImpl) DeleteByReplyIdAndUserId(ctx context.Context, replyId int32, userId int32) error {
	return r.db.WithContext(ctx).Where("reply_id = ? AND user_id = ?", replyId, userId).Delete(&domain.LikeReply{}).Error
}

func (r *LikeReplyRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.LikeReply{}).Count(&count).Error
	return count, err
}

func (r *LikeReplyRepositoryImpl) CountByReplyId(ctx context.Context, replyId int32) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.LikeReply{}).Where("reply_id = ?", replyId).Count(&count).Error
	fmt.Println("CountByReplyId for replyId", replyId, "=", count)
	return count, err
}
