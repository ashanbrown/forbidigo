package main

import (
	"flag"
	"go/ast"
	"log"
	"os"

	"github.com/ashanbrown/forbidigo/forbidigo"
	"golang.org/x/tools/go/packages"
)

func main() {
	log.SetFlags(0) // remove log timestamp

	setExitStatus := flag.Bool("set_exit_status", false, "Set exit status to 1 if any issues are found")
	includeTests := flag.Bool("tests", true, "Include tests")
	excludeGodocExamples := flag.Bool("exclude_godoc_examples", true, "Exclude code in godoc examples")
	analyzeTypes := flag.Bool("analyze_types", false, "Replace the literal source code based on the semantic of the code before matching against patterns")
	flag.Parse()

	var patterns = []string(nil)

	firstPkg := 0
	for n, arg := range flag.Args() {
		if arg == "--" {
			firstPkg = n + 1
			break
		}
		patterns = append(patterns, arg)
	}

	if patterns == nil {
		patterns = forbidigo.DefaultPatterns()
	}
	options := []forbidigo.Option{
		forbidigo.OptionExcludeGodocExamples(*excludeGodocExamples),
		forbidigo.OptionAnalyzeTypes(*analyzeTypes),
	}
	linter, err := forbidigo.NewLinter(patterns, options...)
	if err != nil {
		log.Fatalf("Could not create linter: %s", err)
	}

	cfg := packages.Config{
		Mode:  packages.NeedSyntax | packages.NeedName | packages.NeedFiles | packages.NeedTypes,
		Tests: *includeTests,
	}

	if *analyzeTypes {
		cfg.Mode |= packages.NeedTypesInfo | packages.NeedDeps
	}

	pkgs, err := packages.Load(&cfg, flag.Args()[firstPkg:]...)
	if err != nil {
		log.Fatalf("Could not load packages: %s", err)
	}

	var issues []forbidigo.Issue
	for _, p := range pkgs {
		nodes := make([]ast.Node, 0, len(p.Syntax))
		for _, n := range p.Syntax {
			nodes = append(nodes, n)
		}
		newIssues, err := linter.RunWithConfig(forbidigo.RunConfig{Fset: p.Fset, TypesInfo: p.TypesInfo}, nodes...)
		if err != nil {
			log.Fatalf("failed: %s", err)
		}
		issues = append(issues, newIssues...)
	}

	for _, issue := range issues {
		log.Println(issue)
	}

	if *setExitStatus && len(issues) > 0 {
		os.Exit(1)
	}
}
