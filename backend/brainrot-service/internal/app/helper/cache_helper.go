package helper

import (
	"fmt"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/adapters/outbound/cache"
)

func InvalidateFeedCache(c *cache.MemcachedConnection, userIds ...int32) {
	limits := []int{10, 20, 50}
	maxOffset := 100
	for _, limit := range limits {
		for offset := 0; offset < maxOffset; offset += limit {
			// Invalidate for logged-out
			c.Delete(feedCacheKey(nil, limit, offset))
			// Invalidate for each user if provided
			for _, userId := range userIds {
				c.Delete(feedCacheKey(&userId, limit, offset))
			}
		}
	}
}

func feedCacheKey(userId *int32, limit, offset int) string {
	if userId == nil {
		return "feed:loggedout:" + itoa(limit) + ":" + itoa(offset)
	}
	return "feed:user:" + itoa(int(*userId)) + ":" + itoa(limit) + ":" + itoa(offset)
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
