package main

import (
	"fmt"
)

// for k,v := range f2 {
//
// }
func f2(yield func(int, string) bool) {
	for i := 0; i < 10; i++ {
		if !yield(i, fmt.Sprintf("I'm %d ", i)) {
			return
		}
	}
}

// for k := range f1 {
//
// }
func f1(yield func(int) bool) {
	for i := 0; i < 10; i++ {
		if !yield(i) {
			return
		}
	}
}

// for range functionWithReturnsIsZero {
//
// }
func f0(yield func() bool) {
	for i := 0; i < 10; i++ {
		if !yield() {
			return
		}
	}
}

func main() {
	// 1. Basic usage, accepts iterator functions:
	// func(func() bool)
	// func(func(K) bool)
	// func(func(K, V) bool)
	fmt.Println("Test basic usage: for k,v := range f")
	for k, v := range f2 {
		fmt.Println(k, v)
	}

	fmt.Println("Test basic usage: for k := range f ")
	for k := range f1 {
		fmt.Println(k)
	}

	fmt.Println("Test basic usage: for range f ")
	for range f0 {
	}
	{
		for k, v := range fb {
			if k == 7 {
				fmt.Println("go break")
				break
			}
			fmt.Println(v)
		}
	}

	s := []string{"hello", "world"}
	for i, x := range MyBackward(s) {
		fmt.Println(i, x)
	}
}

func MyBackward[E any](s []E) func(func(int, E) bool) {
	return func(yield func(int, E) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(i, s[i]) {
				return
			}
		}
	}
}

func fb(yield func(int, string) bool) {
	for i := 0; i < 10; i++ {
		if !yield(i, fmt.Sprintf("I'm %d ", i)) {
			fmt.Println("ok break")
			return
		}
	}
}
