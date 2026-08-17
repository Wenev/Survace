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

type PlaylistServiceImpl struct {
	playlistRepo out.PlaylistRepository
	videoRepo    out.VideoRepository
	cache        *cache.MemcachedConnection
}

func NewPlaylistService(
	playlistRepo out.PlaylistRepository,
	_ out.PlaylistVideoRepository,
	videoRepo out.VideoRepository,
) *PlaylistServiceImpl {
	return &PlaylistServiceImpl{
		playlistRepo: playlistRepo,
		videoRepo:    videoRepo,
		cache:        cache.CacheConnection(),
	}
}

func (s *PlaylistServiceImpl) CreatePlaylist(ctx context.Context, userId int32, title string) (int32, string, *domain.Playlist, error) {
	if title == "" {
		return int32(codes.InvalidArgument), "Playlist title cannot be empty", nil, errors.New("playlist title cannot be empty")
	}

	playlist := &domain.Playlist{
		Title:  title,
		UserID: userId,
	}

	err := s.playlistRepo.Create(ctx, playlist)
	if err != nil {
		return int32(codes.Internal), "Failed to create playlist", nil, err
	}
	helper.InvalidatePlaylistCache(s.cache, userId)
	return int32(codes.OK), "Playlist created successfully", playlist, nil
}

func (s *PlaylistServiceImpl) DeletePlaylist(ctx context.Context, playlistId int32, userId int32) (int32, string, error) {
	playlist, err := s.playlistRepo.FindById(ctx, playlistId)
	if err != nil {
		return int32(codes.Internal), "Internal server error", err
	}
	if playlist == nil {
		return int32(codes.NotFound), "Playlist not found", errors.New("playlist not found")
	}

	if playlist.UserID != userId {
		return int32(codes.PermissionDenied), "You do not have permission to delete this playlist", errors.New("permission denied")
	}

	err = s.playlistRepo.Delete(ctx, playlistId)
	if err != nil {
		return int32(codes.Internal), "Failed to delete playlist", err
	}
	helper.InvalidatePlaylistCache(s.cache, userId)
	return int32(codes.OK), "Playlist deleted successfully", nil
}

func playlistCacheKey(userId int32) string {
	return fmt.Sprintf("playlist:user:%d", userId)
}

func (s *PlaylistServiceImpl) GetUserPlaylists(ctx context.Context, userId int32) (int32, string, []*domain.Playlist, error) {
	var cached []*domain.Playlist
	if err := s.cache.Get(playlistCacheKey(userId), &cached); err == nil && cached != nil {
		return int32(codes.OK), "Playlists fetched from cache", cached, nil
	}
	playlists, err := s.playlistRepo.FindByUserId(ctx, userId)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch playlists", nil, err
	}
	_ = s.cache.Set(playlistCacheKey(userId), playlists, 2*time.Minute)
	return int32(codes.OK), "Playlists fetched successfully", playlists, nil
}

func (s *PlaylistServiceImpl) GetPlaylistById(ctx context.Context, playlistId int32) (int32, string, *domain.Playlist, error) {
	playlist, err := s.playlistRepo.FindById(ctx, playlistId)
	if err != nil {
		return int32(codes.Internal), "Failed to fetch playlist", nil, err
	}
	if playlist == nil {
		return int32(codes.NotFound), "Playlist not found", nil, errors.New("playlist not found")
	}
	return int32(codes.OK), "Playlist retrieved successfully", playlist, nil
}

func (s *PlaylistServiceImpl) AddVideoToPlaylist(ctx context.Context, playlistId int32, videoId int32, order int32) (int32, string, error) {
	playlist, err := s.playlistRepo.FindById(ctx, playlistId)
	if err != nil {
		return int32(codes.Internal), "Internal server error", err
	}
	if playlist == nil {
		return int32(codes.NotFound), "Playlist not found", errors.New("playlist not found")
	}

	video, err := s.videoRepo.FindById(ctx, videoId)
	if err != nil {
		return int32(codes.Internal), "Internal server error", err
	}
	if video == nil {
		return int32(codes.NotFound), "Video not found", errors.New("video not found")
	}

	for _, pv := range playlist.Videos {
		if pv.VideoID == videoId {
			return int32(codes.AlreadyExists), "Video is already in the playlist", errors.New("video already in playlist")
		}
	}

	if order == 0 {
		order = int32(len(playlist.Videos)) + 1
	}

	playlistVideo := domain.PlaylistVideo{
		PlaylistID: playlistId,
		VideoID:    videoId,
		Order:      order,
		Video:      *video,
	}

	playlist.Videos = append(playlist.Videos, playlistVideo)
	err = s.playlistRepo.Update(ctx, playlist)
	if err != nil {
		return int32(codes.Internal), "Failed to add video to playlist", err
	}

	return int32(codes.OK), "Video added to playlist successfully", nil
}

