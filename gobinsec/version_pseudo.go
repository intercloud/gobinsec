package gobinsec

import (
	"fmt"
	"time"
)

// PseudoVersionTimeFormat is the time format for pseudo versions
const PseudoVersionTimeFormat = "20060102"
const PseudoVersionMinimumLength = 32

// PseudoVersion for dependencies that don't have a version. Its string
// representation is something like "v0.0.0-20191109021931-daa7c04131f5"
// with time and commit ID
type PseudoVersion struct {
	Text string
	Date time.Time
}

// NewPseudoVersion builds a pseudo version from string
func NewPseudoVersion(text string) (*PseudoVersion, error) {
	if len(text) < PseudoVersionMinimumLength {
		return nil, fmt.Errorf("wrong pseudo version format: %s", text)
	}
	version := PseudoVersion{
		Text: text,
	}
	start := len(text) - 27
	end := start + 8
	date := text[start:end]
	var err error
	version.Date, err = time.Parse(PseudoVersionTimeFormat, date)
	if err != nil {
		return nil, fmt.Errorf("wrong pseudo version time: %s", text)
	}
	return &version, nil
}

// String returns a string representation for pseudo version
func (version *PseudoVersion) String() string {
	return version.Text
}

// Compare two pseudo versions by time
func (version *PseudoVersion) Compare(o interface{}) (int, error) {
	t, err := GetVersionTime(o)
	if err != nil {
		return 0, err
	}
	d := version.Date
	if d.Before(*t) {
		return -1, nil
	}
	if d.After(*t) {
		return 1, nil
	}
	return 0, nil
}
