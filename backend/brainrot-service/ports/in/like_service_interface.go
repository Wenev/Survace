package in

import (
	"context"
)

type LikeService interface {
	LikeVideo(ctx context.Context, userId int32, videoId int32) (int32, string, error)
	UnlikeVideo(ctx context.Context, userId int32, videoId int32) (int32, string, error)
	IsVideoLiked(ctx context.Context, userId int32, videoId int32) (bool, error)
}
