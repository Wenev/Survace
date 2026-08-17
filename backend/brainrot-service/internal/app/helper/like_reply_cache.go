package helper

import "fmt"

func InvalidateLikeReplyCache(cache interface{ Delete(key string) error }, userId int32, replyId int32) {
	userKey := fmt.Sprintf("likereplies:user:%d", userId)
	replyKey := fmt.Sprintf("likereplies:reply:%d", replyId)
	_ = cache.Delete(userKey)
	_ = cache.Delete(replyKey)
}
