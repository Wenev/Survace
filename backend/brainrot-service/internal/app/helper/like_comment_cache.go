package helper

import "fmt"

// Invalidate all cache keys for a comment/user
func InvalidateLikeCommentCache(cache interface{ Delete(key string) error }, userId int32, commentId int32) {
	userKey := fmt.Sprintf("likecomments:user:%d", userId)
	commentKey := fmt.Sprintf("likecomments:comment:%d", commentId)
	_ = cache.Delete(userKey)
	_ = cache.Delete(commentKey)
}
