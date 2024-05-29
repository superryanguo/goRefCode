package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

func main() {
	h := md5.New()
	io.WriteString(h, "md5 加密测试！")
	fmt.Printf("%x\n", h.Sum(nil))
	fmt.Printf("%X\n", h.Sum(nil)) // 大写的 X，代表大写的十六进制字符串

	hh := md5.New()
	hh.Write([]byte("md5 加密测试！"))
	fmt.Print(hex.EncodeToString(hh.Sum(nil)) + "\n")
	fmt.Print(strings.ToTitle(hex.EncodeToString(hh.Sum(nil)))) // strings.ToTitle() 转大写
}
