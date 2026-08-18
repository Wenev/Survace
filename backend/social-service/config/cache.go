package config

import "os"

func CacheAddressFromEnv() string {
	addr := os.Getenv("CACHE_URL")
	if addr == "" {
		addr = "memcached:11212"
	}
	return addr
}
