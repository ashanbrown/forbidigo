// Package renamed intentionally uses a package name that does not match the
// import path to test this situation.
package renamed // import "example.com/some/renamedpkg"

func ForbiddenFunc() {
}

type Struct struct{}

func (s Struct) ForbiddenMethod() {
}
