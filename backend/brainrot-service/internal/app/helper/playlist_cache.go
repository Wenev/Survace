package helper

import "fmt"

// Invalidate all cache keys for a playlist/user
func InvalidatePlaylistCache(cache interface{ Delete(key string) error }, userId int32) {
	userKey := fmt.Sprintf("playlist:user:%d", userId)
	_ = cache.Delete(userKey)
}
