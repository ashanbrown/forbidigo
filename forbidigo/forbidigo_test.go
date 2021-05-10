package forbidigo

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForbiddenIdentifiers(t *testing.T) {
	t.Run("it finds forbidden identifiers", func(t *testing.T) {
		linter, _ := NewLinter([]string{`fmt\.Printf`})
		expectIssues(t, linter, `
package bar

func foo() {
	fmt.Printf("here i am")
}`, "use of `fmt.Printf` forbidden by pattern `fmt\\.Printf` at testing.go:5:2")
	})

	t.Run("it doesn't require a package on the identifier", func(t *testing.T) {
		linter, _ := NewLinter([]string{`Printf`})
		expectIssues(t, linter, `
package bar

func foo() {
	Printf("here i am")
}`, "use of `Printf` forbidden by pattern `Printf` at testing.go:5:2")
	})

	t.Run("allows explicitly permitting otherwise forbidden identifiers", func(t *testing.T) {
		linter, _ := NewLinter([]string{`fmt\.Printf`})
		expectIssues(t, linter, `
package bar

func foo() {
	fmt.Printf("here i am") //permit:fmt.Printf
}`)
	})

	t.Run("allows old notation for explicitly permitting otherwise forbidden identifiers", func(t *testing.T) {
		linter, _ := NewLinter([]string{`fmt\.Printf`})
		expectIssues(t, linter, `
package bar

func foo() {
	fmt.Printf("here i am") // permit:fmt.Printf
}`)
	})

	t.Run("has option to ignore permit directives", func(t *testing.T) {
		linter, _ := NewLinter([]string{`fmt\.Printf`}, OptionIgnorePermitDirectives(true))
		issues := parseFile(t, linter, "file.go", `
package bar

func foo() {
	fmt.Printf("here i am") //permit:fmt.Printf
}`)
		assert.NotEmpty(t, issues)
	})

	t.Run("examples are excluded by default in test files", func(t *testing.T) {
		linter, _ := NewLinter([]string{`fmt\.Printf`})
		issues := parseFile(t, linter, "file_test.go", `
package bar

func ExampleFoo() {
	fmt.Printf("here i am")
}`)
		assert.Empty(t, issues)
	})

	t.Run("whole file examples are excluded by default", func(t *testing.T) {
		linter, _ := NewLinter([]string{`fmt\.Printf`})
		issues := parseFile(t, linter, "file_test.go", `
package bar

func Foo() {
	fmt.Printf("here i am")
}

func Example() {
	Foo()
}`)
		assert.Empty(t, issues)
	})

	t.Run("Test functions prevent a file from being considered a whole file example", func(t *testing.T) {
		linter, _ := NewLinter([]string{`fmt\.Printf`})
		issues := parseFile(t, linter, "file_test.go", `
package bar

func TestFoo() {
	fmt.Printf("here i am")
}

func Example() {
}`)
		assert.NotEmpty(t, issues)
	})

	t.Run("Benchmark functions prevent a file from being considered a whole file example", func(t *testing.T) {
		linter, _ := NewLinter([]string{`fmt\.Printf`})
		issues := parseFile(t, linter, "file_test.go", `
package bar

func BenchmarkFoo() {
	fmt.Printf("here i am")
}

func Example() {
}`)
		assert.NotEmpty(t, issues)
	})

	t.Run("examples can be included", func(t *testing.T) {
		linter, _ := NewLinter([]string{`fmt\.Printf`}, OptionExcludeGodocExamples(false))
		issues := parseFile(t, linter, "file.go", `
package bar

func ExampleFoo() {
	fmt.Printf("here i am")
}`)
		assert.NotEmpty(t, issues)
	})
}

func expectIssues(t *testing.T, linter *Linter, contents string, issues ...string) {
	actualIssues := parseFile(t, linter, "testing.go", contents)
	actualIssueStrs := make([]string, 0, len(actualIssues))
	for _, i := range actualIssues {
		actualIssueStrs = append(actualIssueStrs, i.String())
	}
	assert.ElementsMatch(t, issues, actualIssueStrs)
}

func parseFile(t *testing.T, linter *Linter, fileName, contents string) []Issue {
	fset := token.NewFileSet()
	expr, err := parser.ParseFile(fset, fileName, contents, parser.ParseComments)
	if err != nil {
		t.Fatalf("unable to parse file contents: %s", err)
	}
	issues, err := linter.Run(fset, expr)
	if err != nil {
		t.Fatalf("unable to parse file: %s", err)
	}
	return issues
}
