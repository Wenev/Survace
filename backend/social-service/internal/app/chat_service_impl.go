package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Wenev/Survace/backend/social-service/internal/app/domain"
	"github.com/Wenev/Survace/backend/social-service/ports/in"
	"github.com/Wenev/Survace/backend/social-service/ports/out"
	"google.golang.org/grpc/codes"
)

type ChatServiceImpl struct {
	chatRepo out.ChatRepository
	cache    out.CacheRepository
}

func NewChatService(chatRepo out.ChatRepository, cacheConn out.CacheRepository) in.ChatService {
	return &ChatServiceImpl{
		chatRepo: chatRepo,
		cache:    cacheConn,
	}
}

func (s *ChatServiceImpl) SendMessage(ctx context.Context, senderId int32, receiverId int32, content string) (int32, string, *domain.Message, error) {
	if content == "" {
		return int32(codes.InvalidArgument), "Message content cannot be empty", nil, fmt.Errorf("message content cannot be empty")
	}
	msg := &domain.Message{
		SenderID:   senderId,
		ReceiverID: receiverId,
		Content:    content,
		CreatedAt:  time.Now(),
	}
	err := s.chatRepo.Create(ctx, msg)
	if err != nil {
		return int32(codes.Internal), "Failed to send message", nil, err
	}
	return int32(codes.OK), "Message sent successfully", msg, nil
}

func (s *ChatServiceImpl) GetChatWithUser(ctx context.Context, userId int32, otherUserId int32) (int32, string, []*domain.Message, error) {
	cacheKey := fmt.Sprintf("chat:%d:%d", userId, otherUserId)
	var messages []*domain.Message
	if err := s.cache.Get(cacheKey, &messages); err == nil {
		return int32(codes.OK), "Messages fetched from cache", messages, nil
	}
	msgs1, err := s.chatRepo.FindBySenderId(ctx, userId)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch messages", nil, err
	}
	msgs2, err := s.chatRepo.FindBySenderId(ctx, otherUserId)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch messages", nil, err
	}
	var chat []*domain.Message
	for _, m := range msgs1 {
		if m.ReceiverID == otherUserId {
			chat = append(chat, m)
		}
	}
	for _, m := range msgs2 {
		if m.ReceiverID == userId {
			chat = append(chat, m)
		}
	}
	if err := s.cache.Set(cacheKey, chat, 2*time.Minute); err != nil {
		return int32(codes.Internal), "Failed to cache chat", nil, fmt.Errorf("failed to cache chat: %w", err)
	}
	return int32(codes.OK), "Messages fetched successfully", chat, nil
}

func (s *ChatServiceImpl) UnsendMessage(ctx context.Context, messageId int32, senderId int32) (int32, string, error) {
	msg, err := s.chatRepo.FindById(ctx, messageId)
	if err != nil {
		return int32(codes.NotFound), "Message not found", err
	}
	if msg.SenderID != senderId {
		return int32(codes.PermissionDenied), "You can only unsend your own messages", fmt.Errorf("permission denied")
	}
	err = s.chatRepo.Delete(ctx, messageId)
	if err != nil {
		return int32(codes.Internal), "Failed to unsend message", err
	}
	return int32(codes.OK), "Message unsent successfully", nil
}

func (s *ChatServiceImpl) ListUserMessage(ctx context.Context, userId int32) (int32, string, []int32, error) {
	sentMsgs, err := s.chatRepo.FindBySenderId(ctx, userId)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch sent messages", nil, err
	}
	receivedMsgs, err := s.chatRepo.FindByReceiverId(ctx, userId)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch received messages", nil, err
	}
	partnerSet := make(map[int32]struct{})
	for _, m := range sentMsgs {
		if m.ReceiverID != userId {
			partnerSet[m.ReceiverID] = struct{}{}
		}
	}
	for _, m := range receivedMsgs {
		if m.SenderID != userId {
			partnerSet[m.SenderID] = struct{}{}
		}
	}
	partners := make([]int32, 0, len(partnerSet))
	for id := range partnerSet {
		partners = append(partners, id)
	}
	return int32(codes.OK), "Unique chat partners fetched successfully", partners, nil
}
