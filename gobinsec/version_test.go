package gobinsec

import "testing"

func TestNewVersionSemantic(t *testing.T) {
	v, err := NewVersion("1.2.3")
	if err != nil {
		t.Fatalf("error parsing semantic version: %v", err)
	}
	if v.String() != "1.2.3" {
		t.Fatalf("bad semantic version: %s", v.String())
	}
}

func TestNewVersionPseudo(t *testing.T) {
	v, err := NewVersion("v0.0.0-20191109021931-daa7c04131f5")
	if err != nil {
		t.Fatalf("error parsing pseudo version: %v", err)
	}
	if v.String() != "v0.0.0-20191109021931-daa7c04131f5" {
		t.Fatalf("bad pseudo version: %s", v.String())
	}
}

func TestNewVersionDate(t *testing.T) {
	v, err := NewVersion("2021-12-08")
	if err != nil {
		t.Fatalf("error parsing date version: %v", err)
	}
	if v.String() != "2021-12-08" {
		t.Fatalf("bad date version: %s", v.String())
	}
}

// nolint:errcheck // testing
func TestCompareVersionSemantic(t *testing.T) {
	v1, _ := NewVersion("1.2.3")
	v2, _ := NewVersion("2021-12-08")
	_, err := v1.Compare(v2)
	if err == nil {
		t.Fatalf("should have produced an error comparing semantic version with date version")
	}
	if err.Error() != `can't compare semantic version to other type: *gobinsec.DateVersion` {
		t.Fatalf("bad error message: %s", err.Error())
	}
}

// nolint:errcheck // testing
func TestCompareVersionPseudo(t *testing.T) {
	v1, _ := NewVersion("v0.0.0-20191109021931-daa7c04131f5")
	v2, _ := NewVersion("2019-11-10")
	r, err := v1.Compare(v2)
	if err != nil {
		t.Fatalf("error comparing pseudo version with date version")
	}
	if r != -1 {
		t.Fatalf("bad comparison result: %d", r)
	}
	v2, _ = NewVersion("2019-11-08")
	r, err = v1.Compare(v2)
	if err != nil {
		t.Fatalf("error comparing pseudo version with date version")
	}
	if r != 1 {
		t.Fatalf("bad comparison result: %d", r)
	}
	v2, _ = NewVersion("2019-11-09")
	r, err = v1.Compare(v2)
	if err != nil {
		t.Fatalf("error comparing pseudo version with date version")
	}
	if r != 0 {
		t.Fatalf("bad comparison result: %d", r)
	}
}
