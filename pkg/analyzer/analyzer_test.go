package analyzer_test

import (
	"testing"

	"github.com/ashanbrown/forbidigo/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	a := analyzer.NewAnalyzer()
	analysistest.Run(t, testdata, a, "")
}
