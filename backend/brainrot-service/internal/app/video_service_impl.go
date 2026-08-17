package app

import (
	"context"
	"fmt"
	"github.com/Wenev/Survace/brainrot-service/internal/adapters/inbound/grpc"
	"github.com/Wenev/Survace/brainrot-service/internal/adapters/outbound/cache"
	"github.com/Wenev/Survace/brainrot-service/internal/app/domain"
	"github.com/Wenev/Survace/brainrot-service/internal/app/helper"
	"github.com/Wenev/Survace/brainrot-service/ports/out"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type VideoServiceImpl struct {
	videoRepo         out.VideoRepository
	likeRepo          out.LikeRepository
	commentRepo       out.CommentRepository
	watchHistoryRepo  out.WatchHistoryRepository
	playlistRepo      out.PlaylistRepository
	playlistVideoRepo out.PlaylistVideoRepository
	cache             *cache.MemcachedConnection
	SocialClient      *grpc.SocialClient
}

func NewVideoService(
	videoRepo out.VideoRepository,
	likeRepo out.LikeRepository,
	commentRepo out.CommentRepository,
	watchHistoryRepo out.WatchHistoryRepository,
	playlistRepo out.PlaylistRepository,
	playlistVideoRepo out.PlaylistVideoRepository,
	socialClient *grpc.SocialClient) *VideoServiceImpl {
	return &VideoServiceImpl{
		videoRepo:         videoRepo,
		likeRepo:          likeRepo,
		commentRepo:       commentRepo,
		watchHistoryRepo:  watchHistoryRepo,
		playlistRepo:      playlistRepo,
		playlistVideoRepo: playlistVideoRepo,
		cache:             cache.CacheConnection(),
		SocialClient:      socialClient,
	}
}

func invalidateSearchCache(c *cache.MemcachedConnection) {
	const versionKey = "search:video:version"
	var version int64
	if err := c.Get(versionKey, &version); err == nil {
		_ = c.Set(versionKey, version+1, 0)
	} else {
		_ = c.Set(versionKey, int64(1), 0)
	}
}

func (v *VideoServiceImpl) UploadVideo(ctx context.Context, title string, file []byte, userId int32, enableComment bool, visibility string, thumbnail []byte, isDraft bool, postDate *time.Time) (int32, string, *domain.Video, error) {
	if visibility != "public" && visibility != "unlisted" && visibility != "private" {
		return 3, "Invalid visibility value", nil, status.Errorf(codes.InvalidArgument, "Invalid visibility value")
	}
	createdTime := time.Now()
	objectName := fmt.Sprintf("%d_%d.mp4", userId, createdTime.Unix())
	contentType := "video/mp4"
	url, err := v.videoRepo.UploadMinIO(ctx, objectName, file, contentType)
	if err != nil {
		return 13, "Internal server error", nil, status.Errorf(codes.Internal, "Internal server error")
	}

	var thumbnailUrl string
	if len(thumbnail) > 0 {
		thumbObjectName := fmt.Sprintf("%d_%d_thumb.jpg", userId, createdTime.Unix())
		thumbContentType := "image/jpeg"
		thumbnailUrl, err = v.videoRepo.UploadMinIO(ctx, thumbObjectName, thumbnail, thumbContentType)
		if err != nil {
			return 13, "Failed to upload thumbnail", nil, status.Errorf(codes.Internal, "Failed to upload thumbnail")
		}
	}

	postedAt := createdTime
	if postDate != nil {
		postedAt = *postDate
	}

	video := &domain.Video{
		UserID:        userId,
		ObjectName:    objectName,
		URL:           url,
		Title:         title,
		EnableComment: enableComment,
		Visibility:    visibility,
		CreatedAt:     createdTime,
		ThumbnailUrl:  thumbnailUrl,
		PostedAt:      postedAt,
		IsDraft:       isDraft,
	}

	err = v.videoRepo.Create(ctx, video)
	if err != nil {
		return 13, "Internal server error", nil, status.Errorf(codes.Internal, "Internal server error")
	}
	helper.InvalidateFeedCache(v.cache, userId)
	invalidateSearchCache(v.cache)
	return 0, "Successfully Uploaded Video", video, nil
}

func (v *VideoServiceImpl) GetVideoByID(ctx context.Context, id int32) (int32, string, *domain.Video, int64, int64, int32, error) {
	video, err := v.videoRepo.FindById(ctx, id)
	if err != nil {
		return 13, "Internal server error", nil, 0, 0, 0, status.Errorf(codes.Internal, "Internal server error")
	}
	if video == nil {
		return 5, "Video Not Found", nil, 0, 0, 0, status.Errorf(codes.NotFound, "Video Not Found")
	}
	likeCacheKey := fmt.Sprintf("like:video:%d", video.ID)
	var likeCount int64
	if err := v.cache.Get(likeCacheKey, &likeCount); err == nil {
	} else {
		likeCount, err = v.likeRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return 13, "Failed to count likes", nil, 0, 0, 0, err
		}
		_ = v.cache.Set(likeCacheKey, likeCount, 2*time.Minute)
	}
	commentCount, err := v.commentRepo.CountByVideoId(ctx, video.ID)
	if err != nil {
		return 13, "Failed to count comments", nil, 0, 0, 0, err
	}
	viewCount64, err := v.watchHistoryRepo.CountByVideoId(ctx, video.ID)
	if err != nil {
		return 13, "Failed to count views", nil, 0, 0, 0, err
	}
	return 0, "Successfuly Get Video By ID", video, likeCount, commentCount, int32(viewCount64), nil
}

