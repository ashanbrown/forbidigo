// forbidigo provides a linter for forbidding the use of specific identifiers
package forbidigo

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"log"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type Issue interface {
	Details() string
	Pos() token.Pos
	Position() token.Position
	String() string
}

type UsedIssue struct {
	identifier string
	pattern    string
	pos        token.Pos
	position   token.Position
	customMsg  string
}

func (a UsedIssue) Details() string {
	explanation := fmt.Sprintf(` because %q`, a.customMsg)
	if a.customMsg == "" {
		explanation = fmt.Sprintf(" by pattern `%s`", a.pattern)
	}
	return fmt.Sprintf("use of `%s` forbidden", a.identifier) + explanation
}

func (a UsedIssue) Position() token.Position {
	return a.position
}

func (a UsedIssue) Pos() token.Pos {
	return a.pos
}

func (a UsedIssue) String() string { return toString(a) }

func toString(i UsedIssue) string {
	return fmt.Sprintf("%s at %s", i.Details(), i.Position())
}

type Linter struct {
	cfg      config
	patterns []*Pattern
}

func DefaultPatterns() []string {
	return []string{`^(fmt\.Print(|f|ln)|print|println)$`}
}

//go:generate go-options config
type config struct {
	// don't check inside Godoc examples (see https://blog.golang.org/examples)
	ExcludeGodocExamples   bool `options:",true"`
	IgnorePermitDirectives bool // don't check for `permit` directives(for example, in favor of `nolint`)
}

func NewLinter(patterns []string, options ...Option) (*Linter, error) {
	cfg, err := newConfig(options...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process options")
	}

	if len(patterns) == 0 {
		patterns = DefaultPatterns()
	}
	compiledPatterns := make([]*Pattern, 0, len(patterns))
	for _, ptrn := range patterns {
		p, err := parse(ptrn)
		if err != nil {
			return nil, err
		}
		compiledPatterns = append(compiledPatterns, p)
	}
	return &Linter{
		cfg:      cfg,
		patterns: compiledPatterns,
	}, nil
}

type visitor struct {
	cfg        config
	isTestFile bool // godoc only runs on test files

	linter   *Linter
	comments []*ast.CommentGroup

	runConfig RunConfig
	issues    []Issue
}

// Deprecated: Run was the original entrypoint before RunWithConfig was introduced to support
// additional match patterns that need additional information.
func (l *Linter) Run(fset *token.FileSet, nodes ...ast.Node) ([]Issue, error) {
	return l.RunWithConfig(RunConfig{Fset: fset}, nodes...)
}

// RunConfig provides information that the linter needs for different kinds
// of match patterns. Ideally, all fields should get set. More fields may get
// added in the future as needed.
type RunConfig struct {
	// FSet is required.
	Fset *token.FileSet

	// TypesInfo is needed for "pkg" match patterns. Not providing it
	// disables those patterns.
	TypesInfo *types.Info

	// DebugLog is used to print debug messages. May be nil.
	DebugLog func(format string, args ...interface{})
}

// Patterns returns the parsed patterns.
func (l *Linter) Patterns() []*Pattern {
	return l.patterns
}

