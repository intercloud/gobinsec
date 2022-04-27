package gobinsec

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/memcachier/gomemcache/memcache"
)

const (
	MemcachierDefaultTimeout    = 1 * time.Second
	MemcachierDefaultExpiration = 24 * time.Hour
)

// MemcachierConfig is the configuration for Memcachier
type MemcachierConfig struct {
	Address    string        `yaml:"address"`
	Username   string        `yaml:"username"`
	Password   string        `yaml:"password"`
	Timeout    time.Duration `yaml:"timeout"`
	Expiration time.Duration `yaml:"expiration"`
}

// MemcachierClient is the Cache using Memcachier
type MemcachierClient struct {
	Client     *memcache.Client
	Expiration time.Duration
}

// NewMemcachierCache builds a Memcachier cache
func NewMemcachierCache(config *MemcachierConfig) (Cache, error) {
	config, err := NewMemcachierConfig(config)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}
	client := memcache.New(config.Address)
	client.SetAuth(config.Username, []byte(config.Password))
	client.Timeout = config.Timeout
	cache := MemcachierClient{
		Client:     client,
		Expiration: config.Expiration,
	}
	return &cache, nil
}

// NewMemcachierConfig returns configuration
func NewMemcachierConfig(config *MemcachierConfig) (*MemcachierConfig, error) {
	if config != nil {
		if config.Timeout == 0 {
			config.Timeout = MemcachierDefaultTimeout
		}
		if config.Expiration == 0 {
			config.Expiration = MemcachierDefaultExpiration
		}
		return config, nil
	} else if os.Getenv("MEMCACHIER_ADDRESS") != "" &&
		os.Getenv("MEMCACHIER_USERNAME") != "" &&
		os.Getenv("MEMCACHIER_PASSWORD") != "" {
		address := os.Getenv("MEMCACHIER_ADDRESS")
		username := os.Getenv("MEMCACHIER_USERNAME")
		password := os.Getenv("MEMCACHIER_PASSWORD")
		timeout, err := ParseDuration(os.Getenv("MEMCACHIER_TIMEOUT"), MemcachierDefaultTimeout)
		if err != nil {
			return nil, fmt.Errorf("parsing memcachier timeout: %v", err)
		}
		expiration, err := ParseDuration(os.Getenv("MEMCACHIER_EXPIRATION"), MemcachierDefaultExpiration)
		if err != nil {
			return nil, fmt.Errorf("parsing memcachier expiration: %v", err)
		}
		c := MemcachierConfig{
			Address:    address,
			Username:   username,
			Password:   password,
			Timeout:    timeout,
			Expiration: expiration,
		}
		return &c, nil
	} else {
		return nil, nil
	}
}

// Name returns the name of this cache
func (mc *MemcachierClient) Name() string {
	return "Memcachier"
}

// Get returns NVD response for given dependency
func (mc *MemcachierClient) Get(d *Dependency) ([]byte, error) {
	item, err := mc.Client.Get(d.Key())
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		return nil, fmt.Errorf("getting on memcachier: %v", err)
	}
	if item == nil {
		return nil, nil
	}
	return item.Value, nil
}

// Set put NVD response for given dependency in cache
func (mc *MemcachierClient) Set(d *Dependency, v []byte) error {
	item := memcache.Item{
		Key:        d.Key(),
		Value:      v,
		Expiration: int32(mc.Expiration.Seconds()),
	}
	if err := mc.Client.Set(&item); err != nil {
		return fmt.Errorf("setting on memcachier: %v", err)
	}
	return nil
}

// Ping calls Memcachier
func (mc *MemcachierClient) Open() error {
	item := memcache.Item{
		Key:        "test",
		Value:      []byte("test"),
		Expiration: int32(mc.Expiration.Seconds()),
	}
	return mc.Client.Set(&item)
}

// Clean does nothing
func (mc *MemcachierClient) Close() error {
	return nil
}
