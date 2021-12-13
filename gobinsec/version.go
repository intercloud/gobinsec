package gobinsec

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
func NewVersion(version string) (Version, error) {
	semantic, err := NewSemanticVersion(version)
	if err == nil {
		return semantic, nil
	}
	pseudo, err := NewPseudoVersion(version)
	if err == nil {
		return pseudo, nil
	}
	date, err := NewDateVersion(version)
	if err == nil {
		return date, nil
	}
	return nil, fmt.Errorf("bad version: %s", version)
}

// GetVersionTime extracts time from pseudo or date version
func GetVersionTime(version interface{}) (*time.Time, error) {
	pseudo, ok := version.(*PseudoVersion)
	if ok {
		d := pseudo.Date
		return &d, nil
	} else {
		date, ok := version.(*DateVersion)
		if !ok {
			return nil, fmt.Errorf("unknown version type: %T", version)
		}
		d := date.Date
		return &d, nil
	}
}
