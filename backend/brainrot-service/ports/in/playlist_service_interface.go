package in

import (
	"context"
	"github.com/Wenev/Survace/brainrot-service/internal/app/domain"
)

type PlaylistService interface {
	CreatePlaylist(ctx context.Context, userId int32, title string) (int32, string, *domain.Playlist, error)
	DeletePlaylist(ctx context.Context, playlistId int32, userId int32) (int32, string, error)
	AddVideoToPlaylist(ctx context.Context, playlistId int32, videoId int32, order int32) (int32, string, error)
	RemoveVideoFromPlaylist(ctx context.Context, playlistId int32, videoId int32) (int32, string, error)
	ReorderVideo(ctx context.Context, playlistId int32, videoId int32, newOrder int32) (int32, string, error)
	GetUserPlaylists(ctx context.Context, userId int32) (int32, string, []*domain.Playlist, error)
	GetPlaylistById(ctx context.Context, playlistId int32) (int32, string, *domain.Playlist, error)
	UpdatePlaylistTitle(ctx context.Context, playlistId int32, title string) (int32, string, error)
}
