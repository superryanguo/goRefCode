package main

func main() {
	println(multiply(2, 3))
	println(multiply(100, 2)) // overflows as a negative number
	println(multiply(100, 3)) // overflows as a positive number
	println(multiply(64, 2))  // 128
	println(multiply(127, 1)) // 128
}

const detectOverflows = true

//int8's max value is 127
func multiply(a, b int8) int8 {
	result := a * b
	if detectOverflows && result/a != b {
		println("overflowed!", a, b)
		return -1
	}
	return result
}
