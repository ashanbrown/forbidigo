# forbidigo

forbidigo is a Go static analysis tool to forbidigo use of particular identifiers.

## Installation

    go get -u github.com/ashanbrown/forbidigo

## Usage

    forbidigo [flags...] patterns... -- packages...

Some example patterns would be:

    fmt\.Printf.* -- forbid use of Printf because it is likely just for debugging
    fmt\.Errorf -- forbid Errorf in favor of using github.com/pkg/errors
    ginkgo\.F.* -- forbid ginkgo focused commands (used for debug issues)

### Flags
- **-set_exit_status** (default false) - Set exit status to 1 if any issues are found.

## Purpose

To prevent leaving format statements and temporary statements such as Ginkgo FIt, FDescribe, etc.

## Ignoring issues

You can ignore a particular issue by including the directive `// permit` on that line

## Contributing

Pull requests welcome!
