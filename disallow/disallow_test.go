package disallow

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisallow(t *testing.T) {
	t.Run("it finds disallowed identifiers", func(t *testing.T) {
		linter, _ := NewLinter(`fmt\.Printf`)
		expectIssues(t, linter, `
package bar

func foo() {
	fmt.Printf("here i am")
}`, "use of `fmt.Printf` disallowed by pattern `fmt\\.Printf` at testing.go:5:2")
	})

	t.Run("allows explicitly allowing identifiers", func(t *testing.T) {
		linter, _ := NewLinter(`fmt\.Printf`)
		expectIssues(t, linter, `
package bar

func foo() {
	fmt.Printf("here i am") // allow:fmt.Printf
}`)
	})
}

func expectIssues(t *testing.T, linter *Linter, contents string, issues ...string) {
	actualIssues := parseFile(t, linter, contents)
	actualIssueStrs := make([]string, 0, len(actualIssues))
	for _, i := range actualIssues {
		actualIssueStrs = append(actualIssueStrs, i.String())
	}
	assert.ElementsMatch(t, issues, actualIssueStrs)
}

func parseFile(t *testing.T, linter *Linter, contents string) []Issue {
	fset := token.NewFileSet()
	expr, err := parser.ParseFile(fset, "testing.go", contents, parser.ParseComments)
	if err != nil {
		t.Fatalf("unable to parse file contents: %s", err)
	}
	issues, err := linter.Run(fset, expr)
	if err != nil {
		t.Fatalf("unable to parse file: %s", err)
	}
	return issues
}
