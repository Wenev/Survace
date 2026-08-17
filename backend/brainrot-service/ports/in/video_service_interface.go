package in

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
	"time"
)

type VideoService interface {
	UploadVideo(ctx context.Context, title string, file []byte, userId int32, enableComment bool, visibility string, thumbnail []byte, isDraft bool, postDate *time.Time) (int32, string, *domain.Video, error)
	GetVideoByID(ctx context.Context, id int32) (int32, string, *domain.Video, int64, int64, int32, error)
	GetUserVideo(ctx context.Context, userId int32) (int32, string, []*domain.Video, []int64, []int64, []int32, error)
	GetUserVideoAndDraft(ctx context.Context, userId int32) (int32, string, []*domain.Video, []int64, []int64, []int32, error)
	DeleteUserVideo(ctx context.Context, userId int32) (int32, string, error)
	GetRandomFeed(ctx context.Context, userId *int32, limit int, offset int) (int32, string, []*domain.Video, []int64, []int64, []int32, error)
	GetRandomFeedLoggedOut(ctx context.Context, limit int, offset int) (int32, string, []*domain.Video, []int64, []int64, []int32, error)
	WatchVideo(ctx context.Context, userId int32, videoId int32) (int32, string, error)
	SearchVideo(ctx context.Context, query string, threshold float64, limit int, offset int) (int32, string, []*domain.Video, []int64, []int64, []int32, error)
	DeleteVideoByID(ctx context.Context, videoId int32) (int32, string, error)
	UpdateVideo(ctx context.Context, id int32, userId int32, title string, url string, objectName string, enableComment bool, visibility string, thumbnail []byte, postedAt *time.Time, isDraft bool) (int32, string, *domain.Video, error)
	FollowingVideo(ctx context.Context, userId int32, limit int, offset int) (int32, string, []*domain.Video, []int64, []int64, []int32, error)
	FriendVideo(ctx context.Context, userId int32, limit int, offset int) (int32, string, []*domain.Video, []int64, []int64, []int32, error)
}
