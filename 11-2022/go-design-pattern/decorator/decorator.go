package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

func decorator(f func(s string)) func(s string) {

	return func(s string) {
		fmt.Println("Started")
		f(s)
		fmt.Println("Done")
	}
}

func Hello(s string) {
	fmt.Println(s)
}

func main() {
	decorator(Hello)("Hello, World!")
}

//我们可以看到，我们动用了一个高阶函数 decorator()，在调用的时候，先把 Hello() 函数传进去，然后其返回一个匿名函数，这个匿名函数中除了运行了自己的代码，也调用了被传入的 Hello() 函数。

type SumFunc func(int64, int64) int64

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func timedSumFunc(f SumFunc) SumFunc {
	return func(start, end int64) int64 {

		defer func(t time.Time) {
			fmt.Printf("--- Time Elapsed (%s): %v ---\n",
				getFunctionName(f), time.Since(t))
		}(time.Now())

		return f(start, end)
	}
}

func Sum1(start, end int64) int64 {
	var sum int64
	sum = 0
	if start > end {
		start, end = end, start
	}
	for i := start; i <= end; i++ {
		sum += i
	}
	return sum
}

func Sum2(start, end int64) int64 {
	if start > end {
		start, end = end, start
	}
	return (end - start + 1) * (end + start) / 2
}

func main() {

	sum1 := timedSumFunc(Sum1)
	sum2 := timedSumFunc(Sum2)

	fmt.Printf("%d, %d\n", sum1(-10000, 10000000), sum2(-10000, 10000000))
}

//关于上面的代码，有几个事说明一下：

//1）有两个 Sum 函数，Sum1() 函数就是简单的做个循环，Sum2() 函数动用了数据公式。（注意：start 和 end 有可能有负数的情况）

//2）代码中使用了 Go 语言的反射机器来获取函数名。

//3）修饰器函数是 timedSumFunc()

//运行后输出：

//$ go run time.sum.go
//--- Time Elapsed (main.Sum1): 3.557469ms ---
//--- Time Elapsed (main.Sum2): 291ns ---
//49999954995000, 49999954995000
