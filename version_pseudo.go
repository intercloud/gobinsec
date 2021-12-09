package main

import (
	"fmt"
	"strings"
	"time"
)

// PseudoVersionTimeFormat is the time format for pseudo versions
const PseudoVersionTimeFormat = "20060102150405"

// PseudoVersion for dependencies that don't have a version. Its string
// representation is something like "v0.0.0-20191109021931-daa7c04131f5"
// with time and commit ID
type PseudoVersion struct {
	Time   time.Time // the time of the last commit
	Commit string    // the last commit ID
}

// NewPseudoVersion builds a pseudo version from string
func NewPseudoVersion(s string) (*PseudoVersion, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 3 { // nolint:gomnd // 3 is magic
		return nil, fmt.Errorf("bad pseudo version fields count: %s", s)
	}
	if parts[0] != "v0.0.0" && parts[0] != "0.0.0" {
		return nil, fmt.Errorf("bad pseudo version first part: %s", s)
	}
	var pseudoVersion PseudoVersion
	var err error
	pseudoVersion.Time, err = time.Parse(PseudoVersionTimeFormat, parts[1])
	if err != nil {
		return nil, fmt.Errorf("wrong pseudo version time: %s", s)
	}
	pseudoVersion.Commit = parts[2]
	return &pseudoVersion, nil
}

// String returns a string representation for pseudo version
func (p *PseudoVersion) String() string {
	return fmt.Sprintf("v0.0.0-%s-%s", p.Time.Format(PseudoVersionTimeFormat), p.Commit)
}

// Compare two pseudo versions by time
func (p *PseudoVersion) Compare(o interface{}) (int, error) {
	t, err := GetVersionTime(o)
	if err != nil {
		return 0, err
	}
	d := Time2Date(p.Time)
	if d.Before(*t) {
		return -1, nil
	}
	if d.After(*t) {
		return 1, nil
	}
	return 0, nil
}
