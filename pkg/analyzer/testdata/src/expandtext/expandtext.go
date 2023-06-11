package testdata

import (
	"fmt"
	alias "fmt"

	anotherpkg "example.com/another/pkg"
	somepkg "example.com/some/pkg"
	"example.com/some/renamedpkg" // Package name is "renamed".
	. "example.com/some/thing"
)

func myCustom() somepkg.CustomType {
	return somepkg.CustomType{}
}

var myCustomFunc = myCustom

type myCustomStruct struct {
	somepkg.CustomType
}

type myCustomInterface interface {
	AlsoForbidden() *somepkg.Result
}

var forbiddenFunctionRef = somepkg.Forbidden // want "somepkg.Forbidden.*forbidden by pattern .*\\^pkg.*Forbidden"

var forbiddenVariableRef = Shiny      // want "Shiny.*forbidden by pattern.*\\^Shiny"
var forbiddenVariableRef2 = AlsoShiny // want "Shiny.*forbidden by pattern.*\\^AlsoShiny"

func Foo() {
	fmt.Println("here I am") // want "forbidden by pattern"
	fmt.Printf("this is ok") //permit:fmt.Printf // this is ok
	print("not ok")          // want "forbidden by pattern"
	println("also not ok")   // want "forbidden by pattern"
	alias.Println("hello")   // want "forbidden by pattern"
	somepkg.Forbidden()      // want "somepkg.Forbidden.*forbidden by pattern .*\\^pkg.*Forbidden"

	c := somepkg.CustomType{}
	c.AlsoForbidden() // want "c.AlsoForbidden.*forbidden by pattern.*\\^pkg..CustomType.*Forbidden"
	_ = c.AlsoForbidden().Value // want "c.AlsoForbidden.*forbidden by pattern.*\\^pkg..CustomType.*Forbidden"

	// Selector expression with result of function call in package.
	somepkg.NewCustom().AlsoForbidden() // want "somepkg.NewCustom...AlsoForbidden.*forbidden by pattern.*\\^pkg..CustomType.*Forbidden"

	// Selector expression with result of normal function call.
	myCustom().AlsoForbidden() // want "myCustom...AlsoForbidden.*forbidden by pattern.*\\^pkg..CustomType.*Forbidden"

	// Selector expression with result of normal function call.
	myCustomFunc().AlsoForbidden() // want "myCustomFunc...AlsoForbidden.*forbidden by pattern.*\\^pkg..CustomType.*Forbidden"

	// Type alias and pointer.
	c2 := &anotherpkg.CustomTypeAlias{}
	c2.AlsoForbidden() // want "c2.AlsoForbidden.*forbidden by pattern.*\\^pkg..CustomType.*Forbidden"

	// Interface.
	var ci somepkg.CustomInterface
	ci.StillForbidden() // want "ci.StillForbidden.*forbidden by pattern.*\\^pkg..CustomInterface..Forbidden"

	// Forbidden embedded inside another: must be forbidden separately!
	myCustomStruct{}.AlsoForbidden()    // want "myCustomStruct...AlsoForbidden.*forbidden by pattern.*myCustomStruct"
	_ = myCustomStruct{}.ForbiddenField // want "myCustomStruct...ForbiddenField.*forbidden by pattern.*myCustomStruct"

	// Forbidden method called via interface: must be forbidden separately!
	var ci2 myCustomInterface = somepkg.CustomType{}
	ci2.AlsoForbidden() // want "ci2.AlsoForbidden.*forbidden by pattern.*myCustomInterface"

	// Package name != import path.
	renamed.ForbiddenFunc()            // want "renamed.Forbidden.* by pattern .*renamed..Forbidden"
	renamed.Struct{}.ForbiddenMethod() // want "renamed.Struct...ForbiddenMethod.* by pattern .*renamed.*Struct.*Forbidden"
}

func Bar() string {
	fmt := struct {
		Println string
	}{}
	return fmt.Println // not matched because it expands to `struct{Println string}.Println`
}