func (l *Linter) RunWithConfig(config RunConfig, nodes ...ast.Node) ([]Issue, error) {
	if config.DebugLog == nil {
		config.DebugLog = func(format string, args ...interface{}) {}
	}
	var issues []Issue
	for _, node := range nodes {
		var comments []*ast.CommentGroup
		isTestFile := false
		isWholeFileExample := false
		if file, ok := node.(*ast.File); ok {
			comments = file.Comments
			fileName := config.Fset.Position(file.Pos()).Filename
			isTestFile = strings.HasSuffix(fileName, "_test.go")

			// From https://blog.golang.org/examples, a "whole file example" is:
			// a file that ends in _test.go and contains exactly one example function,
			// no test or benchmark functions, and at least one other package-level declaration.
			if l.cfg.ExcludeGodocExamples && isTestFile && len(file.Decls) > 1 {
				numExamples := 0
				numTestsAndBenchmarks := 0
				for _, decl := range file.Decls {
					funcDecl, isFuncDecl := decl.(*ast.FuncDecl)
					// consider only functions, not methods
					if !isFuncDecl || funcDecl.Recv != nil || funcDecl.Name == nil {
						continue
					}
					funcName := funcDecl.Name.Name
					if strings.HasPrefix(funcName, "Test") || strings.HasPrefix(funcName, "Benchmark") {
						numTestsAndBenchmarks++
						break // not a whole file example
					}
					if strings.HasPrefix(funcName, "Example") {
						numExamples++
					}
				}

				// if this is a whole file example, skip this node
				isWholeFileExample = numExamples == 1 && numTestsAndBenchmarks == 0
			}
		}
		if isWholeFileExample {
			continue
		}
		visitor := visitor{
			cfg:        l.cfg,
			isTestFile: isTestFile,
			linter:     l,
			runConfig:  config,
			comments:   comments,
		}
		ast.Walk(&visitor, node)
		issues = append(issues, visitor.issues...)
	}
	return issues, nil
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case *ast.FuncDecl:
		// don't descend into godoc examples if we are ignoring them
		isGodocExample := v.isTestFile && node.Recv == nil && node.Name != nil && strings.HasPrefix(node.Name.Name, "Example")
		if isGodocExample && v.cfg.ExcludeGodocExamples {
			return nil
		}
		return v
	// The following two are handled below.
	case *ast.SelectorExpr:
	case *ast.Ident:
	// Everything else isn't.
	default:
		return v
	}

	// The text as it appears in the source is always used because issues
	// use that. The other texts to match against are extracted when needed
	// by a pattern. They are nil when not evaluated yet and point to
	// the empty string when there is nothing to match against.
	srcText := v.textFor(node)
	var pkgText *string
	checkedPkgText := false
	for _, p := range v.linter.patterns {
		if p.Match == MatchType && !checkedPkgText {
			pkgText = v.pkgTextFor(node)
			if pkgText != nil {
				v.runConfig.DebugLog("%s: %q -> %q", v.runConfig.Fset.Position(node.Pos()), srcText, *pkgText)
			} else {
				v.runConfig.DebugLog("%s: %q -> not expanded", v.runConfig.Fset.Position(node.Pos()), srcText)
			}
			checkedPkgText = true
		}

		matchText := ""
		switch {
		case p.Match == MatchType:
			if pkgText == nil {
				continue
			}
			matchText = *pkgText
		default:
			matchText = srcText
		}
		if p.re.MatchString(matchText) && !v.permit(node) {
			v.issues = append(v.issues, UsedIssue{
				identifier: srcText, // Always report the expression as it appears in the source code.
				pattern:    p.re.String(),
				pos:        node.Pos(),
				position:   v.runConfig.Fset.Position(node.Pos()),
				customMsg:  p.Msg,
			})
		}
	}
	return nil
}

// textFor returns the expression as it appears in the source code (for
// example, <importname>.<function name>).
func (v *visitor) textFor(node ast.Node) string {
	buf := new(bytes.Buffer)
	if err := printer.Fprint(buf, v.runConfig.Fset, node); err != nil {
		log.Fatalf("ERROR: unable to print node at %s: %s", v.runConfig.Fset.Position(node.Pos()), err)
	}
	return buf.String()
}

// pkgTextFor expands the selector in a selector expression to the full package
// name and (for variables) the type:
//
// - example.com/some/pkg.Function
// - example.com/some/pkg.CustomType.Method
//
// It returns nil when the text is not available, otherwise a pointer to
// the string to match against.
func (v *visitor) pkgTextFor(node ast.Node) *string {
	if v.runConfig.TypesInfo == nil {
		return nil
	}

	// TODO: do type switch here instead of multiple if checks.
	if ident, ok := node.(*ast.Ident); ok {
		object, ok := v.runConfig.TypesInfo.Uses[ident]
		if !ok {
			// No information about the identifier. Should
			// not happen, but perhaps there were compile
			// errors?
			return nil
		}
		str := object.Name()
		if pkg := object.Pkg(); pkg != nil {
			str = pkg.Path() + "." + str
		}
		return &str
	}
	selectorExpr, ok := node.(*ast.SelectorExpr)
	if !ok {

		return nil
	}
	selector := selectorExpr.X
	var pkgText string

	// If we are lucky, the entire selector expression has a known
	// type. We don't care about the value.
	if typeAndValue, ok := v.runConfig.TypesInfo.Types[selector]; ok {
		pkgText = typeAndValue.Type.String()
	} else {
		// Some expressions need special treatment.
		switch selector := selector.(type) {
		case *ast.Ident:
			object, ok := v.runConfig.TypesInfo.Uses[selector]
			if !ok {
				// No information about the identifier. Should
				// not happen, but perhaps there were compile
				// errors?
				return nil
			}
			switch object := object.(type) {
			case *types.PkgName:
				pkgText = object.Imported().Path()
			case *types.Var:
				pkgText = object.Type().String()
			default:
				// Something else?
				return nil
			}
		default:
			return nil
		}
	}

	// Ignore whether it is a pointer.
	pkgText = strings.TrimLeft(pkgText, "*")
	pkgText += "." + selectorExpr.Sel.Name
	return &pkgText
}

func (v *visitor) permit(node ast.Node) bool {
	if v.cfg.IgnorePermitDirectives {
		return false
	}
	nodePos := v.runConfig.Fset.Position(node.Pos())
	nolint := regexp.MustCompile(fmt.Sprintf(`^//\s?permit:%s\b`, regexp.QuoteMeta(v.textFor(node))))
	for _, c := range v.comments {
		commentPos := v.runConfig.Fset.Position(c.Pos())
		if commentPos.Line == nodePos.Line && len(c.List) > 0 && nolint.MatchString(c.List[0].Text) {
			return true
		}
	}
	return false
}
