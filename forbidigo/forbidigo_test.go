package forbidigo

import (
	"go/ast"
	"log"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
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

	t.Run("it finds forbidden, renamed identifiers", func(t *testing.T) {
		linter, _ := NewLinter([]string{`fmt\.Printf`})
		expectIssues(t, linter, `
package bar

import renamed "fmt"

func foo() {
	renamed.Printf("here i am")
}` /* only detected inside golanci-lint */)
	})

	t.Run("displays custom messages", func(t *testing.T) {
		linter, _ := NewLinter([]string{`^fmt\.Printf(# a custom message)?$`})
		expectIssues(t, linter, `
package bar

func foo() {
	fmt.Printf("here i am")
}`, "use of `fmt.Printf` forbidden because \"a custom message\" at testing.go:5:2")
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

	t.Run("import renames not detected by simple pattern", func(t *testing.T) {
		linter, _ := NewLinter([]string{`fmt\.Printf`}, OptionExcludeGodocExamples(false))
		issues := parseFile(t, linter, "file.go", `
package bar

import fmt2 "fmt"

func ExampleFoo() {
	fmt2.Printf("here i am")
}`)
		assert.Empty(t, issues)
	})

	t.Run("import renames detected by package pattern", func(t *testing.T) {
		linter, _ := NewLinter([]string{`(?P<pkg>fmt)\.Printf`}, OptionExcludeGodocExamples(false))
		expectIssues(t, linter, `
package bar

import fmt2 "fmt"

func ExampleFoo() {
	fmt2.Printf("here i am")
}`, "use of `fmt2.Printf` forbidden by pattern `(?P<pkg>fmt)\\.Printf` at testing.go:7:2")
	})

}

// sourcePath matches "at /tmp/TestForbiddenIdentifiersdisplays_custom_messages4260088387/001/testing.go".
var sourcePath = regexp.MustCompile(`at .*/([[:alnum:]]+.go)`)

func expectIssues(t *testing.T, linter *Linter, contents string, issues ...string) {
	actualIssues := parseFile(t, linter, "testing.go", contents)
	actualIssueStrs := make([]string, 0, len(actualIssues))
	for _, i := range actualIssues {
		str := i.String()
		str = sourcePath.ReplaceAllString(str, "at $1")
		actualIssueStrs = append(actualIssueStrs, str)
	}
	assert.ElementsMatch(t, issues, actualIssueStrs)
}

func parseFile(t *testing.T, linter *Linter, fileName, contents string) []Issue {
	// We can use packages.Load if we put a single file into a separate
	// directory and parse it with Go modules of. We have to be in that
	// directory to use "." as pattern, parsing it via the absolute path
	// from the forbidigo project doesn't work ("cannot import absolute
	// path").
	tmpDir := t.TempDir()
	if err := os.WriteFile(path.Join(tmpDir, fileName), []byte(contents), 0644); err != nil {
		t.Fatalf("could not write source file: %v", err)
	}
	env := os.Environ()
	env = append(env, "GO111MODULE=off")
	cfg := packages.Config{
		Mode:  packages.NeedSyntax | packages.NeedName | packages.NeedFiles | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps,
		Env:   env,
		Tests: true,
	}
	pwd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(pwd)
	err = os.Chdir(tmpDir)
	require.NoError(t, err)
	pkgs, err := packages.Load(&cfg, ".")
	if err != nil {
		t.Fatalf("could not load packages: %v", err)
	}
	var issues []Issue
	for _, p := range pkgs {
		nodes := make([]ast.Node, 0, len(p.Syntax))
		for _, n := range p.Syntax {
			nodes = append(nodes, n)
		}
		newIssues, err := linter.RunWithTypes(p.Fset, p.TypesInfo, nodes...)
		if err != nil {
			log.Fatalf("failed: %s", err)
		}
		issues = append(issues, newIssues...)
	}
	if err != nil {
		t.Fatalf("unable to parse file: %s", err)
	}
	return issues
}
