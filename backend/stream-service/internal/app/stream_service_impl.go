package app

import (
	"context"
	"fmt"
	"github.com/Wenev/Survace/stream-service/internal/adapters/outbound/cache"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"os"
	"time"

	"google.golang.org/grpc/codes"
)

type StreamServiceImpl struct {
	cache *cache.MemcachedConnection
}

func NewStreamService(cacheConn *cache.MemcachedConnection) *StreamServiceImpl {
	return &StreamServiceImpl{
		cache: cacheConn,
	}
}

func (s *StreamServiceImpl) GetStreamToken(ctx context.Context, userId string) (int32, string, string, string, error) {
	if userId == "" {
		return int32(codes.InvalidArgument), "userId is required", "", "", nil
	}

	apiKey := os.Getenv("STREAM_API_KEY")
	apiSecret := os.Getenv("STREAM_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		return int32(codes.Internal), "Stream API credentials not set", "", "", nil
	}
	fmt.Print(apiSecret)

	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(apiSecret))
	if err != nil {
		return int32(codes.Internal), "Failed to sign token", "", "", err
	}

	callId := fmt.Sprintf("%08d", rand.Intn(100000000))
	return int32(codes.OK), "success", signedToken, callId, nil
}

func (s *StreamServiceImpl) GetLiveStreamId(ctx context.Context, userId string) (int32, string, string, bool, error) {
	if userId == "" {
		return int32(codes.InvalidArgument), "userId is required", "", false, nil
	}
	key := fmt.Sprintf("live:%s", userId)
	var callId string
	err := s.cache.Get(key, &callId)
	if err == memcache.ErrCacheMiss {
		return int32(codes.NotFound), "user is not live", "", false, nil
	} else if err != nil {
		return int32(codes.Internal), "cache error", "", false, err
	}
	return int32(codes.OK), "success", callId, true, nil
}

func (s *StreamServiceImpl) IsUserLive(ctx context.Context, userId string) (bool, error) {
	if userId == "" {
		return false, nil
	}
	key := fmt.Sprintf("live:%s", userId)
	var callId string
	err := s.cache.Get(key, &callId)
	if err == memcache.ErrCacheMiss {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StreamServiceImpl) ListLiveStreams(ctx context.Context) ([]struct{ UserId, CallId string }, error) {
	var result []struct{ UserId, CallId string }
	var userIds []string
	err := s.cache.Get("live_users", &userIds)
	if err != nil && err != memcache.ErrCacheMiss {
		return nil, err
	}
	for _, userId := range userIds {
		var callId string
		key := fmt.Sprintf("live:%s", userId)
		err := s.cache.Get(key, &callId)
		if err == nil {
			result = append(result, struct{ UserId, CallId string }{UserId: userId, CallId: callId})
		}
	}
	return result, nil
}

func (s *StreamServiceImpl) GoLive(ctx context.Context, userId string, callId string) (int32, string, error) {
	if userId == "" || callId == "" {
		return int32(codes.InvalidArgument), "userId and callId are required", nil
	}
	key := fmt.Sprintf("live:%s", userId)
	err := s.cache.Set(key, callId, 48*time.Hour)
	if err != nil {
		return int32(codes.Internal), "failed to cache live status", err
	}
	var userIds []string
	err = s.cache.Get("live_users", &userIds)
	if err == memcache.ErrCacheMiss {
		userIds = []string{}
	} else if err != nil {
		return int32(codes.Internal), "failed to update live users", err
	}
	found := false
	for _, id := range userIds {
		if id == userId {
			found = true
			break
		}
	}
	if !found {
		userIds = append(userIds, userId)
		err = s.cache.Set("live_users", userIds, 48*time.Hour)
		if err != nil {
			return int32(codes.Internal), "failed to update live users", err
		}
	}
	return int32(codes.OK), "user is now live", nil
}

func (s *StreamServiceImpl) StopLive(ctx context.Context, userId string) (int32, string, error) {
	if userId == "" {
		return int32(codes.InvalidArgument), "userId is required", nil
	}
	key := fmt.Sprintf("live:%s", userId)
	err := s.cache.Delete(key)
	if err != nil && err != memcache.ErrCacheMiss {
		return int32(codes.Internal), "failed to delete live status", err
	}
	var userIds []string
	err = s.cache.Get("live_users", &userIds)
	if err == nil {
		newUserIds := make([]string, 0, len(userIds))
		for _, id := range userIds {
			if id != userId {
				newUserIds = append(newUserIds, id)
			}
		}
		s.cache.Set("live_users", newUserIds, 48*time.Hour)
	}
	return int32(codes.OK), "user is no longer live", nil
}
