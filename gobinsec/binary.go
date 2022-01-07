package gobinsec

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const (
	MinimumBinaryDependencyFields = 3
	MinimumBinaryLines            = 3
)

// NumGoroutines to load vulnerabilities
var NumGoroutines = 4 * runtime.NumCPU()

// Binary represents a binary with its dependencies
type Binary struct {
	Path         string        // path to binary file
	Dependencies []*Dependency // list of dependencies
	Vulnerable   bool          // tells if binary is vulnerable
}

// NewBinary returns a binary
func NewBinary(path string) (*Binary, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	binary := Binary{
		Path: path,
	}
	if err := binary.GetDependencies(); err != nil {
		return nil, err
	}
	return &binary, nil
}

// GetDependencies gets dependencies analyzing binary
func (b *Binary) GetDependencies() error {
	stdout, stderr, err := ExecCommand("go", "version", "-m", b.Path)
	if err != nil {
		return err
	}
	if stderr != "" {
		return fmt.Errorf(stderr)
	}
	lines := strings.Split(stdout, "\n")
	if len(lines) < MinimumBinaryLines {
		return fmt.Errorf(stdout)
	}
	for _, line := range lines[3:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < MinimumBinaryDependencyFields {
			continue
		}
		name := parts[1]
		version := parts[2]
		dependency, err := NewDependency(name, version)
		if err != nil {
			return err
		}
		b.Dependencies = append(b.Dependencies, dependency)
	}
	dependencies := make(chan *Dependency, len(b.Dependencies))
	var wg sync.WaitGroup
	for _, dependency := range b.Dependencies {
		dependencies <- dependency
		wg.Add(1)
	}
	for i := 0; i < NumGoroutines; i++ {
		go LoadVulnerabilities(dependencies, &wg)
	}
	wg.Wait()
	for _, dependency := range b.Dependencies {
		if dependency.Vulnerable {
			b.Vulnerable = true
		}
	}
	return nil
}

// LoadVulnerabilities takes dependencies from channel and loads dependencies for it
func LoadVulnerabilities(dependencies chan *Dependency, wg *sync.WaitGroup) {
	for {
		select {
		case dependency := <-dependencies:
			if err := dependency.LoadVulnerabilities(); err != nil {
				fmt.Printf("ERROR loading vulnerability: %v\n", err)
				os.Exit(1)
			}
			wg.Done()
		default:
			return
		}
	}
}

// Report prints a report on terminal
// nolint:gocyclo // this is life
func (b *Binary) Report(verbose bool) {
	fmt.Printf("binary: '%s'\n", filepath.Base(b.Path))
	fmt.Printf("vulnerable: %t\n", b.Vulnerable)
	if len(b.Dependencies) > 0 && (b.Vulnerable || verbose) {
		fmt.Println("dependencies:")
		for _, dependency := range b.Dependencies {
			if !dependency.Vulnerable && !verbose {
				continue
			}
			fmt.Printf("- name:    '%s'\n", dependency.Name)
			fmt.Printf("  version: '%s'\n", dependency.Version)
			fmt.Printf("  vulnerable: %t\n", dependency.Vulnerable)
			if len(dependency.Vulnerabilities) > 0 {
				fmt.Println("  vulnerabilities:")
				for _, vulnerability := range dependency.Vulnerabilities {
					if !vulnerability.Exposed && !verbose {
						continue
					}
					fmt.Printf("  - id: '%s'\n", vulnerability.ID)
					fmt.Printf("    exposed: %t\n", vulnerability.Exposed)
					fmt.Printf("    ignored: %t\n", vulnerability.Ignored)
					fmt.Println("    references:")
					for _, reference := range vulnerability.References {
						fmt.Printf("    - '%s'\n", reference)
					}
					fmt.Println("    matchs:")
					for _, match := range vulnerability.Matchs {
						var text string
						if match.VersionStartExcluding != nil ||
							match.VersionStartIncluding != nil ||
							match.VersionEndExcluding != nil ||
							match.VersionEndIncluding != nil {
							var parts []string
							if match.VersionStartExcluding != nil {
								parts = append(parts, fmt.Sprintf("%v <", match.VersionStartExcluding))
							}
							if match.VersionStartIncluding != nil {
								parts = append(parts, fmt.Sprintf("%v <=", match.VersionStartIncluding))
							}
							parts = append(parts, "v")
							if match.VersionEndExcluding != nil {
								parts = append(parts, fmt.Sprintf("< %v", match.VersionEndExcluding))
							}
							if match.VersionEndIncluding != nil {
								parts = append(parts, fmt.Sprintf("<= %v", match.VersionEndIncluding))
							}
							text = strings.Join(parts, " ")
						} else {
							text = "?"
						}
						fmt.Printf("    - '%s'\n", text)
					}
				}
			}
		}
	}
}
