package gobinsec

import (
	"sync"
)

type VulnerabilityCache map[string][]Vulnerability

var lock sync.RWMutex
var cache = NewVulnerabilityCache()

func NewVulnerabilityCache() *VulnerabilityCache {
	cache := make(VulnerabilityCache)
	return &cache
}

func (dc *VulnerabilityCache) Get(d *Dependency) []Vulnerability {
	key := d.Key()
	lock.RLock()
	defer lock.RUnlock()
	vulnerabilities, ok := (*dc)[key]
	if ok {
		return vulnerabilities
	}
	return nil
}

func (dc *VulnerabilityCache) Put(d *Dependency, v []Vulnerability) {
	key := d.Key()
	if v == nil {
		v = make([]Vulnerability, 0)
	}
	lock.Lock()
	defer lock.Unlock()
	(*dc)[key] = v
}
