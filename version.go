package main

import (
	"fmt"
	"time"
)

// Version is definition of a version
type Version interface {
	String() string
	Compare(o interface{}) (int, error)
}

// NewVersion from string
func NewVersion(v string) (Version, error) {
	s, err := NewSemanticVersion(v)
	if err == nil {
		return s, nil
	}
	p, err := NewPseudoVersion(v)
	if err == nil {
		return p, nil
	}
	return NewDateVersion(v)
}

// GetVersionTime extracts time from pseudo or date version
func GetVersionTime(v interface{}) (*time.Time, error) {
	p, ok := v.(*PseudoVersion)
	if ok {
		d := Time2Date(p.Time)
		return &d, nil
	} else {
		d, ok := v.(*DateVersion)
		if !ok {
			return nil, fmt.Errorf("unknown version type: %T", v)
		}
		t := Time2Date(time.Time(*d))
		return &t, nil
	}
}

// Time2Date rounds time to date
func Time2Date(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
