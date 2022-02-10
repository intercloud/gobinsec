package gobinsec

import (
	"sync"
)

// MemoryCache is a cache with a synced map
type MemoryCache map[string][]byte

// lock for memory cache
var lock sync.RWMutex

// NewMemoryCache builds a cache in memory
func NewMemoryCache() *MemoryCache {
	cache := make(MemoryCache)
	return &cache
}

// Get returns NVD response for given dependency
func (mc *MemoryCache) Get(d *Dependency) []byte {
	key := d.Key()
	lock.RLock()
	defer lock.RUnlock()
	vulnerabilities, ok := (*mc)[key]
	if ok {
		return vulnerabilities
	}
	return nil
}

// Set put NVD response for given dependency in cache
func (mc *MemoryCache) Set(d *Dependency, v []byte) {
	key := d.Key()
	lock.Lock()
	defer lock.Unlock()
	(*mc)[key] = v
}

// Ping does nothing
func (mc *MemoryCache) Ping() error {
	return nil
}
