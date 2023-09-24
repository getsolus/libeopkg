package pspec

import (
	"testing"
)

const (
	pspec = "../testdata/pspec_x86_64.xml"
)

func TestParse(t *testing.T) {
	p, err := Load(pspec)
	if err != nil {
		t.Fatalf("Should have loaded successfully: %s", err)
	}

	if len(p.Packages) != 1 {
		t.Fatalf("Should have exactly one package")
	}
	if fileCnt := len(p.Packages[0].Files); fileCnt != 94 {
		t.Fatalf("Should have exactly 94 files, got %d", fileCnt)
	}
	if fstFile := p.Packages[0].Files[0].Value; fstFile != "/usr/bin/nano" {
		t.Fatalf("Should have /usr/bin/nano as the first file in the package, but got %s", fstFile)
	}
}
