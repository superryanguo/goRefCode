package main

import "fmt"

func foo(x, y int) (z int) {
	fmt.Printf("x=%d, y=%d, z=%d\n", x, y, z)
	z = x + y

	return
}

func main() {
	x := 99
	y := x * x
	z := foo(x, y)

	fmt.Printf("z=%d\n", z)
}
