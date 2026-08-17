package repository

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/stream-service/internal/app/domain"
	"gorm.io/gorm"
)

type StreamChatRepositoryImpl struct {
	db *gorm.DB
}

func NewStreamChatRepository(db *gorm.DB) *StreamChatRepositoryImpl {
	return &StreamChatRepositoryImpl{db: db}
}

func (r *StreamChatRepositoryImpl) Create(ctx context.Context, message *domain.StreamChat) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *StreamChatRepositoryImpl) FindByCallerId(ctx context.Context, callerId string) ([]*domain.StreamChat, error) {
	var messages []*domain.StreamChat
	err := r.db.WithContext(ctx).Where("call_id = ?", callerId).Order("created_at asc").Find(&messages).Error
	return messages, err
}
