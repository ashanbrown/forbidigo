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
		`{Match: type, Pattern: ^example.com/some/pkg\.Forbidden$}`,
		`{Match: type, Pattern: ^example.com/some/pkg.CustomType\..*Forbidden.*$}`,
		`{Match: type, Pattern: ^example.com/some/pkg.CustomInterface\.StillForbidden$}`,
		`{Match: type, Pattern: example.com/some/thing\.Shiny}`,
		`{Match: type, Pattern: myCustomStruct\..*Forbidden}`,
		`{Match: type, Pattern: myCustomInterface\..*Forbidden}`,
		`{Match: type, Pattern: renamedpkg\.Forbidden}`,
		`{Match: type, Pattern: renamedpkg\.Struct.Forbidden}`,
	)
	a := newAnalyzer(t.Logf)
	for _, pattern := range patterns {
		if err := a.Flags.Set("p", pattern); err != nil {
			t.Fatalf("unexpected error when setting pattern: %v", err)
		}
	}
	analysistest.Run(t, testdata, a, "expandtext")
}
