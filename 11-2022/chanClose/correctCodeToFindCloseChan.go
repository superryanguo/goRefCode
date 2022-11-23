package main

import "fmt"

func main() {
	ch := make(chan int, 1)
	fmt.Printf("chan close is %v\n", isChanCloseCorrect(ch))
	close(ch)
	fmt.Printf("chan close is %v\n", isChanCloseCorrect(ch))
	ch2 := make(chan int, 1)
	fmt.Printf("chan close is %v\n", isChanCloseWrong(ch2))
	close(ch2)
	fmt.Printf("chan close is %v\n", isChanCloseWrong(ch2))
}
func isChanCloseCorrect(ch chan int) bool {
	select {
	case _, received := <-ch:
		return !received
	default:
	}
	return false
}
func isChanCloseWrong(ch chan int) bool {
	_, ok := <-ch
	return ok
}
