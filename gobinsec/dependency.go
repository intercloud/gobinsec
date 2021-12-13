package gobinsec

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const URL = "https://services.nvd.nist.gov/rest/json/cves/1.0/?keyword="

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
	if err := dependency.LoadVulnerabilities(); err != nil {
		return nil, err
	}
	for i := range dependency.Vulnerabilities {
		if dependency.Vulnerabilities[i].IsExposed(v) &&
			!dependency.Vulnerabilities[i].Ignored {
			dependency.Vulnerable = true
		}
	}
	return &dependency, nil
}

// Vulnerabilities return list of vulnerabilities for given dependency
func (d *Dependency) LoadVulnerabilities() error {
	url := URL + d.Name
	if config.APIKey != "" {
		url += "&apiKey=" + config.APIKey
	}
	response, err := http.Get(url) // nolint:noctx,gosec // it's safe!
	if err != nil {
		return fmt.Errorf("calling NVD: %v", err)
	}
	defer response.Body.Close()
	var result Response
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return fmt.Errorf("decoding JSON response: %v", err)
	}
	for _, item := range result.Result.Items {
		vulnerability, err := NewVulnerability(item)
		if err != nil {
			return err
		}
		d.Vulnerabilities = append(d.Vulnerabilities, *vulnerability)
	}
	return nil
}
