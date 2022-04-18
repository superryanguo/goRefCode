package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Test the concurrentAppendSlice")
	concurrentAppendSliceForceIndex()
	concurrentAppendSliceNotForceIndex()
}

func concurrentAppendSliceNotForceIndex() {
	sl := make([]int, 0)
	wg := sync.WaitGroup{}
	for index := 0; index < 100; index++ {
		k := index
		wg.Add(1)
		go func(num int) {
			sl = append(sl, num)
			wg.Done()
		}(k)
	}
	wg.Wait()
	fmt.Printf("final unforce len(sl)=%d cap(sl)=%d\n", len(sl), cap(sl))
}
func concurrentAppendSliceForceIndex() {
	sl := make([]int, 100)
	wg := sync.WaitGroup{}
	for index := 0; index < 100; index++ {
		k := index
		wg.Add(1)
		go func(num int) {
			sl[num] = num
			wg.Done()
		}(k)
	}
	wg.Wait()
	fmt.Printf("final force len(sl)=%d cap(sl)=%d\n", len(sl), cap(sl))
}

//slice支持并发吗？
//我们都知道切片是对数组的抽象，其底层就是数组，在并发下写数据到相同的索引位会被覆盖，并且切片也有自动扩容的功能，当切片要进行扩容时，就要替换底层的数组，在切换底层数组时，多个goroutine是同时运行的，哪个goroutine先运行是不确定的，不论哪个goroutine先写入内存，肯定就有一次写入会覆盖之前的写入，所以在动态扩容时并发写入数组是不安全的；

//所以当别人问你slice支持并发时，你就可以这样回答它：

//当指定索引使用切片时，切片是支持并发读写索引区的数据的，但是索引区的数据在并发时会被覆盖的；当不指定索引切片时，并且切片动态扩容时，并发场景下扩容会被覆盖，所以切片是不支持并发的～。

//github上著名的iris框架也曾遇到过切片动态扩容导致webscoket连接数减少的bug，最终采用sync.map解决了该问题，感兴趣的可以看一下这个issue:https://github.com/kataras/iris/pull/1023#event-1777396646；

//总结
//针对上述问题，我们可以多种方法来解决切片并发安全的问题：

//加互斥锁
//使用channel串行化操作
//使用sync.map代替切片
