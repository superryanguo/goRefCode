package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	// 对字符串加密，方式一
	h := sha256.New()
	h.Write([]byte("sha256 加密测试！"))
	fmt.Printf("%x\n", h.Sum(nil))
	fmt.Printf("%X\n", h.Sum(nil)) // 大写的 X，代表大写的十六进制字符串

	// 对字符串加密，方式二
	hh := sha256.New()
	hh.Write([]byte("sha256 加密测试！"))
	fmt.Print(hex.EncodeToString(hh.Sum(nil)) + "\n")
	fmt.Print(strings.ToTitle(hex.EncodeToString(hh.Sum(nil))) + "\n")

	// 对文件进行加密
	f, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h2 := sha256.New()
	if _, err := io.Copy(h2, f); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%x\n", h2.Sum(nil))
	fmt.Printf("%X", h2.Sum(nil))
}