func (v *VideoServiceImpl) GetUserVideo(ctx context.Context, userId int32) (int32, string, []*domain.Video, []int64, []int64, []int32, error) {
	videos, err := v.videoRepo.FindByUserId(ctx, userId)
	if err != nil {
		return 13, "Internal server error", nil, nil, nil, nil, status.Errorf(codes.Internal, "Internal server error")
	}
	likeCounts := make([]int64, len(videos))
	commentCounts := make([]int64, len(videos))
	viewCounts := make([]int32, len(videos))
	for i, video := range videos {
		likeCacheKey := fmt.Sprintf("like:video:%d", video.ID)
		var likeCount int64
		if err := v.cache.Get(likeCacheKey, &likeCount); err == nil {
		} else {
			likeCount, err = v.likeRepo.CountByVideoId(ctx, video.ID)
			if err != nil {
				return 13, "Failed to count likes", nil, nil, nil, nil, err
			}
			_ = v.cache.Set(likeCacheKey, likeCount, 2*time.Minute)
		}
		commentCount, err := v.commentRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return 13, "Failed to count comments", nil, nil, nil, nil, err
		}
		viewCount64, err := v.watchHistoryRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return 13, "Failed to count views", nil, nil, nil, nil, err
		}
		likeCounts[i] = likeCount
		commentCounts[i] = commentCount
		viewCounts[i] = int32(viewCount64)
	}
	return 0, "Successfully Get User Video", videos, likeCounts, commentCounts, viewCounts, nil
}

func (v *VideoServiceImpl) GetUserVideoAndDraft(ctx context.Context, userId int32) (int32, string, []*domain.Video, []int64, []int64, []int32, error) {
	videos, err := v.videoRepo.FindByUserIdAndAll(ctx, userId)
	if err != nil {
		return 13, "Internal server error", nil, nil, nil, nil, err
	}
	likeCounts := make([]int64, len(videos))
	commentCounts := make([]int64, len(videos))
	viewCounts := make([]int32, len(videos))
	for i, video := range videos {
		likeCacheKey := fmt.Sprintf("like:video:%d", video.ID)
		var likeCount int64
		if err := v.cache.Get(likeCacheKey, &likeCount); err == nil {
		} else {
			likeCount, err = v.likeRepo.CountByVideoId(ctx, video.ID)
			if err != nil {
				return 13, "Failed to count likes", nil, nil, nil, nil, err
			}
			_ = v.cache.Set(likeCacheKey, likeCount, 2*time.Minute)
		}
		commentCount, err := v.commentRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return 13, "Failed to count comments", nil, nil, nil, nil, err
		}
		viewCount64, err := v.watchHistoryRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return 13, "Failed to count views", nil, nil, nil, nil, err
		}
		likeCounts[i] = likeCount
		commentCounts[i] = commentCount
		viewCounts[i] = int32(viewCount64)
	}
	return 0, "Successfully Get User Video And Draft", videos, likeCounts, commentCounts, viewCounts, nil
}

