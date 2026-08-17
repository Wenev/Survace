package helper

import "fmt"

// Invalidate all cache keys for a comment's replies
func InvalidateReplyCache(cache interface{ Delete(key string) error }, commentId int32) {
	commentKey := fmt.Sprintf("replies:comment:%d", commentId)
	_ = cache.Delete(commentKey)
}
