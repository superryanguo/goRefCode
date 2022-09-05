package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func is_prime(value int) bool {
	for i := 2; i <= int(math.Floor(float64(value)/2)); i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func merge(cs []<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	wg.Add(len(cs))
	for k, c := range cs {
		go func(c <-chan int, k int) {
			for n := range c {
				fmt.Printf("From chan %d get the value %d\n", k, n)
				out <- n
			}
			wg.Done()
		}(c, k) //ryan: you should pass the k here, or k will always as 9
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
func echo(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		var sum = 0
		for n := range in {
			sum += n
		}
		out <- sum
		close(out)
	}()
	return out
}

//ryan: generate a gorountine to check if it's prime value
func prime(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if is_prime(n) {
				out <- n
				fmt.Println("number get in=", n)
			}
		}
		close(out)
	}()
	return out
}

func main() {
	nums := makeRange(1, 10000)
	in := echo(nums)

	t0 := time.Now()
	const nProcess = 10
	var chans [nProcess]<-chan int
	for i := range chans {
		chans[i] = sum(prime(in)) //ryan: interesting that in this range, in is a channel
		//And we will automaticlly split the message from "in" in a cocurrent way
	}
	for n := range sum(merge(chans[:])) {
		fmt.Println(n)
	}
	dt := time.Since(t0)

	fmt.Printf("the duration is=%v\n", dt)
}