func (s *PlaylistServiceImpl) RemoveVideoFromPlaylist(ctx context.Context, playlistId int32, videoId int32) (int32, string, error) {
	playlist, err := s.playlistRepo.FindById(ctx, playlistId)
	if err != nil {
		return int32(codes.Internal), "Internal server error", err
	}
	if playlist == nil {
		return int32(codes.NotFound), "Playlist not found", errors.New("playlist not found")
	}

	found := false
	newVideos := make([]domain.PlaylistVideo, 0, len(playlist.Videos))
	for _, pv := range playlist.Videos {
		if pv.VideoID == videoId {
			found = true
			continue
		}
		newVideos = append(newVideos, pv)
	}
	if !found {
		return int32(codes.NotFound), "Video not found in playlist", errors.New("video not in playlist")
	}
	playlist.Videos = newVideos
	err = s.playlistRepo.Update(ctx, playlist)
	if err != nil {
		return int32(codes.Internal), "Failed to remove video from playlist", err
	}
	return int32(codes.OK), "Video removed from playlist successfully", nil
}

func (s *PlaylistServiceImpl) ReorderVideo(ctx context.Context, playlistId int32, videoId int32, newOrder int32) (int32, string, error) {
	playlist, err := s.playlistRepo.FindById(ctx, playlistId)
	if err != nil {
		return int32(codes.Internal), "Internal server error", err
	}
	if playlist == nil {
		return int32(codes.NotFound), "Playlist not found", errors.New("playlist not found")
	}
	if newOrder < 0 || int(newOrder) >= len(playlist.Videos) {
		return int32(codes.InvalidArgument), "Invalid order", errors.New("invalid order")
	}
	var idx int
	found := false
	for i, pv := range playlist.Videos {
		if pv.VideoID == videoId {
			idx = i
			found = true
			break
		}
	}
	if !found {
		return int32(codes.NotFound), "Video not found in playlist", errors.New("video not in playlist")
	}
	video := playlist.Videos[idx]
	playlist.Videos = append(playlist.Videos[:idx], playlist.Videos[idx+1:]...)
	if int(newOrder) >= len(playlist.Videos) {
		playlist.Videos = append(playlist.Videos, video)
	} else {
		playlist.Videos = append(playlist.Videos[:newOrder], append([]domain.PlaylistVideo{video}, playlist.Videos[newOrder:]...)...)
	}
	for i := range playlist.Videos {
		playlist.Videos[i].Order = int32(i + 1)
	}
	err = s.playlistRepo.Update(ctx, playlist)
	if err != nil {
		return int32(codes.Internal), "Failed to reorder video", err
	}
	helper.InvalidatePlaylistCache(s.cache, playlist.UserID)
	return int32(codes.OK), "Video reordered successfully", nil
}

func (s *PlaylistServiceImpl) UpdatePlaylistTitle(ctx context.Context, playlistId int32, title string) (int32, string, error) {
	if title == "" {
		return int32(3), "Playlist title cannot be empty", errors.New("playlist title cannot be empty")
	}
	playlist, err := s.playlistRepo.FindById(ctx, playlistId)
	if err != nil {
		return int32(13), "Internal server error", err
	}
	if playlist == nil {
		return int32(5), "Playlist not found", errors.New("playlist not found")
	}
	playlist.Title = title
	err = s.playlistRepo.Update(ctx, playlist)
	if err != nil {
		return int32(13), "Failed to update playlist title", err
	}
	helper.InvalidatePlaylistCache(s.cache, playlist.UserID)
	return int32(0), "Playlist title updated successfully", nil
}
