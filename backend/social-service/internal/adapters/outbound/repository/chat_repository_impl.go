package repository

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/backend/social-service/internal/app/domain"
	"gorm.io/gorm"
)

type ChatRepositoryImpl struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepositoryImpl {
	return &ChatRepositoryImpl{db: db}
}

func (r *ChatRepositoryImpl) Create(ctx context.Context, message *domain.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *ChatRepositoryImpl) FindById(ctx context.Context, id int32) (*domain.Message, error) {
	var message domain.Message
	result := r.db.WithContext(ctx).First(&message, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &message, nil
}

func (r *ChatRepositoryImpl) FindBySenderId(ctx context.Context, senderId int32) ([]*domain.Message, error) {
	var messages []*domain.Message
	err := r.db.WithContext(ctx).Where("sender_id = ?", senderId).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *ChatRepositoryImpl) FindByReceiverId(ctx context.Context, receiverId int32) ([]*domain.Message, error) {
	var messages []*domain.Message
	err := r.db.WithContext(ctx).Where("receiver_id = ?", receiverId).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *ChatRepositoryImpl) Update(ctx context.Context, message *domain.Message) error {
	return r.db.WithContext(ctx).Save(message).Error
}

func (r *ChatRepositoryImpl) Delete(ctx context.Context, id int32) error {
	return r.db.WithContext(ctx).Delete(&domain.Message{}, id).Error
}

func (r *ChatRepositoryImpl) DeleteUserMessages(ctx context.Context, userId int32) error {
	return r.db.WithContext(ctx).Where("sender_id = ? OR receiver_id = ?", userId, userId).Delete(&domain.Message{}).Error
}
