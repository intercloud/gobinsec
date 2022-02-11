package gobinsec

var cache Cache

type Cache interface {
	Get(d *Dependency) []byte
	Set(d *Dependency, v []byte)
	Ping() error
}

func BuildCache() error {
	if conf := NewMemcachierConfig(config.Memcachier); conf != nil {
		cache = NewMemcachierCache(conf)
	} else if conf := NewMemcachedConfig(config.Memcached); conf != nil {
		cache = NewMemcachedCache(conf)
	} else {
		cache = NewMemoryCache()
	}
	return cache.Ping()
}
