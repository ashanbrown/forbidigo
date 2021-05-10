package examples

import "fmt"

func Foo() {
	fmt.Println("here I am")
	fmt.Printf("this is ok") //permit:fmt.Printf // this is ok
	print("not ok")
	println("also not ok")
}
