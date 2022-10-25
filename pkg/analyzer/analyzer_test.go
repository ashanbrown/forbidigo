package analyzer_test

import (
	"testing"

	"github.com/ashanbrown/forbidigo/forbidigo"
	"github.com/ashanbrown/forbidigo/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	patterns := append(forbidigo.DefaultPatterns(),
		`^(?P<pkg>example.com/some/pkg)\.Forbidden$`,
		`^(?P<pkg>example.com/some/pkg.CustomType)\.AlsoForbidden$`,
		`^(?P<pkg>example.com/some/pkg.CustomInterface)\.StillForbidden$`,
	)
	a := analyzer.NewAnalyzer()
	for _, pattern := range patterns {
		if err := a.Flags.Set("p", pattern); err != nil {
			t.Fatalf("unexpected error when setting pattern: %v", err)
		}
	}
	analysistest.Run(t, testdata, a, "")
}
