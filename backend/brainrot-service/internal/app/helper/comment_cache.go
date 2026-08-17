package helper

import "fmt"

func InvalidateCommentCache(cache interface{ Delete(key string) error }, videoId int32) {
	videoKey := fmt.Sprintf("comments:video:%d", videoId)
	_ = cache.Delete(videoKey)
}
