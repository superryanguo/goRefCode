TODO:string 通常指向字符串字面量，字面量存储的位置是只读段，并不是堆或栈上，所以 string 不能被修改。
package main

string 概念

源代码中 src/builtin/builtin.go string 的描述如下：


 // string is the set of all strings of 8-bit bytes, conventionally but not
// necessarily representing UTF-8-encoded text. A string may be empty, but
// not nil. Values of string type are immutable.
type string string


        string 是所有 8 位字节字符串的集合，通常但不一定代表 UTF-8 编码的文本。

        字符串可以为空（长度为 0），但不会是 nil。

        字符串类型的值是不可变的。

string 数据结构

源码包中 src/runtime/string.go:stringStruct 定义的 string 的数据结构如下：


type stringStruct struct {
  str unsafe.Pointer
  len int
}


    str : 字符串的首地址

    len : 字符串的长度


发现 string 的数据结构有点类似于切片，切片比它多了一个容量成员。string 和 byte 切片经常互转。
string 操作

字符串的构建是先构建 stringStruct，再转换成 string：


//go:nosplit
func gostringnocopy(str *byte) string {
  ss := stringStruct{str: unsafe.Pointer(str), len: findnull(str)} // 先构造 stringStruct
  s := *(*string)(unsafe.Pointer(&ss)) // stringStruct 转换成 string
  return s
}

[]byte 转 string

示例：


func ByteToString(s []byte) string {
  return string(s)
}


注意这里的转换进行一次内存拷贝：


    根据切片长度申请内存空间（假设内存地址为 p，切片长度为 len(b)）

    构建 string（string.str = p; string.len = len）

    拷贝数据（切片中的数据拷贝到新的内存空间）

string 转 []byte

示例：


func StringToByte(str string) []byte {
  return []byte(str)
}


注意这里的转换也会进行一次内存拷贝：


    申请切片内存空间

    将 string 拷贝到切片

字符串拼接

示例：


str := "str1" + "str2" + "str3" + "str4"


    新字符串的内存空间是一次性分配好的，所以即使有很多的字符串进行拼接，性能也会有很好的保证。

    拼接语句在编译时会先放到一个切片中，然后再两次遍历此切片，一次获取字符串长度用来申请内存，一次用来把字符串逐个拷贝过去。

    由于 string 是不能修改的，源码在拼接过程中会用 rawstring() 方法生成一个指定大小的 string，并同时返回一个切片，二者共享同一块内存空间，后面向切片中拷贝数据，也就间接修改了 string。


rawstring()源代码如下：


// rawstring allocates storage for a new string. The returned
// string and byte slice both refer to the same storage.
// The storage is not zeroed. Callers should use
// b to set the string contents and then drop b.
func rawstring(size int) (s string, b []byte) {
  p := mallocgc(uintptr(size), nil, false)

  stringStructOf(&s).str = p
  stringStructOf(&s).len = size

  *(*slice)(unsafe.Pointer(&b)) = slice{p, size, size}

  return
}

string 不能修改

在 Go 中，string 不包含内存空间，它只有一个内存指针，所以 string 非常轻量，很方便进行传递且不用担心内存拷贝。
TODO:string 通常指向字符串字面量，字面量存储的位置是只读段，并不是堆或栈上，所以 string 不能被修改。
[]byte 转 string 不拷贝内存的情况

有时只是临时需要字符串的场景下，byte 切片转换成 string 时并不会拷贝内存，而是直接返回一个 string，这个 string 的指针指向切片的内存。


如：


    使用 m[string(b)] 来查找 map（map 是 string 为 key，临时把切片 b 转成 string）

    字符串拼接，如”<” + “string(b)” + “>”

    字符串比较：string(b) == “foo”

用 []byte 还是 string

[]byte 和 string 都可以表示字符串，它们数据结构不同，其衍生出来的方法也不同。


string 擅长的场景：


    需要字符串比较；

    不需要 nil 字符串；


[]byte 擅长的场景：


    修改字符串的时候；

    函数返回值，需要使用 nil 来表示含义；

    需要切片操作；
