package utils

import (
	"github.com/FishGoddess/cachego"
	"time"
)

type Cache struct {
	cache *cachego.Cache
}

func NewCache() *Cache {
	cache := cachego.NewCache()
	return &Cache{
		cache: cache,
	}
}

func (c *Cache) Set(key string, value interface{}, ttl int64) {
	c.cache.Set(key, value, cachego.WithOpTTL(time.Duration(ttl)))
}

func (c *Cache) Get(key string) (interface{}, error) {
	return c.cache.Get(key)
}

func (c *Cache) Close() {
	c.cache.GC()
}

func (c *Cache) Del(key string) {
	c.cache.Delete(key)
}