func (v *VideoServiceImpl) DeleteUserVideo(ctx context.Context, userId int32) (int32, string, error) {
	err := v.videoRepo.DeleteByUserId(ctx, userId)
	if err != nil {
		return 13, "Internal server error", status.Errorf(codes.Internal, "Internal server error")
	}

	helper.InvalidateFeedCache(v.cache, userId)
	invalidateSearchCache(v.cache)
	return 0, "Successfully Delete User Video", nil
}

func (v *VideoServiceImpl) GetRandomFeed(ctx context.Context, userId *int32, limit int, offset int) (int32, string, []*domain.Video, []int64, []int64, []int32, error) {
	newPopular, err := v.videoRepo.FindNewAndPopular(ctx, limit*2, offset)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch new/popular videos", nil, nil, nil, nil, err
	}

	var unwatched []*domain.Video
	if userId != nil {
		watchHistories, err := v.watchHistoryRepo.FindByUserId(ctx, *userId)
		if err != nil {
			return int32(codes.Internal), "Failed to fetch watched videos", nil, nil, nil, nil, err
		}
		watchedSet := make(map[int32]struct{}, len(watchHistories))
		for _, wh := range watchHistories {
			watchedSet[wh.VideoID] = struct{}{}
		}
		unwatched = helper.FilterUnwatched(newPopular, watchedSet)
	} else {
		unwatched = newPopular
	}

	randoms, err := v.videoRepo.FindRandomVideos(ctx, limit, offset)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch random videos", nil, nil, nil, nil, err
	}

	feed := helper.WeightedShuffle(unwatched, randoms, 0.7, limit)
	feed = helper.DeduplicateVideos(feed)
	if feed == nil {
		feed = []*domain.Video{}
	}
	if len(feed) == 0 {
		return int32(codes.Internal), "mamasita", []*domain.Video{}, []int64{}, []int64{}, []int32{}, nil
	}
	likeCounts := make([]int64, len(feed))
	commentCounts := make([]int64, len(feed))
	viewCounts := make([]int32, len(feed))
	for i, video := range feed {
		likeCount, err := v.likeRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return int32(codes.Internal), "Failed to count likes", nil, nil, nil, nil, err
		}
		commentCount, err := v.commentRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return int32(codes.Internal), "Failed to count comments", nil, nil, nil, nil, err
		}
		viewCount64, err := v.watchHistoryRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return int32(codes.Internal), "Failed to count views", nil, nil, nil, nil, err
		}
		if likeCount != 0 {
			likeCounts[i] = likeCount
		}
		if commentCount != 0 {
			commentCounts[i] = commentCount
		}
		if viewCount64 != 0 {
			viewCounts[i] = int32(viewCount64)
		}
	}

	return int32(codes.OK), "Feed fetched successfully", feed, likeCounts, commentCounts, viewCounts, nil
}

func (v *VideoServiceImpl) GetRandomFeedLoggedOut(ctx context.Context, limit int, offset int) (int32, string, []*domain.Video, []int64, []int64, []int32, error) {
	newPopular, err := v.videoRepo.FindNewAndPopular(ctx, limit*2, offset)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch new/popular videos", nil, nil, nil, nil, err
	}

	randoms, err := v.videoRepo.FindRandomVideos(ctx, limit, offset)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch random videos", nil, nil, nil, nil, err
	}

	feed := helper.WeightedShuffle(newPopular, randoms, 0.7, limit)
	feed = helper.DeduplicateVideos(feed)
	fmt.Print(feed)
	if feed == nil {
		feed = []*domain.Video{}
	}
	if len(feed) == 0 {
		return int32(codes.Internal), "mamsita", []*domain.Video{}, []int64{}, []int64{}, []int32{}, nil
	}
	likeCounts := make([]int64, len(feed))
	commentCounts := make([]int64, len(feed))
	viewCounts := make([]int32, len(feed))
	for i, video := range feed {
		likeCount, err := v.likeRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return int32(codes.Internal), "Failed to count likes", nil, nil, nil, nil, err
		}
		commentCount, err := v.commentRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return int32(codes.Internal), "Failed to count comments", nil, nil, nil, nil, err
		}
		viewCount64, err := v.watchHistoryRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return int32(codes.Internal), "Failed to count views", nil, nil, nil, nil, err
		}
		if likeCount != 0 {
			likeCounts[i] = likeCount
		}
		if commentCount != 0 {
			commentCounts[i] = commentCount
		}
		if viewCount64 != 0 {
			viewCounts[i] = int32(viewCount64)
		}
	}

	return int32(codes.OK), "Feed fetched successfully", feed, likeCounts, commentCounts, viewCounts, nil
}

