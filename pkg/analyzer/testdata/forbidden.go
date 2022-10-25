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
	print("not ok")          // want "forbidden by pattern"
	println("also not ok")   // want "forbidden by pattern"
	alias.Println("hello")   // not matched by default pattern fmt.Println
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
	return fmt.Println // want "forbidden by pattern"
}
