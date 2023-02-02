package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

//拷贝转换
//为什么会有人关注到 SliceHeader、StringHeader 这类运行时细节呢，一大部分原因是业内会有开发者，希望利用其实现零拷贝的 string 到 bytes 的转换。

//常见转换代码如下：

func string2bytes(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))

	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}

	return *(*[]byte)(unsafe.Pointer(&bh))
}

//TODO: 但这其实是错误的，官方明确表示：

//the Data field is not sufficient to guarantee the data it references will not be garbage collected, so programs must keep a separate, correctly typed pointer to the underlying data.

//SliceHeader、StringHeader 的 Data 字段是一个 uintptr 类型。由于 Go 语言只有值传递。

//因此在上述代码中会出现将 Data 作为值拷贝的情况，这就会导致无法保证它所引用的数据不会被垃圾回收（GC）。

//应该使用如下转换方式：

func main1() {
	s := "脑子进煎鱼了"
	v := string2bytes1(s)
	fmt.Println(v)
}

func string2bytes1(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))

	var b []byte
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pbytes.Data = stringHeader.Data
	pbytes.Len = stringHeader.Len
	pbytes.Cap = stringHeader.Len

	return b
}

//Ryan:和上一个版本相比，确实要显示的声明一个b来指向这个data，否则类似第一种只是返回了一个数据，没有声明一个对象
//在程序必须保留一个单独的、正确类型的指向底层数据的指针。

//在性能方面，若只是期望单纯的转换，对容量（cap）等字段值不敏感，也可以使用以下方式：

func string2bytes2(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

//性能对比：

//string2bytes1-1000-4   3.746 ns/op  0 allocs/op
//string2bytes1-1000-4   3.713 ns/op  0 allocs/op
//string2bytes1-1000-4   3.969 ns/op  0 allocs/op

//string2bytes2-1000-4   2.445 ns/op  0 allocs/op
//string2bytes2-1000-4   2.451 ns/op  0 allocs/op
//string2bytes2-1000-4   2.455 ns/op  0 allocs/op
//会相当标准的转换性能会稍快一些，这种强转也会导致一个小问题。

//代码如下：

func main() {
	s := "脑子进煎鱼了"
	v := string2bytes2(s)
	println(len(v), cap(v))
}

//输出结果：

//18 824633927632
//这种强转其会导致 byte 的切片容量非常大，需要特别注意。一般还是推荐使用标准的 SliceHeader、StringHeader 方式就好了，也便于后来的维护者理解。
