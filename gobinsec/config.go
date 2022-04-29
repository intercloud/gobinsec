package gobinsec

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the configuration from YAML config file and command line options
type Config struct {
	APIKey     string            `yaml:"api-key"`
	Memcached  *MemcachedConfig  `yaml:"memcached"`
	Memcachier *MemcachierConfig `yaml:"memcachier"`
	File       *FileConfig       `yaml:"file"`
	Ignore     []string          `yaml:"ignore"`
	Strict     bool              `yaml:"strict"`
	Verbose    bool              `yaml:"verbose"`
	Cache      bool              `yaml:"cache"`
	Wait       bool              `yaml:"wait"`
}

// config is shared
var config Config

// LoadConfig loads configuration from given file and overwrite with command line options
func LoadConfig(path string, strict, wait, verbose, cache bool) error {
	if path != "" {
		bytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("loading configuration file: %v", err)
		}
		if err := yaml.Unmarshal(bytes, &config); err != nil {
			return fmt.Errorf("parsing configuration: %v", err)
		}
	}
	if config.APIKey == "" {
		config.APIKey = os.Getenv("NVD_API_KEY")
	}
	if strict {
		config.Strict = true
	}
	if wait {
		config.Wait = true
	}
	config.Verbose = verbose
	config.Cache = cache
	return nil
}

// IgnoreVulnerability tells if we should ignore given vulnerability
func (c *Config) IgnoreVulnerability(id string) bool {
	for _, ignore := range c.Ignore {
		if ignore == id {
			return true
		}
	}
	return false
}
