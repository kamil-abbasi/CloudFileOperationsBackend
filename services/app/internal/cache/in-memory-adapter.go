package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type InMemoryAdapter struct {
	cache *cache.Cache
}

func NewInMemoryAdapter() ICache {
	return &InMemoryAdapter{}
}

func (c *InMemoryAdapter) Get(key string) (any, bool) {
	return c.cache.Get(key)
}

func (c *InMemoryAdapter) Set(key string, value any, duration time.Duration) {
	c.cache.Set(key, value, duration)
}
