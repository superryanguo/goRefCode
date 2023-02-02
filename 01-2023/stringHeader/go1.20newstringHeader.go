package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

//SliceHeader 是 Slice（切片）的运行时表现；StringHeader 是 String（字符串）的运行时表现。

//背景
//为什么这两个运行时结构体受到了那么多的关注呢？是因为常被用于如下场景：

//将 []byte 转换为 string。
//将 string 转换为 []byte。
//抓取数据指针（data pointer）字段用于 ffi 或其他用途。
//将一种类型的 slice 转换为另一种类型的 slice。
//常见案例，可见如下代码：

//s := "脑子进煎鱼了？重背面试题(doge"
//h := (*reflect.StringHeader)(unsafe.Pointer(&s))
//又或是自己构造一个：

//unsafe.Pointer(&reflect.StringHeader{
//Data: uintptr(unsafe.Pointer(&s.Data[0])),
//Len:  int(s.Size),
//})
//似乎看起来没什么问题，所以在业内打开了一种新的姿势。那就是借助 (String|Slice) Header 来实现零拷贝的 string 到 bytes 的转换，得到了广大开发者的使用。毕竟谁都想性能高一点。

//如下转换代码：

func main() {
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

//当然，还有更多基于 Header 自己写的转换，甚至写错写到泄露没法被 GC 的，又或是抛出 throw 致命错误查了几周的。

//问题
//今年 Go 团队进行了讨论，通过分析、搜索发现 reflect.SliceHeader 和 reflect.StringHeader 在业内经常被滥用，且使用不方便，很容易出错（要命的是很隐性的那种

//这个坑也在于，SliceHeader 和 StringHeader 的 Data 字段（后称数据指针）：

type StringHeader struct {
	Data uintptr
	Len  int
}

//类型是 uintptr不是 unsafe.Pointer。设什么都可以，灵活度过于高，非常容易搞出问题。

//TODO: Go1.20 新特性
//在 Go1.20 起，在 unsafe 标准库新增了 3 个函数来替代前面这两个类型的使用。希望能够进一步标准化，并提供额外的类型安全。
//如下函数签名：

//func String(ptr *byte, len IntegerType) string：根据数据指针和字符长度构造一个新的 string。
//func StringData(str string) *byte：返回指向该 string 的字节数组的数据指针。
//func SliceData(slice []ArbitraryType) *ArbitraryType：返回该 slice 的数据指针。
//新版本的用法将会变成：

func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BytesToString(b []byte) string {
	return unsafe.String(&b[0], len(b))
}

//以往常用的 reflect.SliceHeader 和 reflect.StringHeader 将会被标注为被废弃。
