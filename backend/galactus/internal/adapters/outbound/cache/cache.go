package cache

import (
	"encoding/json"
	"github.com/Wenev/Survace/galactus/config"
	"github.com/Wenev/Survace/galactus/internal/ports/out"
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

var _ out.CacheRepository = (*MemcachedConnection)(nil)

type MemcachedConnection struct {
	client *memcache.Client
}

func NewMemcachedConnection(serverAddress ...string) *MemcachedConnection {
	return &MemcachedConnection{
		client: memcache.New(serverAddress...),
	}
}

func CacheConnection() *MemcachedConnection {
	return NewMemcachedConnection(config.CacheAddressFromEnv())
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
