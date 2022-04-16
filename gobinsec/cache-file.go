package gobinsec

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// FileConfig is the configuration for file cache
type FileConfig struct {
	File       string `yaml:"name"`
	Expiration int32  `yaml:"expiration"`
}

// NewFileConfig builds configuration for file cache
func NewFileConfig(config *FileConfig) *FileConfig {
	if config == nil {
		config = &FileConfig{}
	}
	if config.File == "" {
		config.File = "~/.gobinsec-cache.yml"
	}
	if strings.HasPrefix(config.File, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil
		}
		config.File = filepath.Join(home, config.File[2:])
	}
	if config.Expiration == 0 {
		config.Expiration = 86400
	}
	return config
}

// DependencyCache is an entry of dependency cache
type DependencyCache struct {
	Date            string
	Vulnerabilities string
}

// FileCache is the cache instance
type FileCache struct {
	File       string
	Expiration int32
	Cache      map[string]DependencyCache
}

// NewFileCache builds a cache using file
func NewFileCache(config *FileConfig) Cache {
	cache := make(map[string]DependencyCache)
	file, err := ioutil.ReadFile(config.File)
	if err == nil {
		yaml.Unmarshal(file, cache)
	}
	fileCache := &FileCache{
		File:       config.File,
		Expiration: config.Expiration,
		Cache:      cache,
	}
	fileCache.CleanCache()
	return fileCache
}

// Get returns NVD response for given dependency
func (fc *FileCache) Get(d *Dependency) []byte {
	dependency, ok := fc.Cache[d.Name]
	if ok {
		return []byte(dependency.Vulnerabilities)
	}
	return nil
}

// Set put NVD response for given dependency in cache
func (fc *FileCache) Set(d *Dependency, v []byte) {
	now := time.Now().UTC().Format("2006-01-02T15:04:05")
	cache := DependencyCache{
		Date:            now,
		Vulnerabilities: string(v),
	}
	fc.Cache[d.Name] = cache
}

// Ping does nothing
func (fc *FileCache) Open() error {
	return nil
}

// Clean saves file
func (fc *FileCache) Close() {
	text, err := yaml.Marshal(fc.Cache)
	if err == nil {
		ioutil.WriteFile(fc.File, text, 0644)
	}
}

// CleanCache removes obsolete cache entries
func (fc *FileCache) CleanCache() {
	duration := time.Duration(fc.Expiration) * time.Second
	limit := time.Now().UTC().Add(-duration).Format("2006-01-02T15:04:05")
	for name, cache := range fc.Cache {
		if cache.Date < limit {
			delete(fc.Cache, name)
		}
	}
}
