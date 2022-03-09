package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func Clone(s string) string {
	if len(s) == 0 {
		return ""
	}
	b := make([]byte, len(s))
	copy(b, s)
	return *(*string)(unsafe.Pointer(&b))
}

// Clone returns a copy of b
func Clone(b []byte) []byte {
	b2 := make([]byte, len(b))
	copy(b2, b)
	return b2
}
func equal() {
	b2 := append([]byte(nil), b...)
}

//++++++++++++++++++
type StringHeader struct {
	Data uintptr
	Len  int
}

func main() {
	s0 := "asfbsfsdfdsfsdf"
	s1 := s0[:3]
	s0h := (*reflect.StringHeader)(unsafe.Pointer(&s0))
	s1h := (*reflect.StringHeader)(unsafe.Pointer(&s1))

	fmt.Printf("Len is equal: %t\n", s0h.Len == s1h.Len)
	fmt.Printf("Data is equa: %t\n", s0h.Data == s1h.Data)
}

//supposed to be:
//Len is equal: false
//Data is equa: true
