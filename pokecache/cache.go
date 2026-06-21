package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

type cache struct {
	enteries map[string]cacheEntry
	interval time.Duration
	mu       sync.RWMutex
}

func NewCache(i time.Duration) *cache {
	c := cache{
		interval: i,
		enteries: make(map[string]cacheEntry),
	}
	c.reapLoop()
	return &c
}

func (c *cache) Add(key string, val []byte) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.enteries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.enteries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *cache) reapLoop() {
	ticker := time.NewTicker(c.interval * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			c.mu.RLock()
			for key, entry := range c.enteries {
				if time.Since(entry.createdAt) > c.interval {
					delete(c.enteries, key)
				}
			}
			c.mu.RUnlock()
		}
	}()
}
