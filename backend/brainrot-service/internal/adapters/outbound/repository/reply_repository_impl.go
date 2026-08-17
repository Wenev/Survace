package repository

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
	"gorm.io/gorm"
)

type ReplyRepositoryImpl struct {
	db *gorm.DB
}

func NewReplyRepository(db *gorm.DB) *ReplyRepositoryImpl {
	return &ReplyRepositoryImpl{db: db}
}

func (r *ReplyRepositoryImpl) Create(ctx context.Context, reply *domain.Reply) error {
	return r.db.WithContext(ctx).Create(reply).Error
}

func (r *ReplyRepositoryImpl) FindById(ctx context.Context, id int32) (*domain.Reply, error) {
	var reply domain.Reply
	result := r.db.WithContext(ctx).First(&reply, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &reply, nil
}

func (r *ReplyRepositoryImpl) FindByUserId(ctx context.Context, userId int32) ([]*domain.Reply, error) {
	var replies []*domain.Reply
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&replies).Error
	if err != nil {
		return nil, err
	}
	return replies, nil
}

func (r *ReplyRepositoryImpl) FindByCommentId(ctx context.Context, commentId int32) ([]*domain.Reply, error) {
	var replies []*domain.Reply
	err := r.db.WithContext(ctx).Where("comment_id = ?", commentId).Find(&replies).Error
	if err != nil {
		return nil, err
	}
	return replies, nil
}

func (r *ReplyRepositoryImpl) CountByCommentId(ctx context.Context, commentId int32) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Reply{}).Where("comment_id = ?", commentId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ReplyRepositoryImpl) DeleteUserComment(ctx context.Context, userId int32) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userId).Delete(&domain.Reply{}).Error
}

func (r *ReplyRepositoryImpl) Delete(ctx context.Context, id int32) error {
	return r.db.WithContext(ctx).Delete(&domain.Reply{}, id).Error
}
