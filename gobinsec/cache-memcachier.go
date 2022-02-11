package gobinsec

import (
	"os"
	"strconv"

	"github.com/memcachier/gomemcache/memcache"
)

// MemcachierConfig is the configuration for Memcachier
type MemcachierConfig struct {
	Address    string `yaml:"address"`
	Expiration int32  `yaml:"expiration"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
}

// MemcachierClient is the Cache using Memcachier
type MemcachierClient struct {
	Client     *memcache.Client
	Expiration int32
}

// NewMemcachierConfig returns configuration
func NewMemcachierConfig(config *MemcachierConfig) *MemcachierConfig {
	var username, password, address string
	var expiration int32
	if config != nil {
		return config
	} else if os.Getenv("MEMCACHIER_ADDRESS") != "" {
		username = os.Getenv("MEMCACHIER_USERNAME")
		password = os.Getenv("MEMCACHIER_PASSWORD")
		address = os.Getenv("MEMCACHIER_ADDRESS")
		exp, err := strconv.Atoi(os.Getenv("MEMCACHIER_EXPIRATION"))
		if err != nil {
			return nil
		}
		expiration = int32(exp)
		return &MemcachierConfig{
			Address:    address,
			Expiration: expiration,
			Username:   username,
			Password:   password,
		}
	} else {
		return nil
	}
}

// NewMemcachierCache builds a Memcachier cache
func NewMemcachierCache(config *MemcachierConfig) Cache {
	client := memcache.New(config.Address)
	client.SetAuth(config.Username, []byte(config.Password))
	cache := MemcachierClient{
		Client:     client,
		Expiration: config.Expiration,
	}
	return &cache
}

// Get returns NVD response for given dependency
func (mc *MemcachierClient) Get(d *Dependency) []byte {
	item, err := mc.Client.Get(d.Key())
	if err != nil {
		return nil
	}
	return item.Value
}

// Set put NVD response for given dependency in cache
func (mc *MemcachierClient) Set(d *Dependency, v []byte) {
	item := memcache.Item{
		Key:        d.Key(),
		Value:      v,
		Expiration: mc.Expiration,
	}
	mc.Client.Set(&item) // nolint:errcheck // we ignore error
}

// Ping calls Memcachier
func (mc *MemcachierClient) Ping() error {
	item := memcache.Item{
		Key:        "test",
		Value:      []byte("test"),
		Expiration: mc.Expiration,
	}
	return mc.Client.Set(&item)
}
