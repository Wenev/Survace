package in

import (
	"context"
)

type StreamService interface {
	GetStreamToken(ctx context.Context, userId string) (int32, string, string, string, error)
	GetLiveStreamId(ctx context.Context, userId string) (int32, string, string, bool, error)
	IsUserLive(ctx context.Context, userId string) (bool, error)
	GoLive(ctx context.Context, userId string, callId string) (int32, string, error)
	StopLive(ctx context.Context, userId string) (int32, string, error)
	ListLiveStreams(ctx context.Context) ([]struct{ UserId, CallId string }, error)
}
