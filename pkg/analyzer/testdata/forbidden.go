package testdata

import (
	"fmt"
	alias "fmt"

	anotherpkg "example.com/another/pkg"
	somepkg "example.com/some/pkg"
)

func Foo() {
	fmt.Println("here I am") // want "forbidden by pattern"
	fmt.Printf("this is ok") //permit:fmt.Printf // this is ok
	print("not ok")          // want "use of .print. forbidden by pattern .*print.println"
	println("also not ok")   // want "use of .println. forbidden by pattern .*print.println"
	alias.Println("hello")   // want "forbidden by pattern"
	somepkg.Forbidden()      // want "somepkg.Forbidden.*forbidden by pattern .*example.com/some/pkg.*Forbidden"

	c := somepkg.CustomType{}
	c.AlsoForbidden() // want "c.AlsoForbidden.*forbidden by pattern.*example.com/some/pkg.CustomType.*AlsoForbidden"

	// Type alias and pointer.
	c2 := &anotherpkg.CustomType{}
	c2.AlsoForbidden() // want "c2.AlsoForbidden.*forbidden by pattern.*example.com/some/pkg.CustomType.*AlsoForbidden"

	// Interface.
	var ci somepkg.CustomInterface
	ci.StillForbidden() // want "ci.StillForbidden.*forbidden by pattern.*example.com/some/pkg.CustomInterface.*StillForbidden"
}

func Bar() string {
	fmt := struct {
		Println string
	}{}
	return fmt.Println // not the fmt package, not matched
}
