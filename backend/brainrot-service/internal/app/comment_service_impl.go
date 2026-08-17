package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Wenev/Survace/brainrot-service/internal/adapters/outbound/cache"
	"github.com/Wenev/Survace/brainrot-service/internal/app/domain"
	"github.com/Wenev/Survace/brainrot-service/internal/app/helper"
	"github.com/Wenev/Survace/brainrot-service/ports/out"
	"google.golang.org/grpc/codes"
	"time"
)

type CommentServiceImpl struct {
	repo            out.CommentRepository
	replyRepo       out.ReplyRepository
	likeCommentRepo out.LikeCommentRepository
	cache           *cache.MemcachedConnection
}

func NewCommentService(repo out.CommentRepository, replyRepo out.ReplyRepository, likeCommentRepo out.LikeCommentRepository) *CommentServiceImpl {
	return &CommentServiceImpl{repo: repo, replyRepo: replyRepo, likeCommentRepo: likeCommentRepo, cache: cache.CacheConnection()}
}

func (s *CommentServiceImpl) AddComment(ctx context.Context, userId int32, videoId int32, text string) (int32, string, error) {
	if text == "" {
		return int32(codes.InvalidArgument), "Comment text cannot be empty", errors.New("comment text cannot be empty")
	}
	comment := &domain.Comment{
		Content: text,
		VideoID: videoId,
		UserID:  userId,
	}
	err := s.repo.Create(ctx, comment)
	if err != nil {
		return int32(codes.Internal), "Failed to add comment", err
	}
	helper.InvalidateFeedCache(s.cache, userId)
	helper.InvalidateCommentCache(s.cache, videoId)
	return int32(codes.OK), "Comment added successfully", nil
}

func (s *CommentServiceImpl) GetCommentsByVideo(ctx context.Context, videoId int32) (int32, string, []*domain.Comment, []int64, []int64, error) {
	comments, err := s.repo.FindByVideoId(ctx, videoId)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch comments", nil, nil, nil, err
	}
	replyCounts := make([]int64, len(comments))
	likeCounts := make([]int64, len(comments))
	for i, comment := range comments {
		replyCount, err := s.replyRepo.CountByCommentId(ctx, comment.ID)
		if err != nil {
			return int32(codes.Internal), "Failed to count replies", nil, nil, nil, err
		}
		replyCounts[i] = replyCount

		likeCacheKey := fmt.Sprintf("likecomments:comment:%d", comment.ID)
		var likeCount int64
		if err := s.cache.Get(likeCacheKey, &likeCount); err == nil {
			likeCounts[i] = likeCount
		} else {
			likeCount, err = s.likeCommentRepo.CountByCommentId(ctx, comment.ID)
			if err != nil {
				return int32(codes.Internal), "Failed to count comment likes", nil, nil, nil, err
			}
			likeCounts[i] = likeCount
			_ = s.cache.Set(likeCacheKey, likeCount, 2*time.Minute)
		}
	}
	return int32(codes.OK), "Comments fetched successfully", comments, replyCounts, likeCounts, nil
}

func (s *CommentServiceImpl) DeleteComment(ctx context.Context, commentId int32) (int32, string, error) {
	comment, err := s.repo.FindById(ctx, commentId)
	if err != nil {
		return int32(codes.Internal), "Failed to find comment", err
	}
	if comment == nil {
		return int32(codes.NotFound), "Comment not found", errors.New("comment not found")
	}
	err = s.repo.Delete(ctx, commentId)
	if err != nil {
		return int32(codes.Internal), "Failed to delete comment", err
	}
	helper.InvalidateFeedCache(s.cache, comment.UserID)
	helper.InvalidateCommentCache(s.cache, comment.VideoID)
	return int32(codes.OK), "Comment deleted successfully", nil
}
