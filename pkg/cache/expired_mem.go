package cache

import (
	"time"

	"github.com/maypok86/otter"
)

type ExpiredMemCache[K comparable, V any] struct {
	c otter.Cache[K, V]
}

func NewExpiredMemCache[K comparable, V any](expiredTime time.Duration, capacity int) (*ExpiredMemCache[K, V], error) {
	c, err := otter.MustBuilder[K, V](capacity).
		WithTTL(expiredTime).
		Build()
	if err != nil {
		return nil, err
	}

	return &ExpiredMemCache[K, V]{
		c: c,
	}, nil
}

func (c *ExpiredMemCache[K, V]) Get(k K) (V, bool) {
	return c.c.Get(k)
}

func (c *ExpiredMemCache[K, V]) Has(k K) bool {
	return c.c.Has(k)
}

func (c *ExpiredMemCache[K, V]) Put(k K, v V) {
	c.c.Set(k, v)
}
