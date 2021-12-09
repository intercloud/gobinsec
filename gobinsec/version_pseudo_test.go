package gobinsec

import "testing"

func TestNewPseudoVersion(t *testing.T) {
	pseudoVersion, err := NewPseudoVersion("v0.0.0-20191109021931-daa7c04131f5")
	if err != nil {
		t.Fatalf("parsing pseudo version: %v", err)
	}
	if pseudoVersion.Time.Format(PseudoVersionTimeFormat) != "20191109021931" {
		t.Fatalf("bad pseudo version time: %s", pseudoVersion.Time.Format(PseudoVersionTimeFormat))
	}
	if pseudoVersion.Commit != "daa7c04131f5" {
		t.Fatalf("bad pseudo version commit: %s", pseudoVersion.Commit)
	}
}

func TestPseudoVersionString(t *testing.T) {
	pseudoVersion, err := NewPseudoVersion("v0.0.0-20191109021931-daa7c04131f5")
	if err != nil {
		t.Fatalf("parsing pseudo version: %v", err)
	}
	if pseudoVersion.String() != "v0.0.0-20191109021931-daa7c04131f5" { // nolint:goconst // testing
		t.Fatalf("bad pseudo version string representation: %s", pseudoVersion.String())
	}
}

func TestPseudoVersionCompare(t *testing.T) {
	p1, _ := NewPseudoVersion("v0.0.0-20191109021931-daa7c04131f5") // nolint:errcheck // testing
	p2, _ := NewPseudoVersion("v0.0.0-20191109021931-daa7c04131f5") // nolint:errcheck // testing
	r, err := p1.Compare(p2)
	if err != nil {
		t.Fatalf("performing comparison: %v", err)
	}
	if r != 0 {
		t.Fatalf("bad comparison: %d", r)
	}
	p2, _ = NewPseudoVersion("v0.0.0-20191108021931-daa7c04131f5") // nolint:errcheck // testing
	r, err = p1.Compare(p2)
	if err != nil {
		t.Fatalf("performing comparison: %v", err)
	}
	if r != 1 {
		t.Fatalf("bad comparison: %d", r)
	}
	p2, _ = NewPseudoVersion("v0.0.0-20191110021931-daa7c04131f5") // nolint:errcheck // testing
	r, err = p1.Compare(p2)
	if err != nil {
		t.Fatalf("performing comparison: %v", err)
	}
	if r != -1 {
		t.Fatalf("bad comparison: %d", r)
	}
}
