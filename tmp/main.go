package main

import (
	"fmt"
)

func test() (x int) {
	defer func() {
		fmt.Println("in defer", &x)
		x++
	}()
	x = 1
	return
}

func anotherTest() int {
	var x int
	defer func() {
		fmt.Println("in defer", &x)
		x++
	}()
	x = 1
	return x
}

func main() {
	y := anotherTest()
	fmt.Println(y, &y)
}
