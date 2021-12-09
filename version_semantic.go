package main

import (
	"fmt"
	"strconv"
	"strings"
)

// SemanticVersion is type to represent a semantic version
type SemanticVersion [3]int

// NewSemanticVersion builds a version from string
func NewSemanticVersion(s string) (*SemanticVersion, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 { // nolint:gomnd // 3 is magic
		return nil, fmt.Errorf("wrong version parts count: %s", s)
	}
	var version SemanticVersion
	var err error
	parts[0] = strings.TrimPrefix(parts[0], "v")
	for i := 0; i < 3; i++ {
		version[i], err = strconv.Atoi(parts[i])
		if err != nil {
			return nil, fmt.Errorf("wrong version number: %s", s)
		}
		if version[i] < 0 {
			return nil, fmt.Errorf("negative version number: %s", s)
		}
	}
	return &version, nil
}

// String returns string representation of version
func (s *SemanticVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", s[0], s[1], s[2])
}

// Compare two semantic versions s and o:
// - if s < o : returns -1
// - if s > o : returns +1
// - if s = o : returns 0
func (s *SemanticVersion) Compare(o interface{}) (int, error) {
	os, ok := o.(*SemanticVersion)
	if !ok {
		return 0, fmt.Errorf("can't compare semantic version to other type: %T", o)
	}
	for i := 0; i < 3; i++ {
		if s[i] < os[i] {
			return -1, nil
		}
		if s[i] > os[i] {
			return 1, nil
		}
	}
	return 0, nil
}
