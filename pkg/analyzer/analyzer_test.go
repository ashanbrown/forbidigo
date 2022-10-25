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
		`^AlsoShiny`,
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
		`{p: ^pkg\.Forbidden$, pkg: ^example.com/some/pkg$}`,
		`{p: ^pkg\.CustomType.*Forbidden.*$, pkg: ^example.com/some/pkg$}`,
		`{p: ^pkg\.CustomInterface.*Forbidden$, pkg: ^example.com/some/pkg$}`,
		`{p: ^Shiny, pkg: ^example.com/some/thing$}`,
		`{p: ^AlsoShiny}`,
		`{p: myCustomStruct\..*Forbidden, pkg: ^expandtext$}`,
		`{p: myCustomInterface\.AlsoForbidden, pkg: ^expandtext$}`,
		`{p: renamed\.Forbidden, pkg: ^example.com/some/renamedpkg$}`,
		`{p: renamed\.Struct.Forbidden, pkg: ^example.com/some/renamedpkg$}`,
	)
	a := newAnalyzer(t.Logf)
	for _, pattern := range patterns {
		if err := a.Flags.Set("p", pattern); err != nil {
			t.Fatalf("unexpected error when setting pattern: %v", err)
		}
	}
	if err := a.Flags.Set("analyze_types", "true"); err != nil {
		t.Fatalf("unexpected error when enabling expression expansion: %v", err)
	}
	analysistest.Run(t, testdata, a, "expandtext")
}
