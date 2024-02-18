package main

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

//我们的种子语料库里每个符号都是单个字节。但是像 泃这样的中文符号由多个字节组成，如果以字节为维度进行反转，就会得到无效的结果。
func ReverseBug(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

//上面的例子里，输入的字符串是只有1个byte的字节切片，这1个byte是\x91。

//当我们把这个输入的字符串转成[]rune时，Go会把字节切片编码为UTF-8，于是就把\x91替换成了'�'，'�'饭庄后还是'�'，一次就导致原字符串\x91反转后得到的字符串是'�'了。

//现在问题明确了，是因为输入的数据是非法的unicode。那接下来我们就可以修正Reverse函数的实现了。
func Reverse(s string) (string, error) {
	if !utf8.ValidString(s) {
		return s, errors.New("input is not valid UTF-8")
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}

func main() {
	input := "The quick brown fox jumped over the lazy dog"
	rev := Reverse(input)
	doubleRev := Reverse(rev)
	fmt.Printf("original: %q\n", input)
	fmt.Printf("reversed: %q\n", rev)
	fmt.Printf("reversed again: %q\n", doubleRev)
}
