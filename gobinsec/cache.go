package gobinsec

import "fmt"

var CacheInstance Cache

// Cache is the interface for caching
type Cache interface {
	Name() string
	Get(d *Dependency) ([]byte, error)
	Set(d *Dependency, v []byte) error
	Open() error
	Close() error
}

// NewCache builds a cache instance depending on configuration and environment
func NewCache() error {
	var err error
	CacheInstance, err = NewMemcachierCache(config.Memcachier)
	if err != nil {
		return fmt.Errorf("configuring memcachier: %v", err)
	}
	if CacheInstance != nil {
		return nil
	}
	CacheInstance, err = NewMemcachedCache(config.Memcached)
	if err != nil {
		return fmt.Errorf("configuring memcached: %v", err)
	}
	if CacheInstance != nil {
		return nil
	}
	CacheInstance, err = NewFileCache(config.File)
	if err != nil {
		return fmt.Errorf("configuring file cache: %v", err)
	}
	return nil
}

// BuildCache and open it
func BuildCache() error {
	if err := NewCache(); err != nil {
		return err
	}
	if config.Cache {
		CacheInstance = NewCacheLogger(CacheInstance)
	}
	return CacheInstance.Open()
}
