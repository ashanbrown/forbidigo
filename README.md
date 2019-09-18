# disallow

disallow is a Go static analysis tool to disallow use of particular identifiers.

## Installation

    go get -u github.com/ashanbrown/disallow

## Usage

    disallow [flags...] patterns... -- packages...

### Flags
- **-set_exit_status** (default false) - Set exit status to 1 if any issues are found.

## Purpose

To prevent leaving format statements and temporary statements such as Ginkgo FIt, FDescribe, etc.

## Ignoring issues

You can ignore a particular issue by including the directive `// allow` on that line

## Contributing

Pull requests welcome!
