package analyzer

import (
	"testing"

	"github.com/ashanbrown/forbidigo/forbidigo"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLiteralAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	patterns := append(forbidigo.DefaultPatterns(),
		`^pkg\.Forbidden$`,
		`^Shiny`,
		`^renamed\.Forbidden`,
	)
	a := newAnalyzer(t.Logf)
	for _, pattern := range patterns {
		if err := a.Flags.Set("p", pattern); err != nil {
			t.Fatalf("unexpected error when setting pattern: %v", err)
		}
	}
	analysistest.Run(t, testdata, a, "matchtext")
}

func TestExpandAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	patterns := append(forbidigo.DefaultPatterns(),
		`{pattern: ^pkg\.Forbidden$, package: ^example.com/some/pkg$}`,
		`{pattern: ^pkg\.CustomType.*Forbidden.*$, package: ^example.com/some/pkg$}`,
		`{pattern: ^pkg\.CustomInterface.*Forbidden$, package: ^example.com/some/pkg$}`,
		`{pattern: ^thing\.Shiny, package: ^example.com/some/thing$}`,
		`{pattern: myCustomStruct\..*Forbidden, package: ^expandtext$}`,
		`{pattern: myCustomInterface\.AlsoForbidden, package: ^expandtext$}`,
		`{pattern: renamedpkg\.Forbidden, package: ^example.com/some/renamedpkg$}`,
		`{pattern: renamedpkg\.Struct.Forbidden, package: ^example.com/some/renamedpkg$}`,
	)
	a := newAnalyzer(t.Logf)
	for _, pattern := range patterns {
		if err := a.Flags.Set("p", pattern); err != nil {
			t.Fatalf("unexpected error when setting pattern: %v", err)
		}
	}
	if err := a.Flags.Set("expand_expressions", "true"); err != nil {
		t.Fatalf("unexpected error when enabling expression expansion: %v", err)
	}
	analysistest.Run(t, testdata, a, "expandtext")
}
