# forbidigo

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/ashanbrown/forbidigo/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/ashanbrown/forbidigo/tree/master)

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
forbidigo also supports a more advanced mode where it uses type information to
identify what an expression references. This needs to be enabled through the
`analyze_types` command line parameter. Beware this may have a performance
impact because additional information is required for the analysis.

Replacing the literal source code works for items in a package as in the
`fmt2.Print` example above and also for struct fields and methods. For those,
`<package name>.<type name>.<field or method name>` replaces the source code
text. `<package name>` is what the package declares in its `package` statement,
which may be different from last part of the import path:

      import "example.com/some/pkg" // pkg uses `package somepkg`
      s := somepkg.SomeStruct{}
      s.SomeMethod() // -> somepkg.SomeStruct.SomeMethod

Pointers are treated like the type they point to:

      var cf *spew.ConfigState = ...
      cf.Dump() // -> spew.ConfigState.Dump

When a type is an alias for a type in some other package, the name of that
other package will be used.

An imported identifier gets replaced as if it had been imported without `import .`
*and* also gets matched literally, so in this example both `^ginkgo.FIt$`
and `^FIt$` would catch the usage of `FIt`:

     import . "github.com/onsi/ginkgo/v2"
     FIt(...) // -> ginkgo.FIt, FIt

Beware that looking up the package name has limitations. When a struct embeds
some other type, references to the inherited fields or methods get resolved
with the outer struct as type:

     package foo

     type InnerStruct {
         SomeField int
     }

     func (i innerStruct) SomeMethod() {}

     type OuterStruct {
         InnerStruct
     }

     s := OuterStruct{}
     s.SomeMethod() // -> foo.OuterStruct.SomeMethod
     i := s.SomeField // -> foo.OuterStruct.SomeField

When a method gets called via some interface, that invocation also only
gets resolved to the interface, not the underlying implementation:

    // innerStruct as above

    type myInterface interface {
        SomeMethod()
    }

    var i myInterface = InnerStruct{}
    i.SomeMethod() // -> foo.myInterface.SomeMethod

Using the package name is simple, but the name is not necessarily unique. For
more advanced cases, it is possible to specify more complex patterns. Such
patterns are strings that contain JSON or YAML for a struct.

The full pattern struct has the following fields:

* `msg`: an additional comment that gets added to the error message when a
  pattern matches.
* `p`: the regular expression that matches the source code or expanded
  expression, depending on the global flag.
* `pkg`: a regular expression for the full package import path. The package
  path includes the package version if the package has a version >= 2. This is
  only supported when `analyze_types` is enabled.

To distinguish such patterns from traditional regular expression patterns, the
encoding must start with a `{` or contain line breaks. When using just JSON
encoding, backslashes must get quoted inside strings. When using YAML, this
isn't necessary. The following pattern strings are equivalent:

    {p: "^fmt\\.Println$", msg: "do not write to stdout"}

    {p: ^fmt\.Println$,
     msg: do not write to stdout,
    }

    {p: ^fmt\.Println$, msg: do not write to stdout}

    p: ^fmt\.Println$
    msg: do not write to stdout

A larger set of interesting patterns might include:

-* `^fmt\.Print.*$` -- forbid use of Print statements because they are likely just for debugging
-* `^ginkgo\.F[A-Z].*$` -- forbid ginkgo focused commands (used for debug issues)
-* `^spew\.Dump$` -- forbid dumping detailed data to stdout
-* `^spew.ConfigState\.Dump$` -- also forbid it via a `ConfigState`
-* `^spew\.Dump(# please do not spew to stdout)?$` -- forbid spewing, with a custom message
-* `{p: ^spew\.Dump$, msg: please do not spew to stdout}` -- the same with separate msg field

### Flags
- **-set_exit_status** (default false) - Set exit status to 1 if any issues are found.
- **-exclude_godoc_examples** (default true) - Controls whether godoc examples are identified and excluded
- **-tests** (default true) - Controls whether tests are included
- **-analyze_types** (default false) - Replace literal source code before matching

## Purpose

To prevent leaving format statements and temporary statements such as Ginkgo FIt, FDescribe, etc.

## Ignoring issues

You can ignore a particular issue by including the directive `//permit` on that line.  *This feature is disabled inside `golangci-lint` to encourage ignoring issues using the `// nolint` directive common for all linters (nolinting well is hard and I didn't want to make an effort do it exactly right within this linter).*

## Contributing

Pull requests welcome!
