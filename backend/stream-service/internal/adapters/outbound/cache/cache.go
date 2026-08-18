package cache

import (
	"encoding/json"
	"github.com/Wenev/Survace/stream-service/ports/out"
	"github.com/bradfitz/gomemcache/memcache"
	"log"
	"os"
	"time"
)

type MemcachedConnection struct {
	client *memcache.Client
}

var _ out.CacheRepository = (*MemcachedConnection)(nil)

func NewMemcachedConnection(serverAddress ...string) *MemcachedConnection {
	return &MemcachedConnection{
		client: memcache.New(serverAddress...),
	}
}

func CacheConnection() *MemcachedConnection {
	memcachedAddr := os.Getenv("CACHE_URL")
	log.Print(memcachedAddr + "hello")
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
	if err == memcache.ErrCacheMiss {
		return out.ErrCacheMiss
	}
	if err != nil {
		return err
	}
	log.Print(item.Value)
	log.Print("HELLO")
	return json.Unmarshal(item.Value, dest)
}

func (m *MemcachedConnection) Delete(key string) error {
	err := m.client.Delete(key)
	if err == memcache.ErrCacheMiss {
		return out.ErrCacheMiss
	}
	return err
}
