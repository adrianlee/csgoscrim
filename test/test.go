package main

import (
	"fmt"
)

func main() {
	// print
	var i int = 0
	fmt.Println("hello", i)

	// for
	j := 0
	for j < 3 {
		// if
		if j == 3 {
			break
		} else {
			fmt.Println(j)
		}
		j++
	}

	// switch/function
	fmt.Println(switchTest("asdl;fjkasd"))
	fmt.Println(switchTest("omg"))

	// multi return
	a, o := multiReturn()
	fmt.Println(a, o)

	// interaces

}

func switchTest(k string) string {
	switch k {
	case "omg":
		return "omg"
	default:
		return "default"
	}
}

func multiReturn() (string, int) {
	return "hello", 123
}

// interfaces
type geometry interface {
	area() float64
	perim() float64
}
