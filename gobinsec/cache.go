package gobinsec

var CacheInstance Cache

type Cache interface {
	Get(d *Dependency) []byte
	Set(d *Dependency, v []byte)
	Open() error
	Close()
}

func BuildCache() error {
	if conf := NewMemcachierConfig(config.Memcachier); conf != nil {
		CacheInstance = NewMemcachierCache(conf)
	} else if conf := NewMemcachedConfig(config.Memcached); conf != nil {
		CacheInstance = NewMemcachedCache(conf)
	} else {
		conf := NewFileConfig(config.File)
		CacheInstance = NewFileCache(conf)
	}
	return CacheInstance.Open()
}
