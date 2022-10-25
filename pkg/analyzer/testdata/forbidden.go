package testdata

import (
	"fmt"
	alias "fmt"

	"example.com/some/pkg"
)

func Foo() {
	fmt.Println("here I am") // want "forbidden by pattern"
	fmt.Printf("this is ok") //permit:fmt.Printf // this is ok
	print("not ok")          // want "forbidden by pattern"
	println("also not ok")   // want "forbidden by pattern"
	alias.Println("hello")   // not matched by default pattern fmt.Println
	pkg.Forbidden()          // want "pkg.Forbidden.*forbidden by pattern .*example.com/some/pkg.*Forbidden"
}

func Bar() string {
	fmt := struct {
		Println string
	}{}
	return fmt.Println // want "forbidden by pattern"
}