func (v *VideoServiceImpl) WatchVideo(ctx context.Context, userId int32, videoId int32) (int32, string, error) {
	if userId <= 0 || videoId <= 0 {
		return int32(codes.InvalidArgument), "Invalid user or video ID", fmt.Errorf("invalid user or video ID")
	}
	video, err := v.videoRepo.FindById(ctx, videoId)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch video", err
	}
	if video == nil {
		return int32(codes.NotFound), "Video not found", fmt.Errorf("video not found")
	}
	video.ViewCount++
	err = v.videoRepo.Update(ctx, video)
	if err != nil {
		return int32(codes.Internal), "Failed to update view count", err
	}
	watchHistory := &domain.WatchHistory{
		UserID:    userId,
		VideoID:   videoId,
		CreatedAt: time.Now(),
	}
	err = v.watchHistoryRepo.Create(ctx, watchHistory)
	if err != nil {
		return int32(codes.Internal), "Failed to mark video as watched", err
	}
	return int32(codes.OK), "Video marked as watched", nil
}

func (v *VideoServiceImpl) SearchVideo(ctx context.Context, query string, threshold float64, limit int, offset int) (int32, string, []*domain.Video, []int64, []int64, []int32, error) {
	if query == "" {
		return int32(codes.InvalidArgument), "Search query cannot be empty", nil, nil, nil, nil,
			fmt.Errorf("empty search query")
	}

	allVideos, err := v.videoRepo.FindNewAndPopular(ctx, 1000, 1)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch videos", nil, nil, nil, nil, err
	}
	log.Print(allVideos)
	log.Print("save it")
	filteredVideos := helper.FilterVideosByJaroDistance(allVideos, query, threshold)
	start := offset
	end := offset + limit
	if start >= len(filteredVideos) {
		return int32(codes.OK), "No matching videos found", []*domain.Video{}, []int64{}, []int64{}, []int32{}, nil
	}
	if end > len(filteredVideos) {
		end = len(filteredVideos)
	}
	pagedVideos := filteredVideos[start:end]

	likeCounts := make([]int64, len(pagedVideos))
	commentCounts := make([]int64, len(pagedVideos))
	viewCounts := make([]int32, len(pagedVideos))

	for i, video := range pagedVideos {
		likeCount, err := v.likeRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return int32(codes.Internal), "Failed to count likes", nil, nil, nil, nil, err
		}
		commentCount, err := v.commentRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return int32(codes.Internal), "Failed to count comments", nil, nil, nil, nil, err
		}

		viewCount64, err := v.watchHistoryRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return int32(codes.Internal), "Failed to count views", nil, nil, nil, nil, err
		}

		likeCounts[i] = likeCount
		commentCounts[i] = commentCount
		viewCounts[i] = int32(viewCount64)
	}

	return int32(codes.OK), "Videos found successfully", pagedVideos, likeCounts, commentCounts, viewCounts, nil
}

func (v *VideoServiceImpl) DeleteVideoByID(ctx context.Context, videoId int32) (int32, string, error) {

	if err := v.playlistVideoRepo.DeleteByVideoId(ctx, videoId); err != nil {
		return 13, "Failed to delete video from playlists", err
	}

	if err := v.likeRepo.DeleteByVideoId(ctx, videoId); err != nil {
		return 13, "Failed to delete likes", err
	}

	if err := v.commentRepo.DeleteByVideoId(ctx, videoId); err != nil {
		return 13, "Failed to delete comments", err
	}

	if err := v.watchHistoryRepo.DeleteByVideoId(ctx, videoId); err != nil {
		return 13, "Failed to delete watch history", err
	}

	if err := v.videoRepo.Delete(ctx, videoId); err != nil {
		return 13, "Internal server error", err
	}

	helper.InvalidateFeedCache(v.cache, 0)
	invalidateSearchCache(v.cache)

	return 0, "Successfully deleted video and all related data", nil
}

