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

type LikeServiceImpl struct {
	likeRepo  out.LikeRepository
	videoRepo out.VideoRepository
	cache     out.CacheRepository
}

func NewLikeService(likeRepo out.LikeRepository, videoRepo out.VideoRepository, cache out.CacheRepository) *LikeServiceImpl {
	return &LikeServiceImpl{
		likeRepo:  likeRepo,
		videoRepo: videoRepo,
		cache:     cache,
	}
}

func likeCacheKey(userId int32) string {
	return fmt.Sprintf("likes:user:%d", userId)
}

func (s *LikeServiceImpl) LikeVideo(ctx context.Context, userId int32, videoId int32) (int32, string, error) {
	video, err := s.videoRepo.FindById(ctx, videoId)
	if err != nil {
		return int32(codes.Internal), "Internal server error", err
	}
	if video == nil {
		return int32(codes.NotFound), "Video not found", errors.New("video not found")
	}

	likes, err := s.likeRepo.FindByUserId(ctx, userId)
	if err != nil {
		return int32(codes.Internal), "Internal server error", err
	}

	for _, like := range likes {
		if like.VideoID == videoId {
			return int32(codes.AlreadyExists), "User has already liked this video", errors.New("user has already liked this video")
		}
	}

	newLike := &domain.Like{
		UserID:  userId,
		VideoID: videoId,
	}
	if err := s.likeRepo.Create(ctx, newLike); err != nil {
		return int32(codes.Internal), "Failed to like video", err
	}
	likeCacheKey := fmt.Sprintf("like:video:%d", videoId)
	s.cache.Delete(likeCacheKey)
	helper.InvalidateFeedCache(s.cache, userId)
	return int32(codes.OK), "Video liked successfully", nil
}

func (s *LikeServiceImpl) UnlikeVideo(ctx context.Context, userId int32, videoId int32) (int32, string, error) {
	video, err := s.videoRepo.FindById(ctx, videoId)
	if err != nil {
		return int32(codes.Internal), "Internal server error", err
	}
	if video == nil {
		return int32(codes.NotFound), "Video not found", errors.New("video not found")
	}

	likes, err := s.likeRepo.FindByUserId(ctx, userId)
	if err != nil {
		return int32(codes.Internal), "Internal server error", err
	}

	var likeToDelete *domain.Like
	for _, like := range likes {
		if like.VideoID == videoId {
			likeToDelete = like
			break
		}
	}

	if likeToDelete == nil {
		return int32(codes.NotFound), "User has not liked this video", errors.New("user has not liked this video")
	}

	if err := s.likeRepo.Delete(ctx, likeToDelete.ID); err != nil {
		return int32(codes.Internal), "Failed to unlike video", err
	}
	likeCacheKey := fmt.Sprintf("like:video:%d", videoId)
	s.cache.Delete(likeCacheKey)
	helper.InvalidateFeedCache(s.cache, userId)
	return int32(codes.OK), "Video unliked successfully", nil
}

func (s *LikeServiceImpl) IsVideoLiked(ctx context.Context, userId int32, videoId int32) (bool, error) {
	var cached []*domain.Like
	if err := s.cache.Get(likeCacheKey(userId), &cached); err == nil && cached != nil {
		for _, like := range cached {
			if like.VideoID == videoId {
				return true, nil
			}
		}
		return false, nil
	}
	likes, err := s.likeRepo.FindByUserId(ctx, userId)
	if err != nil {
		return false, err
	}
	_ = s.cache.Set(likeCacheKey(userId), likes, 2*time.Minute)
	for _, like := range likes {
		if like.VideoID == videoId {
			return true, nil
		}
	}
	return false, nil
}
