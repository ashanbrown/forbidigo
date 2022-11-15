package testdata

import (
	"fmt"
	alias "fmt"

	"example.com/some/pkg"
)

func Foo() {
	fmt.Println("here I am") // want "forbidden by pattern"
	fmt.Printf("this is ok") //permit:fmt.Printf // this is ok
	print("not ok")          // no package, not matched
	println("also not ok")   // no package, not matched
	alias.Println("hello")   // want "forbidden by pattern"
	pkg.Forbidden()          // want "pkg.Forbidden.*forbidden by pattern .*example.com/some/pkg.*Forbidden"
}

func Bar() string {
	fmt := struct {
		Println string
	}{}
	return fmt.Println // not the fmt package, not matched
}
