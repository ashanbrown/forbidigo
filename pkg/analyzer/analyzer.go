package analyzer

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"

	"github.com/ashanbrown/forbidigo/forbidigo"
	"golang.org/x/tools/go/analysis"
)

type listVar struct {
	values *[]string
}

func (v *listVar) Set(value string) error {
	*v.values = append(*v.values, value)
	if value == "" {
		return errors.New("value cannot be empty")
	}
	return nil
}

func (v *listVar) String() string {
	return ""
}

type analyzer struct {
	patterns           []string
	usePermitDirective bool
	includeExamples    bool
	analyzeTypes       bool
	debugLog           func(format string, args ...interface{})
}

// NewAnalyzer returns a go/analysis-compatible analyzer
// The "-p" argument can be used to add a pattern.
// Set "-examples" to analyze godoc examples
// Set "-permit=false" to ignore "//permit:<identifier>" directives.
func NewAnalyzer() *analysis.Analyzer {
	return newAnalyzer(nil /* no debug output */)
}

func newAnalyzer(debugLog func(format string, args ...interface{})) *analysis.Analyzer {
	var flags flag.FlagSet
	a := analyzer{
		usePermitDirective: true,
		includeExamples:    true,
		debugLog:           debugLog,
	}

	flags.Var(&listVar{values: &a.patterns}, "p", "pattern")
	flags.BoolVar(&a.includeExamples, "examples", false, "check godoc examples")
	flags.BoolVar(&a.usePermitDirective, "permit", true, `when set, lines with "//permit" directives will be ignored`)
	flags.BoolVar(&a.analyzeTypes, "analyze_types", false, `when set, expressions get expanded instead of matching the literal source code`)
	return &analysis.Analyzer{
		Name:  "forbidigo",
		Doc:   "forbid identifiers",
		Run:   a.runAnalysis,
		Flags: flags,
	}
}

func (a *analyzer) runAnalysis(pass *analysis.Pass) (interface{}, error) {
	if a.patterns == nil {
		a.patterns = forbidigo.DefaultPatterns()
	}
	linter, err := forbidigo.NewLinter(a.patterns,
		forbidigo.OptionIgnorePermitDirectives(!a.usePermitDirective),
		forbidigo.OptionExcludeGodocExamples(!a.includeExamples),
		forbidigo.OptionAnalyzeTypes(a.analyzeTypes),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to configure linter: %w", err)
	}
	nodes := make([]ast.Node, 0, len(pass.Files))
	for _, f := range pass.Files {
		nodes = append(nodes, f)
	}
	config := forbidigo.RunConfig{Fset: pass.Fset, DebugLog: a.debugLog}
	if a.analyzeTypes {
		config.TypesInfo = pass.TypesInfo
	}
	issues, err := linter.RunWithConfig(config, nodes...)
	if err != nil {
		return nil, err
	}
	reportIssues(pass, issues)
	return nil, nil
}

func reportIssues(pass *analysis.Pass, issues []forbidigo.Issue) {
	for _, i := range issues {
		diag := analysis.Diagnostic{
			Pos:      i.Pos(),
			Message:  i.Details(),
			Category: "restriction",
		}
		pass.Report(diag)
	}
}
