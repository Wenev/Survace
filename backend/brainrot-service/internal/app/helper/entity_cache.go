package helper

import "fmt"

type CacheDeleter interface {
	Delete(key string) error
}

func InvalidateEntityCache(cache CacheDeleter, prefix string, id any) error {
	return cache.Delete(fmt.Sprintf("%s:%v", prefix, id))
}
