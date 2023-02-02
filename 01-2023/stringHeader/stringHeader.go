package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

//type SliceHeader struct {
//Data uintptr
//Len  int
//Cap  int
//}
//Data：指向具体的底层数组。
//Len：代表切片的长度。
//Cap：代表切片的容量。
func main() {

	// 初始化底层数组
	ss := [4]string{"脑子", "进", "煎鱼", "了"}
	ss1 := ss[0:1]
	ss2 := ss[:]

	// 构造 SliceHeader
	sh1 := (*reflect.SliceHeader)(unsafe.Pointer(&ss1))
	sh2 := (*reflect.SliceHeader)(unsafe.Pointer(&ss2))
	fmt.Println(sh1.Len, sh1.Cap, sh1.Data)
	fmt.Println(sh2.Len, sh2.Cap, sh2.Data)
	//两个切片的 Data 属性所指向的底层数组是一致的，Len 属性的值不一样，sh1 和 sh2 分别是两个切片。
	//疑问
	//为什么两个新切片所指向的 Data 是同一个地址的呢？

	//这其实是 Go 语言本身为了减少内存占用，提高整体的性能才这么设计的。

	//将切片复制到任意函数的时候，对底层数组大小都不会影响。复制时只会复制切片本身（值传递），不会涉及底层数组。

	//也就是在函数间传递切片，其只拷贝 24 个字节（指针字段 8 个字节，长度和容量分别需要 8 个字节），效率很高。

	//坑
	//这种设计也引出了新的问题，在平时通过 s[i:j] 所生成的新切片，两个切片底层指向的是同一个底层数组。

	s := "脑子进煎鱼了"
	s1 := "脑子进煎鱼了"
	s2 := "脑子进煎鱼了"[7:]

	fmt.Printf("%d \n", (*reflect.StringHeader)(unsafe.Pointer(&s)).Data)
	fmt.Printf("%d \n", (*reflect.StringHeader)(unsafe.Pointer(&s1)).Data)
	fmt.Printf("%d \n", (*reflect.StringHeader)(unsafe.Pointer(&s2)).Data)
	//从输出结果来看，变量 s 和 s1 指向同一个内存地址。变量 s2 虽稍有偏差，但本质上也是指向同一块。

	//因为其是字符串的切片操作，是从第 7 位索引开始，因此正好的 17608234-17608227 = 7。也就是三个变量都是指向同一块内存空间，这是为什么呢？

	//这是因为在 Go 语言中，字符串都是只读的，为了节省内存，相同字面量的字符串通常对应于同一字符串常量，因此指向同一个底层数组。
}
