package in

import "github.com/Acad600-TPA/WEB-WE-251/backend/social-service/internal/app/domain"

type FollowService interface {
	GetFollowers(userID int32) (int32, string, []*domain.Follow, error)
	GetFollowing(userID int32) (int32, string, []*domain.Follow, error)
	GetFriends(userID int32) (int32, string, []*domain.Follow, error)
	Follow(userID, targetID int32) (int32, string, error)
	Unfollow(userID, targetID int32) (int32, string, error)
	DeleteUserFollowing(userID int32) (int32, string, error)
}