func (v *VideoServiceImpl) UpdateVideo(ctx context.Context, id int32, userId int32, title string, url string, objectName string, enableComment bool, visibility string, thumbnail []byte, postedAt *time.Time, isDraft bool) (int32, string, *domain.Video, error) {
	video, err := v.videoRepo.FindById(ctx, id)
	if err != nil {
		return 13, "Internal server error", nil, err
	}
	if video == nil {
		return 5, "Video Not Found", nil, nil
	}
	video.UserID = userId
	video.Title = title
	video.URL = url
	video.ObjectName = objectName
	video.EnableComment = enableComment
	video.Visibility = visibility
	if len(thumbnail) > 0 {
		thumbObjectName := fmt.Sprintf("%d_%d_thumb.jpg", userId, time.Now().Unix())
		thumbContentType := "image/jpeg"
		thumbnailUrl, err := v.videoRepo.UploadMinIO(ctx, thumbObjectName, thumbnail, thumbContentType)
		if err != nil {
			return 13, "Failed to upload thumbnail", nil, err
		}
		video.ThumbnailUrl = thumbnailUrl
	}
	if postedAt != nil {
		video.PostedAt = *postedAt
	}
	video.IsDraft = isDraft

	err = v.videoRepo.Update(ctx, video)
	if err != nil {
		return 13, "Internal server error", nil, err
	}
	return 0, "Successfully Updated Video", video, nil
}

func (v *VideoServiceImpl) FriendVideo(ctx context.Context, userId int32, limit int, offset int) (int32, string, []*domain.Video, []int64, []int64, []int32, error) {
	cacheKey := fmt.Sprintf("friend:video:%d:%d:%d", userId, limit, offset)
	var cached struct {
		Feed          []*domain.Video
		LikeCounts    []int64
		CommentCounts []int64
		ViewCounts    []int32
	}
	if err := v.cache.Get(cacheKey, &cached); err == nil && cached.Feed != nil {
		return 0, "Successfully fetched friends' videos (from cache)", cached.Feed, cached.LikeCounts, cached.CommentCounts, cached.ViewCounts, nil
	}

	resp, err := v.SocialClient.ListFriend(ctx, userId)
	if err != nil {
		return 13, "Failed to fetch friends", nil, nil, nil, nil, err
	}
	var allVideos []*domain.Video
	for _, friend := range resp.Data {
		userVideos, err := v.videoRepo.FindByUserId(ctx, friend.UserId)
		if err != nil {
			return 13, "Failed to fetch friend's videos", nil, nil, nil, nil, err
		}
		allVideos = append(allVideos, userVideos...)
	}
	watchHistories, err := v.watchHistoryRepo.FindByUserId(ctx, userId)
	if err != nil {
		return 13, "Failed to fetch watched videos", nil, nil, nil, nil, err
	}
	watchedSet := make(map[int32]struct{}, len(watchHistories))
	for _, wh := range watchHistories {
		watchedSet[wh.VideoID] = struct{}{}
	}
	unwatched := helper.FilterUnwatched(allVideos, watchedSet)
	randoms, err := v.videoRepo.FindRandomVideos(ctx, limit, offset)
	if err != nil {
		return 13, "Failed to fetch random videos", nil, nil, nil, nil, err
	}
	feed := helper.WeightedShuffle(unwatched, randoms, 0.7, limit)
	feed = helper.DeduplicateVideos(feed)
	if feed == nil {
		feed = []*domain.Video{}
	}
	likeCounts := make([]int64, len(feed))
	commentCounts := make([]int64, len(feed))
	viewCounts := make([]int32, len(feed))
	for i, video := range feed {
		likeCount, err := v.likeRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return 13, "Failed to count likes", nil, nil, nil, nil, err
		}
		commentCount, err := v.commentRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return 13, "Failed to count comments", nil, nil, nil, nil, err
		}
		viewCount64, err := v.watchHistoryRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return 13, "Failed to count views", nil, nil, nil, nil, err
		}
		likeCounts[i] = likeCount
		commentCounts[i] = commentCount
		viewCounts[i] = int32(viewCount64)
	}
	cacheValue := struct {
		Feed          []*domain.Video
		LikeCounts    []int64
		CommentCounts []int64
		ViewCounts    []int32
	}{
		Feed:          feed,
		LikeCounts:    likeCounts,
		CommentCounts: commentCounts,
		ViewCounts:    viewCounts,
	}
	_ = v.cache.Set(cacheKey, cacheValue, 2*time.Minute)
	return 0, "Successfully fetched friends' videos", feed, likeCounts, commentCounts, viewCounts, nil
}

