package cache

import (
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"os"
	"time"
)

type MemcachedConnection struct {
	client *memcache.Client
}

func NewMemcachedConnection(serverAddress ...string) *MemcachedConnection {
	return &MemcachedConnection{
		client: memcache.New(serverAddress...),
	}
}

func CacheConnection() *MemcachedConnection {
	memcachedAddr := os.Getenv("CACHE_URL")
	if memcachedAddr == "" {
		memcachedAddr = "memcached:11212"
	}
	return NewMemcachedConnection(memcachedAddr)
}

func (m *MemcachedConnection) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return m.client.Set(&memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: int32(expiration.Seconds()),
	})
}

func (m *MemcachedConnection) Get(key string, dest interface{}) error {
	item, err := m.client.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(item.Value, dest)
}

func (m *MemcachedConnection) Delete(key string) error {
	return m.client.Delete(key)
}
