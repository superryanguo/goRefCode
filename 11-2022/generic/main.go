//For go version>go1.18
//key:  the [T] should be use as the generice type!!!
//key: any, comparable
//$ go run -gcflags=-G=3 ./main.go
package main

import "fmt"

func print[T any] (arr []T) {
  for _, v := range arr {
    fmt.Print(v)
    fmt.Print(" ")
  }
  fmt.Println("")
}

func main() {
  strs := []string{"Hello", "World",  "Generics"}
  decs := []float64{3.14, 1.14, 1.618, 2.718 }
  nums := []int{2,4,6,8}

  print(strs)
  print(decs)
  print(nums)
}

func find[T comparable] (arr []T, elem T) int {
  for i, v := range arr {
    if  v == elem {
      return i
    }
  }
  return -1
}
//从上面的这两个小程序来看，Go语言的泛型已基本可用了，只不过，还有三个问题：

//一个是 fmt.Printf()中的泛型类型是 %v 还不够好，不能像c++ iostream重载 >> 来获得程序自定义的输出。
//另外一个是，go不支持操作符重载，所以，你也很难在泛型算法中使用“泛型操作符”如：== 等
//最后一个是，上面的 find() 算法依赖于“数组”，对于hash-table、tree、graph、link等数据结构还要重写。也就是说，没有一个像C++ STL那样的一个泛型迭代器（这其中的一部分工作当然也需要通过重载操作符（如：++ 来实现）
