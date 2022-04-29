package gobinsec

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	URL                  = "https://services.nvd.nist.gov/rest/json/cves/1.0/?keyword="
	StatusCodeLimit      = 300
	RateMinuteWithoutKey = 10
	RateMinuteWithKey    = 100
)

// timeLastCall is time for last NVD API call
var timeLastCall time.Time

// Dependency is a dependency with vulnerabilities
type Dependency struct {
	Name            string
	Version         Version
	Vulnerabilities []Vulnerability
	Vulnerable      bool
}

// NewDependency builds a new dependency and loads its vulnerabilities
func NewDependency(name, version string) (*Dependency, error) {
	v := NewVersion(version)
	dependency := Dependency{
		Name:    name,
		Version: v,
	}
	return &dependency, nil
}

// Vulnerabilities return list of vulnerabilities for given dependency
func (d *Dependency) LoadVulnerabilities() error {
	vulnerabilities, err := CacheInstance.Get(d)
	if err != nil {
		return err
	}
	if vulnerabilities == nil {
		WaitBeforeCall()
		url := URL + d.Name
		if config.APIKey != "" {
			url += "&apiKey=" + config.APIKey
		}
		response, err := http.Get(url) // nolint:gosec // it's safe!
		if err != nil {
			return fmt.Errorf("calling NVD: %v", err)
		}
		defer response.Body.Close()
		if response.StatusCode >= StatusCodeLimit {
			return fmt.Errorf("bad status code calling NVD: %d", response.StatusCode)
		}
		vulnerabilities, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		if err := CacheInstance.Set(d, vulnerabilities); err != nil {
			return err
		}
	}
	var result Response
	if err := json.Unmarshal(vulnerabilities, &result); err != nil {
		return fmt.Errorf("decoding JSON response: %v", err)
	}
	for _, item := range result.Result.Items {
		vulnerability, err := NewVulnerability(item)
		if err != nil {
			return err
		}
		if vulnerability.IsExposed(d.Version) &&
			!vulnerability.Ignored {
			d.Vulnerable = true
		}
		d.Vulnerabilities = append(d.Vulnerabilities, *vulnerability)
	}
	return nil
}

// WaitBeforeCall waits in order not to exceed NVD call rate limit
func WaitBeforeCall() {
	if config.Wait {
		elapsedSinceLastCall := time.Since(timeLastCall)
		timeToSleep := (60000 / RateMinuteWithoutKey) * time.Millisecond
		if config.APIKey != "" {
			timeToSleep = (60000 / RateMinuteWithKey) * time.Millisecond
		}
		timeToWait := timeToSleep - elapsedSinceLastCall
		if timeToWait > 0 {
			time.Sleep(timeToWait)
		}
		timeLastCall = time.Now()
	}
}

// Key returns a key as a string for caching
func (d *Dependency) Key() string {
	return d.Name
}
