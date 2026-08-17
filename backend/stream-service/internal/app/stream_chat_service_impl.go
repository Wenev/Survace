package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Wenev/Survace/stream-service/internal/app/domain"
	"github.com/Wenev/Survace/stream-service/ports/in"
	"github.com/Wenev/Survace/stream-service/ports/out"
	"google.golang.org/grpc/codes"
)

type StreamChatServiceImpl struct {
	repo out.StreamChatRepository
}

func NewStreamChatService(repo out.StreamChatRepository) in.StreamChatService {
	return &StreamChatServiceImpl{repo: repo}
}

func (s *StreamChatServiceImpl) SendMessage(ctx context.Context, callId string, senderId int32, content string) (int32, string, *domain.StreamChat, error) {
	if content == "" {
		return int32(codes.InvalidArgument), "Message content cannot be empty", nil, fmt.Errorf("message content cannot be empty")
	}
	callIDInt, err := parseCallID(callId)
	if err != nil {
		return int32(codes.InvalidArgument), "Invalid callId", nil, err
	}
	msg := &domain.StreamChat{
		SenderID:  senderId,
		CallID:    callIDInt,
		Content:   content,
		CreatedAt: time.Now(),
	}
	err = s.repo.Create(ctx, msg)
	if err != nil {
		return int32(codes.Internal), "Failed to save message", nil, err
	}
	return int32(codes.OK), "Message sent successfully", msg, nil
}

func (s *StreamChatServiceImpl) GetMessagesByCallID(ctx context.Context, callId string) (int32, string, []*domain.StreamChat, error) {
	msgs, err := s.repo.FindByCallerId(ctx, callId)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch messages", nil, err
	}
	return int32(codes.OK), "Messages fetched successfully", msgs, nil
}

func parseCallID(callId string) (int32, error) {
	var id int32
	_, err := fmt.Sscanf(callId, "%d", &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
