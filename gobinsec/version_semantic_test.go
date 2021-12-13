package gobinsec

import "testing"

func TestNewSemanticVersion(t *testing.T) {
	version, err := NewSemanticVersion("1.2.3")
	if err != nil {
		t.Fatalf("error parsing version: %v", err)
	}
	if version.Fields[0] != 1 {
		t.Fatalf("bad major version: %d", version.Fields[0])
	}
	if version.Fields[1] != 2 {
		t.Fatalf("bad minor version: %d", version.Fields[1])
	}
	if version.Fields[2] != 3 {
		t.Fatalf("bad debug version: %d", version.Fields[2])
	}
}

func TestNewSemanticVersionErrors(t *testing.T) { // nolint:gocyclo // testing
	_, err := NewSemanticVersion("x")
	if err == nil || err.Error() != "wrong version parts count: x" {
		t.Fatalf("wrong error parsing version: %v", err)
	}
	_, err = NewSemanticVersion("1.2")
	if err == nil || err.Error() != "wrong version parts count: 1.2" {
		t.Fatalf("wrong error parsing version: %v", err)
	}
	_, err = NewSemanticVersion("x.y.z")
	if err == nil || err.Error() != "wrong version number: x.y.z" {
		t.Fatalf("wrong error parsing version: %v", err)
	}
	_, err = NewSemanticVersion("1.y.z")
	if err == nil || err.Error() != "wrong version number: 1.y.z" {
		t.Fatalf("wrong error parsing version: %v", err)
	}
	_, err = NewSemanticVersion("1.2.z")
	if err == nil || err.Error() != "wrong version number: 1.2.z" {
		t.Fatalf("wrong error parsing version: %v", err)
	}
	_, err = NewSemanticVersion("-1.2.3")
	if err == nil || err.Error() != "negative version number: -1.2.3" {
		t.Fatalf("wrong error parsing version: %v", err)
	}
	_, err = NewSemanticVersion("1.-2.3")
	if err == nil || err.Error() != "negative version number: 1.-2.3" {
		t.Fatalf("wrong error parsing version: %v", err)
	}
	_, err = NewSemanticVersion("1.2.-3")
	if err == nil || err.Error() != "negative version number: 1.2.-3" {
		t.Fatalf("wrong error parsing version: %v", err)
	}
}

func TestVersionSemanticString(t *testing.T) {
	version, err := NewSemanticVersion("1.2.3")
	if err != nil {
		t.Fatalf("error parsing version: %v", err)
	}
	if version.String() != "1.2.3" { // nolint:goconst // testing
		t.Fatalf("bad string version: %s", version.String())
	}
}

// nolint:errcheck // testing
func TestSemanticVersionCompare(t *testing.T) {
	v1, _ := NewSemanticVersion("1.2.3")
	v2, _ := NewSemanticVersion("1.2.3")
	r, err := v1.Compare(v2)
	if err != nil {
		t.Fatalf("performing comparison: %d", err)
	}
	if r != 0 {
		t.Fatalf("wrong comparison: %d", r)
	}
	v2, _ = NewSemanticVersion("2.2.3")
	r, err = v1.Compare(v2)
	if err != nil {
		t.Fatalf("performing comparison: %d", err)
	}
	if r != -1 {
		t.Fatalf("wrong comparison: %d", r)
	}
	v2, _ = NewSemanticVersion("0.2.3")
	r, err = v1.Compare(v2)
	if err != nil {
		t.Fatalf("performing comparison: %d", err)
	}
	if r != 1 {
		t.Fatalf("wrong comparison: %d", r)
	}
	v2, _ = NewSemanticVersion("1.3.3")
	r, err = v1.Compare(v2)
	if err != nil {
		t.Fatalf("performing comparison: %d", err)
	}
	if r != -1 {
		t.Fatalf("wrong comparison: %d", r)
	}
	v2, _ = NewSemanticVersion("1.1.3")
	r, err = v1.Compare(v2)
	if err != nil {
		t.Fatalf("performing comparison: %d", err)
	}
	if r != 1 {
		t.Fatalf("wrong comparison: %d", r)
	}
	v2, _ = NewSemanticVersion("1.2.4")
	r, err = v1.Compare(v2)
	if err != nil {
		t.Fatalf("performing comparison: %d", err)
	}
	if r != -1 {
		t.Fatalf("wrong comparison: %d", r)
	}
	v2, _ = NewSemanticVersion("1.2.2")
	r, err = v1.Compare(v2)
	if err != nil {
		t.Fatalf("performing comparison: %d", err)
	}
	if r != 1 {
		t.Fatalf("wrong comparison: %d", r)
	}
}
