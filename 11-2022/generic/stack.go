package main

import "fmt"
//（注：目前Go的泛型函数不支持 export，所以只能使用第一个字符是小写的函数名）

type stack [T any] []T

func (s *stack[T]) push(elem T) {
  *s = append(*s, elem)
}

func (s *stack[T]) pop() {
  if len(*s) > 0 {
    *s = (*s)[:len(*s)-1]
  }
}
func (s *stack[T]) top() *T{
  if len(*s) > 0 {
    return &(*s)[len(*s)-1]
  }
  return nil
}

func (s *stack[T]) len() int{
  return len(*s)
}

func (s *stack[T]) print() {
  for _, elem := range *s {
    fmt.Print(elem)
    fmt.Print(" ")
  }
  fmt.Println("")
}
func main() {

  ss := stack[string]{}
  ss.push("Hello")
  ss.push("Hao")
  ss.push("Chen")
  ss.print()
  fmt.Printf("stack top is - %v\n", *(ss.top()))
  ss.pop()
  ss.pop()
  ss.print()


  ns := stack[int]{}
  ns.push(10)
  ns.push(20)
  ns.print()
  ns.pop()
  ns.print()
  *ns.top() += 1
  ns.print()
  ns.pop()
  fmt.Printf("stack top is - %v\n", ns.top())

}
