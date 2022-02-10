package gobinsec

var cache Cache

type Cache interface {
	Get(d *Dependency) []byte
	Set(d *Dependency, v []byte)
	Ping() error
}

func BuildCache() error {
	if config.Memcached == nil {
		cache = NewMemoryCache()
	} else {
		cache = NewMemcachedCache()
	}
	return cache.Ping()
}
