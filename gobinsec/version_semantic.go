package gobinsec

import (
	"fmt"
	"strconv"
	"strings"
)

// SemanticVersion is type to represent a semantic version
type SemanticVersion struct {
	Text   string
	Fields [3]int
}

// NewSemanticVersion builds a version from string
func NewSemanticVersion(text string) (*SemanticVersion, error) {
	version := SemanticVersion{Text: text}
	parts := strings.TrimPrefix(text, "v")
	parts = strings.TrimSuffix(parts, "+incompatible")
	numbers := strings.Split(parts, ".")
	if len(numbers) != 3 { // nolint:gomnd // 3 is magic
		return nil, fmt.Errorf("wrong version parts count: %s", text)
	}
	var err error
	for i := 0; i < 3; i++ {
		version.Fields[i], err = strconv.Atoi(numbers[i])
		if err != nil {
			return nil, fmt.Errorf("wrong version number: %s", text)
		}
		if version.Fields[i] < 0 {
			return nil, fmt.Errorf("negative version number: %s", text)
		}
	}
	return &version, nil
}

// String returns string representation of version
func (version *SemanticVersion) String() string {
	return version.Text
}

// Compare two semantic versions s and o:
// - if s < o : returns -1
// - if s > o : returns +1
// - if s = o : returns 0
func (version *SemanticVersion) Compare(other interface{}) (int, error) {
	otherVersion, ok := other.(*SemanticVersion)
	if !ok {
		return 0, fmt.Errorf("can't compare semantic version to other type: %T", other)
	}
	for i := 0; i < 3; i++ {
		if version.Fields[i] < otherVersion.Fields[i] {
			return -1, nil
		}
		if version.Fields[i] > otherVersion.Fields[i] {
			return 1, nil
		}
	}
	return 0, nil
}
