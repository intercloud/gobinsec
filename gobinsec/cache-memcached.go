package gobinsec

import (
	"os"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcachedConfig is the configuration for memcached
type MemcachedConfig struct {
	Address    string `yaml:"address"`
	Expiration int32  `yaml:"expiration"`
}

// NewMemcachedConfig returns configuration
func NewMemcachedConfig(config *MemcachedConfig) *MemcachedConfig {
	var address string
	var expiration int32
	if config != nil {
		return config
	} else if os.Getenv("MEMCACHED_ADDRESS") != "" {
		address = os.Getenv("MEMCACHED_ADDRESS")
		exp, err := strconv.Atoi(os.Getenv("MEMCACHED_EXPIRATION"))
		if err != nil {
			return nil
		}
		expiration = int32(exp)
		return &MemcachedConfig{
			Address:    address,
			Expiration: expiration,
		}
	} else {
		return nil
	}
}

// MemcachedClient is the Cache using memcached
type MemcachedClient struct {
	Client     *memcache.Client
	Expiration int32
}

// NewMemcachedCache builds a memcached cache
func NewMemcachedCache(config *MemcachedConfig) Cache {
	cache := MemcachedClient{
		Client:     memcache.New(config.Address),
		Expiration: config.Expiration,
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
	mc.Client.Set(&item) // nolint:errcheck // we ignore error
}

// Ping calls memcached
func (mc *MemcachedClient) Ping() error {
	return mc.Client.Ping()
}
