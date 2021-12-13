package gobinsec

import "fmt"

type UnknownVersion string

func NewUnknownVersion(version string) *UnknownVersion {
	unknown := UnknownVersion(version)
	return &unknown
}

func (version *UnknownVersion) String() string {
	return string(*version)
}

func (version *UnknownVersion) Compare(other interface{}) (int, error) {
	return 0, fmt.Errorf("can't compare unknown version to any other")
}
