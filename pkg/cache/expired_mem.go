package cache

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type ExpiredMemCache struct {
	c *cache.Cache
}

func NewExpiredMemCache(expiredTime time.Duration, cleanupInterval time.Duration) (*ExpiredMemCache, error) {
	c := cache.New(expiredTime, cleanupInterval)
	return &ExpiredMemCache{
		c: c,
	}, nil
}

func ExpiredMemCacheGet[T any](c *ExpiredMemCache, ns string, key string) (value T, exist bool) {
	v, ok := c.c.Get(genKey(ns, key))
	if !ok {
		return value, false
	}
	value, ok = v.(T)
	return value, ok
}

func ExpiredMemCachePut[T any](c *ExpiredMemCache, ns string, k string, v T) {
	c.c.Set(genKey(ns, k), v, cache.DefaultExpiration)
}

func genKey(ns string, k string) string {
	return fmt.Sprintf("%s_%s", ns, k)
}
