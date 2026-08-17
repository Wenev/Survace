package helper

import "fmt"

func InvalidateReplyCache(cache interface{ Delete(key string) error }, commentId int32) {
	commentKey := fmt.Sprintf("replies:comment:%d", commentId)
	_ = cache.Delete(commentKey)
}
