package internal

import (
	"sync"
	"time"
)

type Cache struct {
	Map   map[string]cacheEntry
	Mutex *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

var (
	cache        Cache
	reapInterval = time.Second * 5
)

func NewCache(interval time.Duration) {
	cache = Cache{
		Map:   map[string]cacheEntry{},
		Mutex: &sync.Mutex{},
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		cache.reapLoop(interval)
	}
}

func (c Cache) Add(key string, val []byte) {
	newCacheEntry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	cache.Mutex.Lock()
	cache.Map[key] = newCacheEntry
	cache.Mutex.Unlock()
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.Mutex.Lock()
	e, exists := cache.Map[key]
	c.Mutex.Unlock()
	if exists {
		return e.val, exists
	} else {
		return nil, exists
	}
}

func (c Cache) reapLoop(interval time.Duration) {
	c.Mutex.Lock()
	for _, entry := range cache.Map {
		if time.Since(entry.createdAt) > interval {
		}
	}
	c.Mutex.Unlock()
}
