package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

type CacheItemWithId struct {
	CacheItem
}

type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		items: make(map[string]CacheItem),
	}
	go cache.StartCleanup(interval)
	return cache
}

func (c *Cache) Set(key string, data interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found || item.ExpiresAt.Before(time.Now()) {
		return nil, false
	}
	return item.Data, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, item := range c.items {
		if item.ExpiresAt.Before(time.Now()) {
			delete(c.items, key)
		}
	}
}

func (c *Cache) StartCleanup(interval time.Duration) {
	for {
		time.Sleep(interval)
		c.Cleanup()
	}
}

var (
	AppCache *Cache
	JWT      *Cache
)

func Init() {
	AppCache = NewCache(10 * time.Minute)
	JWT = NewCache(1000 * time.Hour)
}
