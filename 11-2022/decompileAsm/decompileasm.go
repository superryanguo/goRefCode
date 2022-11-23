package main

import "fmt"

func main() {
	fmt.Println("vim-go")
}

//go汇编代码
//go tool compile -S x.go
//go build -gcflags -S x.go

//查看反汇编
//go tool objdump
//这是Go语言自带的反编译命令。

//$ go build x.go
//$ go tool objdump -s main.main x

//objdump
//这是Linux环境中一个通用的反编译工具，不仅仅适用于Go程序。

//$ objdump --disassemble=main.main x

//https://lrita.github.io/2017/12/12/golang-asm/
