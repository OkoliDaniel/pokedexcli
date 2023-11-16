package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.Mutex
}

func (c Cache) Add(key string, value []byte) {
	defer c.mu.Unlock()
	c.mu.Lock()
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
	c.cache[key] = entry
}

func (c Cache) Get(key string) (val []byte, found bool) {
	defer c.mu.Unlock()
	c.mu.Lock()
	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c Cache) reapLoop(tickChan <-chan time.Time, dur time.Duration) {
	for tick := range tickChan {
		c.mu.Lock()
		for k, v := range c.cache {
			timeCreated := v.createdAt
			if tick.Sub(timeCreated) >= dur {
				delete(c.cache, k)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(d time.Duration) Cache {
	ticker := time.NewTicker(d)
	tickChan := ticker.C
	cache := Cache{
		cache: make(map[string]cacheEntry),
		mu:    &sync.Mutex{},
	}
	go cache.reapLoop(tickChan, d)

	return cache
}
