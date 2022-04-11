package gobinsec

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIKey     string            `yaml:"api-key"`
	Memcached  *MemcachedConfig  `yaml:"memcached"`
	Memcachier *MemcachierConfig `yaml:"memcachier"`
	SQLite     *SQLiteConfig     `yaml:"sqlite"`
	Ignore     []string          `yaml:"ignore"`
	Strict     bool              `yaml:"strict"`
}

var config Config

// LoadConfig loads configuration from given file
func LoadConfig(path string, strict bool) error {
	if path == "" {
		config.Strict = strict
		return nil
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("loading configuration file: %v", err)
	}
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return fmt.Errorf("parsing configuration: %v", err)
	}
	if config.APIKey == "" {
		config.APIKey = os.Getenv("NVD_API_KEY")
	}
	if strict {
		config.Strict = strict
	}
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
