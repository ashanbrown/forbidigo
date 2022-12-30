# forbidigo

`forbidigo` is a Go static analysis tool to forbidigo use of particular identifiers.

`forbidigo` is recommended to be run as part of [golangci-lint](https://github.com/golangci/golangci-lint) where it can be controlled using file-based configuration and `//nolint` directives, but it can also be run as a standalone tool.

## Installation

    go get -u github.com/ashanbrown/forbidigo

## Usage

    forbidigo [flags...] patterns... -- packages...

If no patterns are specified, the default pattern of `^(fmt\.Print.*|print|println)$` is used to eliminate debug statements.  By default,
functions (and whole files), that are identifies as Godoc examples (https://blog.golang.org/examples) are excluded from 
checking.

By default, patterns get matched against the actual expression as it appears in
the source code. The effect is that ``^fmt\.Print.*$` will not match when that
package gets imported with `import fmt2 "fmt"` and then the function gets
called with `fmt2.Print`.

This makes it hard to match packages that may get imported under a variety of
different names, for example because there is no established convention or the
name is so generic that import aliases have to be used. To solve this,
forbidigo also supports more complex patterns. Such patterns are strings that
contain JSON or YAML for a struct.

The full pattern struct has the following fields:

* `Msg`: an additional comment that gets added to the error message when a
  pattern matches.
* `Pattern`: the regular expression itself.
* `Match`: a string which defines what the regular expression is matched
  against. Valid values are:
  * `text`: the traditional, literal source code match.
  * `type`: a semantic match that uses type information to enable precise
    matches against what is being used (a function in a certain package, or a
    method in a certain type) instead of how that thing is called in the source
    code.

A pattern with `Match: type` expands selector expressions (`<some>.<thing>`)
and identifiers. Those expressions get expanded as follows:

* An imported package gets replaced with the full package path, including the
  version if there is one. Example: `ginkgo.FIt` ->
  `github.com/onsi/ginkgo/v2.FIt`.

* For a method call, the type is inserted. Pointers are treated like the type
  they point to. When a type is an alias for a type in some other package, the
  name of that other package will be used. Example:

     var cf *spew.ConfigState = ...
     cf.Dump() // -> github.com/davecgh/go-spew/spew.ConfigState.Dump

* A simple identifier gets replaced with full package path and name. Example:

     . "github.com/onsi/ginkgo/v2"

     FIt(...) // -> github.com/onsi/ginkgo/v2.FIt

To distinguish such patterns from traditional regular expression patterns, the
encoding must start with a `{` or contain line breaks. When using just JSON
encoding, backslashes must get quoted inside strings. When using YAML, this
isn't necessary. The following pattern strings are equivalent:

    {Match: "type", Pattern: "^fmt\\.Println$"}

    {Match: type,
    Pattern: ^fmt\.Println$
    }

    {Match: type, Pattern: ^fmt\.Println$}

    Match: type
    Pattern: ^fmt\.Println$

A larger set of interesting patterns might include:

* `{Match: "type", Pattern: "^fmt\\.Print.*$"}` -- forbid use of Print statements because they are likely just for debugging
* `{Match: "type", Pattern: "^fmt\\.Errorf$", Msg: "use github.com/pkg/errors"}` -- forbid Errorf in favor of using github.com/pkg/errors
* `{Match: "type", Pattern: "^github.com/onsi/ginkgo(/v[[:digit:]]*)?)\\.F[A-Z].*$"}` -- forbid ginkgo focused commands (used for debug issues)
* `{Match: "type", Pattern: "^github.com/davecgh/go-spew/spew\\.Dump$"}` -- forbid dumping detailed data to stdout
* `{Match: "type", Pattern: "^github.com/davecgh/go-spew/spew.ConfigState\\.Dump"}` -- also forbid it via a `ConfigState`

For backwards compatibility, the message may also get encoded inside the
regular expression:

* `{Match: "type", Pattern: "^fmt\\.Errorf(# please use github\.com/pkg/errors)?$` -- forbid Errorf, with a custom message

### Flags
- **-set_exit_status** (default false) - Set exit status to 1 if any issues are found.
- **-exclude_godoc_examples** (default true) - Controls whether godoc examples are identified and excluded
- **-tests** (default true) - Controls whether tests are included

## Purpose

To prevent leaving format statements and temporary statements such as Ginkgo FIt, FDescribe, etc.

## Ignoring issues

You can ignore a particular issue by including the directive `//permit` on that line.  *This feature is disabled inside `golangci-lint` to encourage ignoring issues using the `// nolint` directive common for all linters (nolinting well is hard and I didn't want to make an effort do it exactly right within this linter).*

## Contributing

Pull requests welcome!
