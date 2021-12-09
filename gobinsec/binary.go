package gobinsec

import (
	"fmt"
	"os"
	"strings"
)

// Binary represents a binary with its dependencies
type Binary struct {
	Path         string       // path to binary file
	Dependencies []Dependency // list of dependencies
	Vulnerable   bool         // tells if binary is vulnerable
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
	if stderr == fmt.Sprintf("%s: not executable file", b.Path) {
		return fmt.Errorf(stderr)
	}
	for _, line := range strings.Split(stdout, "\n")[3:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		name := parts[1]
		version := parts[2]
		dependency, err := NewDependency(name, version)
		if err != nil {
			return err
		}
		if dependency.Vulnerable {
			b.Vulnerable = true
		}
		b.Dependencies = append(b.Dependencies, *dependency)
	}
	return nil
}

// Report prints a report on terminal
// nolint:gocyclo // this is life
func (b *Binary) Report(verbose bool) {
	fmt.Printf("binary: '%s'\n", b.Path)
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
					fmt.Println("    versions:")
					for _, match := range vulnerability.Matchs {
						if match.VersionStartExcluding != nil {
							fmt.Printf("    - >: '%v'\n", match.VersionStartExcluding)
						}
						if match.VersionStartIncluding != nil {
							fmt.Printf("    - >=: '%v'\n", match.VersionStartIncluding)
						}
						if match.VersionEndExcluding != nil {
							fmt.Printf("    - <: '%v'\n", match.VersionEndExcluding)
						}
						if match.VersionEndIncluding != nil {
							fmt.Printf("    - <=: '%v'\n", match.VersionEndIncluding)
						}
					}
				}
			}
		}
	}
}
