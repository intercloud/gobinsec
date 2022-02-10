package gobinsec

import (
	"sync"
)

type DependencyCache map[string][]Vulnerability

var lock sync.RWMutex
var cache = NewDependencyCache()

func NewDependencyCache() *DependencyCache {
	cache := make(DependencyCache)
	return &cache
}

func (dc *DependencyCache) Get(d *Dependency) []Vulnerability {
	key := d.Hash()
	lock.RLock()
	defer lock.RUnlock()
	vulnerabilities, ok := (*dc)[key]
	if ok {
		return vulnerabilities
	}
	return nil
}

func (dc *DependencyCache) Put(d *Dependency, v []Vulnerability) {
	key := d.Hash()
	if v == nil {
		v = make([]Vulnerability, 0)
	}
	lock.Lock()
	defer lock.Unlock()
	(*dc)[key] = v
}
