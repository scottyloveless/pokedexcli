package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Map   map[string]CacheEntry
	Mutex *sync.Mutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

var GlobalCache Cache

const reapInterval = time.Second * 5

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		Map:   map[string]CacheEntry{},
		Mutex: &sync.Mutex{},
	}

	go cache.reapLoop(interval)

	return cache
}

func (c Cache) Add(key string, val []byte) {
	newCacheEntry := CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.Mutex.Lock()
	c.Map[key] = newCacheEntry
	c.Mutex.Unlock()
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.Mutex.Lock()
	e, exists := c.Map[key]
	c.Mutex.Unlock()
	if exists {
		return e.val, exists
	} else {
		return nil, exists
	}
}

func (c Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		c.Mutex.Lock()
		toDelete := []string{}
		for k, v := range c.Map {
			if time.Since(v.createdAt) > interval {
				toDelete = append(toDelete, k)
			}
		}
		for _, v := range toDelete {
			delete(c.Map, v)
		}
		c.Mutex.Unlock()

	}
}
