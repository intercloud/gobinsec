package gobinsec

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

const (
	MemcachedDefaultTimeout    = 1 * time.Second
	MemcachedDefaultExpiration = 24 * time.Hour
)

// MemcachedConfig is the configuration for memcached
type MemcachedConfig struct {
	Address    string        `yaml:"address"`
	Timeout    time.Duration `yaml:"timeout"`
	Expiration time.Duration `yaml:"expiration"`
}

// MemcachedClient is the Cache using memcached
type MemcachedClient struct {
	Client     *memcache.Client
	Expiration time.Duration
}

// NewMemcachedCache builds a memcached cache
func NewMemcachedCache(config *MemcachedConfig) (Cache, error) {
	config, err := NewMemcachedConfig(config)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}
	client := memcache.New(config.Address)
	client.Timeout = config.Timeout
	cache := MemcachedClient{
		Client:     client,
		Expiration: config.Expiration,
	}
	return &cache, nil
}

// NewMemcachedConfig returns configuration
func NewMemcachedConfig(config *MemcachedConfig) (*MemcachedConfig, error) {
	if config != nil {
		if config.Timeout == 0 {
			config.Timeout = MemcachedDefaultTimeout
		}
		if config.Expiration == 0 {
			config.Expiration = MemcachedDefaultExpiration
		}
		return config, nil
	} else if os.Getenv("MEMCACHED_ADDRESS") != "" {
		address := os.Getenv("MEMCACHED_ADDRESS")
		timeout, err := ParseDuration(os.Getenv("MEMCACHED_TIMEOUT"), MemcachedDefaultTimeout)
		if err != nil {
			return nil, fmt.Errorf("parsing memcached timeout: %v", err)
		}
		expiration, err := ParseDuration(os.Getenv("MEMCACHED_EXPIRATION"), MemcachedDefaultExpiration)
		if err != nil {
			return nil, fmt.Errorf("parsing memcached expiration: %v", err)
		}
		c := MemcachedConfig{
			Address:    address,
			Timeout:    timeout,
			Expiration: expiration,
		}
		return &c, nil
	} else {
		return nil, nil
	}
}

// Name returns the name of this cache
func (mc *MemcachedClient) Name() string {
	return "Memcached"
}

// Get returns NVD response for given dependency
func (mc *MemcachedClient) Get(d *Dependency) ([]byte, error) {
	item, err := mc.Client.Get(d.Key())
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		return nil, fmt.Errorf("getting on memcached: %v", err)
	}
	if item == nil {
		return nil, nil
	}
	return item.Value, nil
}

// Set put NVD response for given dependency in cache
func (mc *MemcachedClient) Set(d *Dependency, v []byte) error {
	item := memcache.Item{
		Key:        d.Key(),
		Value:      v,
		Expiration: int32(mc.Expiration.Seconds()),
	}
	if err := mc.Client.Set(&item); err != nil {
		return fmt.Errorf("setting on memcached: %v", err)
	}
	return nil
}

// Ping calls memcached
func (mc *MemcachedClient) Open() error {
	return mc.Client.Ping()
}

// Clean does nothing
func (mc *MemcachedClient) Close() error {
	return nil
}
