package app

import (
	"context"
	"errors"
	"github.com/Wenev/Survace/brainrot-service/internal/adapters/outbound/cache"
	"github.com/Wenev/Survace/brainrot-service/internal/app/domain"
	"github.com/Wenev/Survace/brainrot-service/internal/app/helper"
	"github.com/Wenev/Survace/brainrot-service/ports/out"
	"google.golang.org/grpc/codes"
)

type LikeReplyServiceImpl struct {
	repo  out.LikeReplyRepository
	cache *cache.MemcachedConnection
}

func NewLikeReplyService(repo out.LikeReplyRepository) *LikeReplyServiceImpl {
	return &LikeReplyServiceImpl{
		repo:  repo,
		cache: cache.CacheConnection(),
	}
}

func (s *LikeReplyServiceImpl) LikeReply(ctx context.Context, replyId int32, userId int32) (int32, string, error) {
	like, err := s.repo.FindByReplyIdAndUserId(ctx, replyId, userId)
	if err != nil {
		return int32(codes.Internal), "Failed to check like", err
	}
	if like != nil {
		return int32(codes.AlreadyExists), "Already liked", errors.New("already liked")
	}
	likeReply := &domain.LikeReply{
		ReplyID: replyId,
		UserID:  userId,
	}
	err = s.repo.Create(ctx, likeReply)
	if err != nil {
		return int32(codes.Internal), "Failed to like reply", err
	}
	helper.InvalidateLikeReplyCache(s.cache, userId, replyId)
	return int32(codes.OK), "Reply liked successfully", nil
}

func (s *LikeReplyServiceImpl) UnlikeReply(ctx context.Context, replyId int32, userId int32) (int32, string, error) {
	err := s.repo.DeleteByReplyIdAndUserId(ctx, replyId, userId)
	if err != nil {
		return int32(codes.Internal), "Failed to unlike reply", err
	}
	helper.InvalidateLikeReplyCache(s.cache, userId, replyId)
	return int32(codes.OK), "Reply unliked successfully", nil
}

func (s *LikeReplyServiceImpl) IsReplyLiked(ctx context.Context, replyId int32, userId int32) (bool, error) {
	like, err := s.repo.FindByReplyIdAndUserId(ctx, replyId, userId)
	if err != nil {
		return false, err
	}
	return like != nil, nil
}
