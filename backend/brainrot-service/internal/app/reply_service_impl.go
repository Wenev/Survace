package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/adapters/outbound/cache"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/helper"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/ports/out"
	"google.golang.org/grpc/codes"
	"time"
)

type ReplyServiceImpl struct {
	repo          out.ReplyRepository
	likeReplyRepo out.LikeReplyRepository
	cache         *cache.MemcachedConnection
}

func NewReplyService(repo out.ReplyRepository, likeReplyRepo out.LikeReplyRepository) *ReplyServiceImpl {
	return &ReplyServiceImpl{repo: repo, likeReplyRepo: likeReplyRepo, cache: cache.CacheConnection()}
}

func (s *ReplyServiceImpl) AddReply(ctx context.Context, userId int32, commentId int32, text string) (int32, string, error) {
	if text == "" {
		return int32(codes.InvalidArgument), "Reply text cannot be empty", errors.New("reply text cannot be empty")
	}
	reply := &domain.Reply{
		Content:   text,
		CommentID: commentId,
		UserID:    userId,
	}
	err := s.repo.Create(ctx, reply)
	if err != nil {
		return int32(codes.Internal), "Failed to add reply", err
	}
	helper.InvalidateReplyCache(s.cache, commentId)
	return int32(codes.OK), "Reply added successfully", nil
}

func replyCacheKey(commentId int32) string {
	return fmt.Sprintf("replies:comment:%d", commentId)
}

func (s *ReplyServiceImpl) GetRepliesByComment(ctx context.Context, commentId int32) (int32, string, []*domain.Reply, []int64, error) {
	replies, err := s.repo.FindByCommentId(ctx, commentId)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch replies", nil, nil, err
	}
	likeCounts := make([]int64, len(replies))
	for i, reply := range replies {
		likeCacheKey := fmt.Sprintf("likereplies:reply:%d", reply.ID)
		var likeCount int64
		if err := s.cache.Get(likeCacheKey, &likeCount); err == nil {
			likeCounts[i] = likeCount
		} else {
			likeCount, err = s.likeReplyRepo.CountByReplyId(ctx, reply.ID)
			if err != nil {
				return int32(codes.Internal), "Failed to count reply likes", nil, nil, err
			}
			likeCounts[i] = likeCount
			_ = s.cache.Set(likeCacheKey, likeCount, 2*time.Minute)
		}
	}
	return int32(codes.OK), "Replies fetched successfully", replies, likeCounts, nil
}

func (s *ReplyServiceImpl) DeleteReply(ctx context.Context, replyId int32) (int32, string, error) {
	commentId := int32(0)
	reply, err := s.repo.FindById(ctx, replyId)
	if err == nil && reply != nil {
		commentId = reply.CommentID
	}
	err = s.repo.Delete(ctx, replyId)
	if err != nil {
		return int32(codes.Internal), "Failed to delete reply", err
	}
	if commentId != 0 {
		helper.InvalidateReplyCache(s.cache, commentId)
	}
	return int32(codes.OK), "Reply deleted successfully", nil
}
