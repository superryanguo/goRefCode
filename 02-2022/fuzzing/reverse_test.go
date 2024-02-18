package main

import (
	"testing"
	"unicode/utf8"
)

func TestReverse(t *testing.T) {
	testcases := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{" ", " "},
		{"!12345", "54321!"},
	}
	for _, tc := range testcases {
		rev := Reverse(tc.in)
		if rev != tc.want {
			t.Errorf("Reverse: %q, want %q", rev, tc.want)
		}
	}
}

//单元测试有局限性，每个测试输入必须由开发者指定加到单元测试的测试用例里。

//fuzzing的优点之一是可以基于开发者代码里指定的测试输入作为基础数据，进一步自动生成新的随机测试数据，用来发现指定测试输入没有覆盖到的边界情况。

//在这，我们会把单元测试转换成模糊测试，这样可以更轻松地生成更多的测试输入。
//注意：fuzzing模糊测试和Go已有的单元测试以及性能测试框架是互为补充的，并不是替代关系。

//比如我们实现的Reverse函数如果是一个错误的版本，直接return返回输入的字符串，是完全可以通过上面的模糊测试的，但是没法通过我们前面编写的单元测试。因此单元测试和模糊测试是互为补充的，不是替代关系。

//Go模糊测试和单元测试在语法上有如下差异：

//Go模糊测试函数以FuzzXxx开头，单元测试函数以TestXxx开头
//Go模糊测试函数以 *testing.F作为入参，单元测试函数以*testing.T作为入参
//Go模糊测试会调用f.Add函数和f.Fuzz函数。

//f.Add函数把指定输入作为模糊测试的种子语料库(seed corpus)，fuzzing基于种子语料库生成随机输入。

//f.Fuzz函数接收一个fuzz target函数作为入参。fuzz target函数有多个参数，第一个参数是*testing.T，其它参数是被模糊的类型(注意：被模糊的类型目前只支持部分内置类型, 列在 Go Fuzzing docs，未来会支持更多的内置类型)。
//f.Add(5,"hello")
//f.Fuzz(func(t t*testing.T, i int, s string))

func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rev := Reverse(orig)
		doubleRev := Reverse(rev)
		t.Logf("Number of runes: orig=%d, rev=%d, doubleRev=%d", utf8.RuneCountInString(orig), utf8.RuneCountInString(rev), utf8.RuneCountInString(doubleRev))
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}

//执行如下命令来运行模糊测试。
//这个方式只会使用种子语料库，而不会生成随机测试数据。通过这种方式可以用来验证种子语料库的测试数据是否可以测试通过。(fuzz test without fuzzing)
//$ go test
//PASS
//ok      example/fuzz  0.013s
//如果reverse_test.go文件里有其它单元测试函数或者模糊测试函数，但是只想运行FuzzReverse模糊测试函数，我们可以执行go test -run=FuzzReverse命令。

//注意：go test默认会执行所有以TestXxx开头的单元测试函数和以FuzzXxx开头的模糊测试函数，默认不运行以BenchmarkXxx开头的性能测试函数，如果我们想运行 benchmark用例，则需要加上 -bench 参数。

//如果要基于种子语料库生成随机测试数据用于模糊测试，需要给go test命令增加 -fuzz参数。(fuzz test with fuzzing)
//$ go test -fuzz=Fuzz
//上面的fuzzing测试结果是FAIL，引起FAIL的输入数据被写到了一个语料库文件里。下次运行go test命令的时候，即使没有-fuzz参数，这个语料库文件里的测试数据也会被用到。
//我们实现的Reverse函数是按照字节(byte)为维度进行字符串反转，这就是问题所在。

//比如中文里的字符 泃其实是由3个字节组成的，如果按照字节反转，反转后得到的就是一个无效的字符串了。

//因此为了保证字符串反转后得到的仍然是一个有效的UTF-8编码的字符串，我们要按照rune进行字符串反转。
func FuzzReverseFull(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rev, err1 := Reverse(orig)
		if err1 != nil {
			return
		}
		doubleRev, err2 := Reverse(rev)
		if err2 != nil {
			return
		}
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}

//fuzz test如果没有遇到错误，默认会一直运行下去，需要使用 ctrl-C 结束测试。

//也可以传递-fuzztime参数来指定测试时间，这样就不用 ctrl-C 了。

//指定测试时间。 go test -fuzz=Fuzz -fuzztime 30s 如果没有遇到错误会执行30s后自动结束。
