package app

import (
	"context"

	"github.com/Wenev/Survace/backend/social-service/internal/app/domain"
	"github.com/Wenev/Survace/backend/social-service/ports/out"
	"google.golang.org/grpc/codes"
)

type FollowServiceImpl struct {
	repo out.FollowRepository
}

func NewFollowService(repo out.FollowRepository) *FollowServiceImpl {
	return &FollowServiceImpl{repo: repo}
}

func (s *FollowServiceImpl) GetFollowers(userID int32) (int32, string, []*domain.Follow, error) {
	ctx := context.Background()
	follows, err := s.repo.FindByFolloweeId(ctx, userID)
	if err != nil {
		return int32(codes.Internal), "Failed to get followers", nil, err
	}
	return int32(codes.OK), "Followers fetched successfully", follows, nil
}

func (s *FollowServiceImpl) GetFollowing(userID int32) (int32, string, []*domain.Follow, error) {
	ctx := context.Background()
	follows, err := s.repo.FindByFollowerId(ctx, userID)
	if err != nil {
		return int32(codes.Internal), "Failed to get following", nil, err
	}
	return int32(codes.OK), "Following fetched successfully", follows, nil
}

func (s *FollowServiceImpl) Follow(userID, targetID int32) (int32, string, error) {
	ctx := context.Background()
	follow := &domain.Follow{
		FollowerID: userID,
		FolloweeID: targetID,
	}
	err := s.repo.Create(ctx, follow)
	if err != nil {
		return int32(codes.Internal), "Failed to follow user", err
	}
	return int32(codes.OK), "Followed successfully", nil
}

func (s *FollowServiceImpl) Unfollow(userID, targetID int32) (int32, string, error) {
	ctx := context.Background()
	follows, err := s.repo.FindByFollowerId(ctx, userID)
	if err != nil {
		return int32(codes.Internal), "Failed to get following", err
	}
	for _, follow := range follows {
		if follow.FolloweeID == targetID {
			err := s.repo.Delete(ctx, follow.ID)
			if err != nil {
				return int32(codes.Internal), "Failed to unfollow user", err
			}
			return int32(codes.OK), "Unfollowed successfully", nil
		}
	}
	return int32(codes.NotFound), "Follow relationship not found", nil
}

func (s *FollowServiceImpl) DeleteUserFollowing(userID int32) (int32, string, error) {
	ctx := context.Background()
	err := s.repo.DeleteUserFollows(ctx, userID)
	if err != nil {
		return int32(codes.Internal), "Failed to delete user following", err
	}
	return int32(codes.OK), "User following deleted successfully", nil
}

func (s *FollowServiceImpl) GetFriends(userID int32) (int32, string, []*domain.Follow, error) {
	_, _, followers, err := s.GetFollowers(userID)
	if err != nil {
		return int32(codes.Internal), "Failed to get followers", nil, err
	}
	_, _, following, err := s.GetFollowing(userID)
	if err != nil {
		return int32(codes.Internal), "Failed to get following", nil, err
	}
	friendSet := make(map[int32]*domain.Follow)
	for _, follower := range followers {
		friendSet[follower.FollowerID] = follower
	}
	var friends []*domain.Follow
	for _, followee := range following {
		if _, exists := friendSet[followee.FolloweeID]; exists {
			friends = append(friends, followee)
		}
	}
	return int32(codes.OK), "Friends fetched successfully", friends, nil
}
