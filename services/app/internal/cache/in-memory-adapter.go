package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type InMemoryAdapter struct {
	cache *cache.Cache
}

func NewInMemoryAdapter() ICache {
	return &InMemoryAdapter{
		cache: cache.New(time.Minute*15, time.Minute*10),
	}
}

func (c *InMemoryAdapter) Get(key string) (any, bool) {
	return c.cache.Get(key)
}

func (c *InMemoryAdapter) Set(key string, value any, duration time.Duration) {
	c.cache.Set(key, value, duration)
}