func (v *VideoServiceImpl) FollowingVideo(ctx context.Context, userId int32, limit int, offset int) (int32, string, []*domain.Video, []int64, []int64, []int32, error) {
	cacheKey := fmt.Sprintf("following:video:%d:%d:%d", userId, limit, offset)
	var cached struct {
		Feed          []*domain.Video
		LikeCounts    []int64
		CommentCounts []int64
		ViewCounts    []int32
	}
	if err := v.cache.Get(cacheKey, &cached); err == nil && cached.Feed != nil {
		return 0, "Successfully fetched following's videos (from cache)", cached.Feed, cached.LikeCounts, cached.CommentCounts, cached.ViewCounts, nil
	}

	resp, err := v.SocialClient.ListFollowing(ctx, userId)
	if err != nil {
		fmt.Print("TIDAK BISA")
		return 13, "Failed to fetch following", nil, nil, nil, nil, err
	}
	var allVideos []*domain.Video
	for _, following := range resp.Data {
		userVideos, err := v.videoRepo.FindByUserId(ctx, following.FolloweeId)
		if err != nil {
			return 13, "Failed to fetch following's videos", nil, nil, nil, nil, err
		}
		allVideos = append(allVideos, userVideos...)
	}
	watchHistories, err := v.watchHistoryRepo.FindByUserId(ctx, userId)
	if err != nil {
		return 13, "Failed to fetch watched videos", nil, nil, nil, nil, err
	}
	watchedSet := make(map[int32]struct{}, len(watchHistories))
	for _, wh := range watchHistories {
		watchedSet[wh.VideoID] = struct{}{}
	}
	unwatched := helper.FilterUnwatched(allVideos, watchedSet)
	randoms, err := v.videoRepo.FindRandomVideos(ctx, limit, offset)
	if err != nil {
		return 13, "Failed to fetch random videos", nil, nil, nil, nil, err
	}
	feed := helper.WeightedShuffle(unwatched, randoms, 0.7, limit)
	feed = helper.DeduplicateVideos(feed)
	if feed == nil {
		feed = []*domain.Video{}
	}
	likeCounts := make([]int64, len(feed))
	commentCounts := make([]int64, len(feed))
	viewCounts := make([]int32, len(feed))
	for i, video := range feed {
		likeCount, err := v.likeRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return 13, "Failed to count likes", nil, nil, nil, nil, err
		}
		commentCount, err := v.commentRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return 13, "Failed to count comments", nil, nil, nil, nil, err
		}
		viewCount64, err := v.watchHistoryRepo.CountByVideoId(ctx, video.ID)
		if err != nil {
			return 13, "Failed to count views", nil, nil, nil, nil, err
		}
		likeCounts[i] = likeCount
		commentCounts[i] = commentCount
		viewCounts[i] = int32(viewCount64)
	}
	cacheValue := struct {
		Feed          []*domain.Video
		LikeCounts    []int64
		CommentCounts []int64
		ViewCounts    []int32
	}{
		Feed:          feed,
		LikeCounts:    likeCounts,
		CommentCounts: commentCounts,
		ViewCounts:    viewCounts,
	}
	_ = v.cache.Set(cacheKey, cacheValue, 2*time.Minute)
	return 0, "Successfully fetched following's videos", feed, likeCounts, commentCounts, viewCounts, nil
}
