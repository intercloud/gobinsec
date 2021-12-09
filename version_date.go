package main

import (
	"fmt"
	"time"
)

// DateVersionTimeFormat is the time format for date versions
const DateVersionTimeFormat = "2006-01-02"

// DateVersion for dependencies that don't have a version
type DateVersion time.Time

// NewDateVersion builds a date version from string
func NewDateVersion(s string) (*DateVersion, error) {
	t, err := time.Parse(DateVersionTimeFormat, s)
	if err != nil {
		return nil, fmt.Errorf("parsing date version: %v", err)
	}
	v := DateVersion(t)
	return &v, nil
}

// String returns a string representation for date version
func (d *DateVersion) String() string {
	t := time.Time(*d)
	return t.Format(DateVersionTimeFormat)
}

// Compare two date versions by time
func (d *DateVersion) Compare(o interface{}) (int, error) {
	t1 := Time2Date(time.Time(*d))
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
