package gobinsec

import "fmt"

// CacheLogger logs cache calls on embedded cache
type CacheLogger struct {
	Cache Cache
}

// NewCacheLogger returns a cache logger embedding a cache
func NewCacheLogger(cache Cache) *CacheLogger {
	logger := CacheLogger{
		Cache: cache,
	}
	return &logger
}

// Name return embedded cache name
func (c *CacheLogger) Name() string {
	return c.Cache.Name()
}

// Get calls embedded cache and logs result
func (c *CacheLogger) Get(d *Dependency) ([]byte, error) {
	bytes, err := c.Cache.Get(d)
	if err != nil {
		return nil, err
	}
	if bytes != nil {
		fmt.Printf("<<< %s\n", d.Name)
	} else {
		fmt.Printf("!!! %s\n", d.Name)
	}
	return bytes, nil
}

// Set calls embedded cache and logs call
func (c *CacheLogger) Set(d *Dependency, v []byte) error {
	err := c.Cache.Set(d, v)
	if err != nil {
		return err
	}
	fmt.Printf(">>> %s\n", d.Name)
	return nil
}

// Open prints an information message on terminal
func (c *CacheLogger) Open() error {
	fmt.Printf("Caching with %s\n", c.Name())
	return c.Cache.Open()
}

// Close does nothing
func (c *CacheLogger) Close() error {
	return c.Cache.Close()
}
