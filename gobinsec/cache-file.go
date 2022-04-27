package gobinsec

import (
	"fmt"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	FileCacheDefaultFile       = "~/.gobinsec-cache.yml"
	FileCacheDefaultExpiration = 24 * time.Hour
)

// FileConfig is the configuration for file cache
type FileConfig struct {
	File       string        `yaml:"name"`
	Expiration time.Duration `yaml:"expiration"`
}

// DependencyCache is an entry of dependency cache
type DependencyCache struct {
	Date            string
	Vulnerabilities string
}

// FileCache is the cache instance
type FileCache struct {
	File       string
	Expiration time.Duration
	Cache      map[string]DependencyCache
}

// lock for memory cache
var lock sync.RWMutex

// NewFileCache builds a cache using file
func NewFileCache(config *FileConfig) (Cache, error) {
	var err error
	config, err = NewFileConfig(config)
	if err != nil {
		return nil, err
	}
	cache := make(map[string]DependencyCache)
	fileCache := FileCache{
		File:       config.File,
		Expiration: config.Expiration,
		Cache:      cache,
	}
	return &fileCache, nil
}

// NewFileConfig builds configuration for file cache
func NewFileConfig(config *FileConfig) (*FileConfig, error) {
	if config != nil {
		if config.File == "" {
			config.File = FileCacheDefaultFile
		}
		if config.Expiration == 0 {
			config.Expiration = FileCacheDefaultExpiration
		}
	} else {
		file := os.Getenv("FILECACHE_FILE")
		if file == "" {
			file = FileCacheDefaultFile
		}
		expiration, err := ParseDuration(os.Getenv("FILECACHE_EXPIRATION"), FileCacheDefaultExpiration)
		if err != nil {
			return nil, fmt.Errorf("parsing file cache expiration: %v", err)
		}
		config = &FileConfig{
			File:       file,
			Expiration: expiration,
		}
	}
	config.File = ExpandHome(config.File)
	return config, nil
}

// Name returns the name of this cache
func (fc *FileCache) Name() string {
	return "File"
}

// Get returns NVD response for given dependency
func (fc *FileCache) Get(d *Dependency) ([]byte, error) {
	key := d.Key()
	lock.RLock()
	defer lock.RUnlock()
	dependency, ok := fc.Cache[key]
	if ok {
		return []byte(dependency.Vulnerabilities), nil
	}
	return nil, nil
}

// Set put NVD response for given dependency in cache
func (fc *FileCache) Set(d *Dependency, v []byte) error {
	key := d.Key()
	lock.Lock()
	defer lock.Unlock()
	now := time.Now().UTC().Format("2006-01-02T15:04:05")
	cache := DependencyCache{
		Date:            now,
		Vulnerabilities: string(v),
	}
	fc.Cache[key] = cache
	return nil
}

// Open and load file cache if exists
func (fc *FileCache) Open() error {
	if FileExists(fc.File) {
		file, err := os.ReadFile(fc.File)
		if err != nil {
			return fmt.Errorf("reading file cache: %v", err)
		} else {
			if err := yaml.Unmarshal(file, fc.Cache); err != nil {
				return fmt.Errorf("unmarshaling file cache: %v", err)
			}
		}
		fc.CleanCache()
	}
	return nil
}

// Close saves cache in file
func (fc *FileCache) Close() error {
	text, err := yaml.Marshal(fc.Cache)
	if err != nil {
		return fmt.Errorf("marshaling file cache: %v", err)
	}
	if err := os.WriteFile(fc.File, text, 0644); err != nil { // nolint:gosec,gomnd // no confidential data in file
		return fmt.Errorf("saving file cache: %v", err)
	}
	return nil
}

// CleanCache removes obsolete cache entries
func (fc *FileCache) CleanCache() {
	limit := time.Now().UTC().Add(-fc.Expiration).Format("2006-01-02T15:04:05")
	for name, cache := range fc.Cache {
		if cache.Date < limit {
			delete(fc.Cache, name)
		}
	}
}
