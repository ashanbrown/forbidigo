# forbidigo

forbidigo is a Go static analysis tool to forbidigo use of particular identifiers.

## Installation

    go get -u github.com/ashanbrown/forbidigo

## Usage

    forbidigo [flags...] patterns... -- packages...

If no patterns are specified, the default pattern of `^(fmt\.Print.*|print|println)$` is used to eliminate debug statements.  By default,
functions (and whole files), that are identifies as Godoc examples (https://blog.golang.org/examples) are excluded from 
checking.

A larger set of interesting patterns might include:

* `^fmt\.Print.*$` -- forbid use of Print statements because they are likely just for debugging
* `^fmt\.Errorf$` -- forbid Errorf in favor of using github.com/pkg/errors
* `^ginkgo\.F[A-Z].*$` -- forbid ginkgo focused commands (used for debug issues)
* `^spew\.Dump$` -- forbid dumping detailed data to stdout

Note that the linter has no knowledge of what packages were actually imported, so aliased imports will match these patterns.

### Flags
- **-set_exit_status** (default false) - Set exit status to 1 if any issues are found.
- **-exclude_godoc_examples** (default true) - Controls whether godoc examples are identified and excluded

## Purpose

To prevent leaving format statements and temporary statements such as Ginkgo FIt, FDescribe, etc.

## Ignoring issues

You can ignore a particular issue by including the directive `//permit` on that line

## Contributing

Pull requests welcome!
