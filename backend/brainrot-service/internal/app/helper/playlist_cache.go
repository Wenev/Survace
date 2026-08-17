package helper

import "fmt"

func InvalidatePlaylistCache(cache interface{ Delete(key string) error }, userId int32) {
	userKey := fmt.Sprintf("playlist:user:%d", userId)
	_ = cache.Delete(userKey)
}
