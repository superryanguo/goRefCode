package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type Visitor func(shape Shape)

//我们的代码中有一个Visitor的函数定义，还有一个Shape接口，其需要使用 Visitor函数做为参数。
//我们的实例的对象 Circle和 Rectangle实现了 Shape 的接口的 accept() 方法，这个方法就是等外面给我传递一个Visitor。

type Shape interface {
	accept(Visitor)
}

type Circle struct {
	Radius int
}

func (c Circle) accept(v Visitor) {
	v(c)
}

type Rectangle struct {
	Width, Heigh int
}

func (r Rectangle) accept(v Visitor) {
	v(r)
}

//其实，这段代码的目的就是想解耦 数据结构和 算法，使用 Strategy 模式也是可以完成的，而且会比较干净。但是在有些情况下，多个Visitor是来访问一个数据结构的不同部分，这种情况下，数据结构有点像一个数据库，而各个Visitor会成为一个个小应用。 kubectl就是这种情况。
func JsonVisitor(shape Shape) {
	bytes, err := json.Marshal(shape)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func XmlVisitor(shape Shape) {
	bytes, err := xml.Marshal(shape)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func main() {
	c := Circle{10}
	r := Rectangle{100, 200}
	shapes := []Shape{c, r}

	for _, s := range shapes {
		s.accept(JsonVisitor)
		s.accept(XmlVisitor)
	}

}
