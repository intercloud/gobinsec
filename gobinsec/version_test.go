package gobinsec

import "testing"

func TestNewVersionSemantic(t *testing.T) {
	if NewVersion("1.2.3").String() != "1.2.3" {
		t.Fatalf("bad semantic version: %s", NewVersion("1.2.3").String())
	}
}

func TestNewVersionPseudo(t *testing.T) {
	if NewVersion("v0.0.0-20191109021931-daa7c04131f5").String() != "v0.0.0-20191109021931-daa7c04131f5" {
		t.Fatalf("bad pseudo version: %s", NewVersion("v0.0.0-20191109021931-daa7c04131f5").String())
	}
}

func TestNewVersionDate(t *testing.T) {
	if NewVersion("2021-12-08").String() != "2021-12-08" {
		t.Fatalf("bad date version: %s", NewVersion("2021-12-08").String())
	}
}

func TestCompareVersionSemantic(t *testing.T) {
	_, err := NewVersion("1.2.3").Compare(NewVersion("2021-12-08"))
	if err == nil {
		t.Fatalf("should have produced an error comparing semantic version with date version")
	}
	if err.Error() != `can't compare semantic version to other type: *gobinsec.DateVersion` {
		t.Fatalf("bad error message: %s", err.Error())
	}
}

func TestCompareVersionPseudo(t *testing.T) {
	r, err := NewVersion("v0.0.0-20191109021931-daa7c04131f5").Compare(NewVersion("2019-11-10"))
	if err != nil {
		t.Fatalf("error comparing pseudo version with date version")
	}
	if r != -1 {
		t.Fatalf("bad comparison result: %d", r)
	}
	r, err = NewVersion("v0.0.0-20191109021931-daa7c04131f5").Compare(NewVersion("2019-11-08"))
	if err != nil {
		t.Fatalf("error comparing pseudo version with date version")
	}
	if r != 1 {
		t.Fatalf("bad comparison result: %d", r)
	}
	r, err = NewVersion("v0.0.0-20191109021931-daa7c04131f5").Compare(NewVersion("2019-11-09"))
	if err != nil {
		t.Fatalf("error comparing pseudo version with date version")
	}
	if r != 0 {
		t.Fatalf("bad comparison result: %d", r)
	}
}
