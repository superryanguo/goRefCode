package main

import "runtime"

func main() {
	runtime.Stack ///usr/lib/go-1.18/src/runtime/mprof.go

	///usr/lib/go-1.18/src/runtime/chan.go
	v := <-c
	v, ok := <-c
	//对应函数分别是 chanrecv1 和 chanrecv2 ，位于 runtime/chan.go 文件

	c <- x
	//对应函数实现 chansend ，位于 runtime/chan.go 文件。
}
