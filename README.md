# forbidigo

forbidigo is a Go static analysis tool to forbidigo use of particular identifiers.

## Installation

    go get -u github.com/ashanbrown/forbidigo

## Usage

    forbidigo [flags...] patterns... -- packages...

If no patterns are specified, the default pattern of `fmt\.Printf.*` is used to eliminate debug statememts.

A larger set of interesting patterns might include:

* `fmt\.Printf.*` -- forbid use of Printf because it is likely just for debugging
* `fmt\.Errorf` -- forbid Errorf in favor of using github.com/pkg/errors
* `ginkgo\.F.*` -- forbid ginkgo focused commands (used for debug issues)
* `spew\.Dump` -- forbid dumping detailed data to stdout

Note that the linter has no knowledge of what packages were actually imported, so aliased imports will match these patterns.

### Flags
- **-set_exit_status** (default false) - Set exit status to 1 if any issues are found.

## Purpose

To prevent leaving format statements and temporary statements such as Ginkgo FIt, FDescribe, etc.

## Ignoring issues

You can ignore a particular issue by including the directive `// permit` on that line

## Contributing

Pull requests welcome!
