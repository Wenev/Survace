package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Wenev/Survace/brainrot-service/internal/app/domain"
	"github.com/Wenev/Survace/brainrot-service/internal/app/helper"
	"github.com/Wenev/Survace/brainrot-service/ports/out"
	"google.golang.org/grpc/codes"
	"time"
)

type LikeCommentServiceImpl struct {
	likeCommentRepo out.LikeCommentRepository
	commentRepo     out.CommentRepository
	cache           out.CacheRepository
}

func NewLikeCommentService(likeCommentRepo out.LikeCommentRepository, commentRepo out.CommentRepository, cache out.CacheRepository) *LikeCommentServiceImpl {
	return &LikeCommentServiceImpl{
		likeCommentRepo: likeCommentRepo,
		commentRepo:     commentRepo,
		cache:           cache,
	}
}

func (s *LikeCommentServiceImpl) LikeComment(ctx context.Context, userId int32, commentId int32) (int32, string, error) {
	comment, err := s.commentRepo.FindById(ctx, commentId)
	if err != nil {
		return int32(codes.Internal), "Internal server error", err
	}
	if comment == nil {
		return int32(codes.NotFound), "Comment not found", errors.New("comment not found")
	}

	likes, err := s.likeCommentRepo.FindByUserId(ctx, userId)
	if err != nil {
		return int32(codes.Internal), "Internal server error", err
	}

	for _, like := range likes {
		if like.CommentID == commentId {
			return int32(codes.AlreadyExists), "User has already liked this comment", errors.New("user has already liked this comment")
		}
	}

	newLike := &domain.LikeComment{
		UserID:    userId,
		CommentID: commentId,
	}
	if err := s.likeCommentRepo.Create(ctx, newLike); err != nil {
		return int32(codes.Internal), "Failed to like comment", err
	}

	helper.InvalidateEntityCache(s.cache, "likecomments:user", userId)
	helper.InvalidateEntityCache(s.cache, "likecomments:comment", commentId)
	helper.InvalidateEntityCache(s.cache, "comments:video", comment.VideoID)

	return int32(codes.OK), "Comment liked successfully", nil
}

func (s *LikeCommentServiceImpl) UnlikeComment(ctx context.Context, userId int32, commentId int32) (int32, string, error) {
	comment, err := s.commentRepo.FindById(ctx, commentId)
	if err != nil {
		return int32(codes.Internal), "Internal server error", err
	}
	if comment == nil {
		return int32(codes.NotFound), "Comment not found", errors.New("comment not found")
	}

	likes, err := s.likeCommentRepo.FindByUserId(ctx, userId)
	if err != nil {
		return int32(codes.Internal), "Internal server error", err
	}

	var likeToDelete *domain.LikeComment
	for _, like := range likes {
		if like.CommentID == commentId {
			likeToDelete = like
			break
		}
	}

	if likeToDelete == nil {
		return int32(codes.NotFound), "User has not liked this comment", errors.New("user has not liked this comment")
	}

	if err := s.likeCommentRepo.Delete(ctx, likeToDelete.ID); err != nil {
		return int32(codes.Internal), "Failed to unlike comment", err
	}

	helper.InvalidateEntityCache(s.cache, "likecomments:user", userId)
	helper.InvalidateEntityCache(s.cache, "likecomments:comment", commentId)
	helper.InvalidateEntityCache(s.cache, "comments:video", comment.VideoID)

	return int32(codes.OK), "Comment unliked successfully", nil
}

func (s *LikeCommentServiceImpl) IsCommentLiked(ctx context.Context, userId int32, commentId int32) (bool, error) {
	var cachedByUser []*domain.LikeComment
	userKey := fmt.Sprintf("likecomments:user:%d", userId)
	if err := s.cache.Get(userKey, &cachedByUser); err == nil && cachedByUser != nil {
		for _, like := range cachedByUser {
			if like.CommentID == commentId {
				return true, nil
			}
		}
		return false, nil
	}
	var cachedByComment []*domain.LikeComment
	commentKey := fmt.Sprintf("likecomments:comment:%d", commentId)
	if err := s.cache.Get(commentKey, &cachedByComment); err == nil && cachedByComment != nil {
		for _, like := range cachedByComment {
			if like.UserID == userId {
				return true, nil
			}
		}
		return false, nil
	}
	likes, err := s.likeCommentRepo.FindByUserId(ctx, userId)
	if err != nil {
		return false, err
	}
	_ = s.cache.Set(userKey, likes, 2*time.Minute)
	returnValue := false
	for _, like := range likes {
		if like.CommentID == commentId {
			_ = s.cache.Set(commentKey, []*domain.LikeComment{like}, 2*time.Minute)
			returnValue = true
			break
		}
	}
	return returnValue, nil
}
