package gobinsec

import "github.com/bradfitz/gomemcache/memcache"

// MemcachedClient is the Cache using memcached
type MemcachedClient struct {
	Client     *memcache.Client
	Expiration int32
}

// NewMemcachedCache builds a memcached cache
func NewMemcachedCache() Cache {
	cache := MemcachedClient{
		Client:     memcache.New(config.Memcached.Address),
		Expiration: config.Memcached.Expiration,
	}
	return &cache
}

// Get returns NVD response for given dependency
func (mc *MemcachedClient) Get(d *Dependency) []byte {
	item, err := mc.Client.Get(d.Key())
	if err != nil {
		return nil
	}
	return item.Value
}

// Set put NVD response for given dependency in cache
func (mc *MemcachedClient) Set(d *Dependency, v []byte) {
	item := memcache.Item{
		Key:        d.Key(),
		Value:      v,
		Expiration: mc.Expiration,
	}
	_ = mc.Client.Set(&item) // we ignore error
}

// Ping calls memcached
func (mc *MemcachedClient) Ping() error {
	return mc.Client.Ping()
}
