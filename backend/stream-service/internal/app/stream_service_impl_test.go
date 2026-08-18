package app

import (
	"context"
	"testing"
	"time"

	"github.com/Wenev/Survace/stream-service/ports/out"
)

type fakeCache struct {
	store map[string]interface{}
}

func newFakeCache() *fakeCache {
	return &fakeCache{store: make(map[string]interface{})}
}

func (f *fakeCache) Get(key string, dest interface{}) error {
	v, ok := f.store[key]
	if !ok {
		return out.ErrCacheMiss
	}
	switch d := dest.(type) {
	case *string:
		*d = v.(string)
	case *[]string:
		*d = v.([]string)
	}
	return nil
}

func (f *fakeCache) Set(key string, value interface{}, expiration time.Duration) error {
	f.store[key] = value
	return nil
}

func (f *fakeCache) Delete(key string) error {
	if _, ok := f.store[key]; !ok {
		return out.ErrCacheMiss
	}
	delete(f.store, key)
	return nil
}

func TestIsUserLive_CacheMiss(t *testing.T) {
	s := NewStreamService(newFakeCache())
	live, err := s.IsUserLive(context.Background(), "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if live {
		t.Fatalf("expected user to not be live on cache miss")
	}
}

func TestIsUserLive_CacheHit(t *testing.T) {
	cache := newFakeCache()
	cache.store["live:user1"] = "call123"
	s := NewStreamService(cache)
	live, err := s.IsUserLive(context.Background(), "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !live {
		t.Fatalf("expected user to be live on cache hit")
	}
}

func TestIsUserLive_EmptyUserId(t *testing.T) {
	s := NewStreamService(newFakeCache())
	live, err := s.IsUserLive(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if live {
		t.Fatalf("expected false for empty userId")
	}
}

func TestGoLive_SetsKeyAndTracksUser(t *testing.T) {
	cache := newFakeCache()
	s := NewStreamService(cache)
	code, _, err := s.GoLive(context.Background(), "user1", "call123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected OK code, got %d", code)
	}
	if cache.store["live:user1"] != "call123" {
		t.Fatalf("expected live key to be set, got %v", cache.store["live:user1"])
	}
	userIds := cache.store["live_users"].([]string)
	if len(userIds) != 1 || userIds[0] != "user1" {
		t.Fatalf("expected live_users to contain user1, got %v", userIds)
	}
}

func TestGoLive_DoesNotDuplicateUser(t *testing.T) {
	cache := newFakeCache()
	s := NewStreamService(cache)
	if _, _, err := s.GoLive(context.Background(), "user1", "call123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, _, err := s.GoLive(context.Background(), "user1", "call456"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	userIds := cache.store["live_users"].([]string)
	if len(userIds) != 1 {
		t.Fatalf("expected live_users to contain exactly one entry, got %v", userIds)
	}
}

func TestGoLive_MissingArgs(t *testing.T) {
	s := NewStreamService(newFakeCache())
	code, _, err := s.GoLive(context.Background(), "", "call123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code == 0 {
		t.Fatalf("expected non-OK code for missing userId")
	}
}

func TestStopLive_RemovesKeyAndUser(t *testing.T) {
	cache := newFakeCache()
	s := NewStreamService(cache)
	if _, _, err := s.GoLive(context.Background(), "user1", "call123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code, _, err := s.StopLive(context.Background(), "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected OK code, got %d", code)
	}
	if _, ok := cache.store["live:user1"]; ok {
		t.Fatalf("expected live key to be removed")
	}
	userIds := cache.store["live_users"].([]string)
	if len(userIds) != 0 {
		t.Fatalf("expected live_users to be empty, got %v", userIds)
	}
}

func TestStopLive_CacheMissIsNotError(t *testing.T) {
	s := NewStreamService(newFakeCache())
	code, _, err := s.StopLive(context.Background(), "user1")
	if err != nil {
		t.Fatalf("expected cache miss on delete to be tolerated, got error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected OK code, got %d", code)
	}
}

func TestStopLive_MissingUserId(t *testing.T) {
	s := NewStreamService(newFakeCache())
	code, _, err := s.StopLive(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code == 0 {
		t.Fatalf("expected non-OK code for missing userId")
	}
}
