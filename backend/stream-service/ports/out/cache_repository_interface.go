package out

import (
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache: key not found")

type CacheRepository interface {
	Get(key string, dest interface{}) error
	Set(key string, value interface{}, expiration time.Duration) error
	Delete(key string) error
}
