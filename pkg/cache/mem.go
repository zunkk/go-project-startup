package cache

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type MemCache struct {
	c *cache.Cache
}

func NewMemCache(expiredTime time.Duration, cleanupInterval time.Duration) (*MemCache, error) {
	c := cache.New(expiredTime, cleanupInterval)
	return &MemCache{
		c: c,
	}, nil
}

func GetFromMemCache[T any](c *MemCache, ns string, key string) (value T, exist bool) {
	v, ok := c.c.Get(genKey(ns, key))
	if !ok {
		return value, false
	}
	value, ok = v.(T)
	return value, ok
}

func PutToMemCache[T any](c *MemCache, ns string, k string, v T) {
	c.c.Set(genKey(ns, k), v, cache.DefaultExpiration)
}

func genKey(ns string, k string) string {
	return fmt.Sprintf("%s_%s", ns, k)
}
