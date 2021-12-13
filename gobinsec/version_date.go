package gobinsec

import (
	"fmt"
	"time"
)

// DateVersionTimeFormat is the time format for date versions
const DateVersionTimeFormat = "2006-01-02"

// DateVersion for dependencies that don't have a version
type DateVersion struct {
	Text string
	Date time.Time
}

// NewDateVersion builds a date version from string
func NewDateVersion(text string) (*DateVersion, error) {
	date, err := time.Parse(DateVersionTimeFormat, text)
	if err != nil {
		return nil, fmt.Errorf("parsing date version: %v", err)
	}
	version := DateVersion{Text: text, Date: date}
	return &version, nil
}

// String returns a string representation for date version
func (version *DateVersion) String() string {
	return version.Text
}

// Compare two date versions by time
func (version *DateVersion) Compare(o interface{}) (int, error) {
	t1 := version.Date
	t2, err := GetVersionTime(o)
	if err != nil {
		return 0, err
	}
	if t1.Before(*t2) {
		return -1, nil
	}
	if t1.After(*t2) {
		return 1, nil
	}
	return 0, nil
}
