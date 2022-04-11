package gobinsec

var CacheInstance Cache

type Cache interface {
	Get(d *Dependency) []byte
	Set(d *Dependency, v []byte)
	Ping() error
	Clean()
}

func BuildCache() error {
	if conf := NewMemcachedConfig(config.Memcached); conf != nil {
		CacheInstance = NewMemcachedCache(conf)
	} else if conf := NewMemcachierConfig(config.Memcachier); conf != nil {
		CacheInstance = NewMemcachierCache(conf)
	} else {
		conf := NewSQLiteConfig(config.SQLite)
		CacheInstance = NewSQLiteCache(conf)
	}
	return CacheInstance.Ping()
}
