package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

type Cache struct {
	enteries map[string]cacheEntry
	interval time.Duration
	mu       sync.RWMutex
}

func NewCache(i time.Duration) *Cache {
	c := Cache{
		interval: i,
		enteries: make(map[string]cacheEntry),
	}
	c.reapLoop()
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.enteries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.enteries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			c.mu.Lock()
			for key, entry := range c.enteries {
				if time.Since(entry.createdAt) > c.interval {
					delete(c.enteries, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}
