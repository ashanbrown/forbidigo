package testdata

import (
	"fmt"
	alias "fmt"

	anotherpkg "example.com/another/pkg"
	somepkg "example.com/some/pkg"
	renamed "example.com/some/renamedpkg" // Package name is "renamed".
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

var forbiddenFunctionRef = somepkg.Forbidden

var forbiddenVariableRef = Shiny      // want "Shiny.*forbidden by pattern.*\\^Shiny"
var forbiddenVariableRef2 = AlsoShiny // want "Shiny.*forbidden by pattern.*\\^AlsoShiny"

func Foo() {
	fmt.Println("here I am") // want "forbidden by pattern"
	fmt.Printf("this is ok") //permit:fmt.Printf // this is ok
	print("not ok")          // want "forbidden by pattern"
	println("also not ok")   // want "forbidden by pattern"
	alias.Println("hello")   // not matched by default pattern fmt.Println
	somepkg.Forbidden()

	c := somepkg.CustomType{}
	c.AlsoForbidden()

	// Selector expression with result of function call in package.
	somepkg.NewCustom().AlsoForbidden()

	// Selector expression with result of normal function call.
	myCustom().AlsoForbidden()

	// Selector expression with result of normal function call.
	myCustomFunc().AlsoForbidden()

	// Type alias and pointer.
	c2 := &anotherpkg.CustomTypeAlias{}
	c2.AlsoForbidden()

	// Interface.
	var ci somepkg.CustomInterface
	ci.StillForbidden()

	// Struct embedded inside another: must be forbidden separately!
	myCustomStruct{}.AlsoForbidden()
	_ = myCustomStruct{}.ForbiddenField

	// Forbidden method called via interface: must be forbidden separately!
	var ci2 myCustomInterface = somepkg.CustomType{}
	ci2.AlsoForbidden()

	// Package name != import path.
	renamed.ForbiddenFunc() // want "renamed.Forbidden.* by pattern .*renamed..Forbidden"
	renamed.Struct{}.ForbiddenMethod()
}

func Bar() string {
	fmt := struct {
		Println string
	}{}
	return fmt.Println // want "forbidden by pattern"
}
