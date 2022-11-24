package main

import "fmt"

//k8s相关背景
//接下来，我们再来了解一下相关的知识背景：

//对于Kubernetes，其抽象了很多种的Resource，比如：Pod, ReplicaSet, ConfigMap, Volumes, Namespace, Roles …. 种类非常繁多，这些东西构成为了Kubernetes的数据模型（点击 Kubernetes Resources 地图 查看其有多复杂）
//kubectl 是Kubernetes中的一个客户端命令，操作人员用这个命令来操作Kubernetes。kubectl 会联系到 Kubernetes 的API Server，API Server会联系每个节点上的 kubelet ，从而达到控制每个结点。
//kubectl 主要的工作是处理用户提交的东西（包括，命令行参数，yaml文件等），然后其会把用户提交的这些东西组织成一个数据结构体，然后把其发送给 API Server。
//相关的源代码在 src/k8s.io/cli-runtime/pkg/resource/visitor.go 中（源码链接）
//kubectl 的代码比较复杂，不过，其本原理简单来说，它从命令行和yaml文件中获取信息，通过Builder模式并把其转成一系列的资源，最后用 Visitor 模式模式来迭代处理这些Reources。

//下面我们来看看 kubectl 的实现，为了简化，我用一个小的示例来表明 ，而不是直接分析复杂的源码。

//kubectl的实现方法
//Visitor模式定义
//首先，kubectl 主要是用来处理 Info结构体，下面是相关的定义：

type VisitorFunc func(*Info, error) error
type Visitor interface {
	Visit(VisitorFunc) error
}
type Info struct {
	Namespace   string
	Name        string
	OtherThings string
}

func (info *Info) Visit(fn VisitorFunc) error {
	return fn(info, nil)
}

//我们可以看到，

//有一个 VisitorFunc 的函数类型的定义
//一个 Visitor 的接口，其中需要 Visit(VisitorFunc) error  的方法（这就像是我们上面那个例子的 Shape ）
//最后，为Info 实现 Visitor 接口中的 Visit() 方法，实现就是直接调用传进来的方法（与前面的例子相仿）
//我们再来定义几种不同类型的 Visitor。

//Name Visitor
//这个Visitor 主要是用来访问 Info 结构中的 Name 和 NameSpace 成员

type NameVisitor struct {
	visitor Visitor
}

func (v NameVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("NameVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> Name=%s, NameSpace=%s\n", info.Name, info.Namespace)
		}
		fmt.Println("NameVisitor() after call function")
		return err
	})
}

//我们可以看到，上面的代码：

//声明了一个 NameVisitor 的结构体，这个结构体里有一个 Visitor 接口成员，这里意味着多态。
//在实现 Visit() 方法时，其调用了自己结构体内的那个 Visitor的 Visitor() 方法，这其实是一种修饰器的模式，用另一个Visitor修饰了自己（关于修饰器模式，参看《Go编程模式：修饰器》）
//Other Visitor
//这个Visitor主要用来访问 Info 结构中的 OtherThings 成员

type OtherThingsVisitor struct {
	visitor Visitor
}

func (v OtherThingsVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("OtherThingsVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> OtherThings=%s\n", info.OtherThings)
		}
		fmt.Println("OtherThingsVisitor() after call function")
		return err
	})
}

//实现逻辑同上，我就不再重新讲了

//Log Visitor
type LogVisitor struct {
	visitor Visitor
}

func (v LogVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("LogVisitor() before call function")
		err = fn(info, err)
		fmt.Println("LogVisitor() after call function")
		return err
	})
}

//使用方代码
//现在我们看看如果使用上面的代码：

func main() {
	info := Info{}
	var v Visitor = &info
	v = LogVisitor{v}
	v = NameVisitor{v}
	v = OtherThingsVisitor{v}
	loadFile := func(info *Info, err error) error {
		info.Name = "Hao Chen"
		info.Namespace = "MegaEase"
		info.OtherThings = "We are running as remote team."
		return nil
	}
	v.Visit(loadFile)
}

//上面的代码，我们可以看到

//Visitor们一层套一层
//我用 loadFile 假装从文件中读如数据
//最后一条 v.Visit(loadfile) 我们上面的代码就全部开始激活工作了。
//上面的代码输出如下的信息，你可以看到代码的执行顺序是怎么执行起来了

//LogVisitor() before call function
//NameVisitor() before call function
//OtherThingsVisitor() before call function
//==> OtherThings=We are running as remote team.
//OtherThingsVisitor() after call function
//==> Name=Hao Chen, NameSpace=MegaEase
//NameVisitor() after call function
//LogVisitor() after call function
//我们可以看到，上面的代码有以下几种功效：

//解耦了数据和程序。
//使用了修饰器模式
//还做出来pipeline的模式
//所以，其实，我们是可以把上面的代码重构一下的。

//Visitor修饰器
//下面，我们用修饰器模式来重构一下上面的代码。

type DecoratedVisitor struct {
	visitor    Visitor
	decorators []VisitorFunc
}

func NewDecoratedVisitor(v Visitor, fn ...VisitorFunc) Visitor {
	if len(fn) == 0 {
		return v
	}
	return DecoratedVisitor{v, fn}
}

//Visit implements Visitor
func (v DecoratedVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		if err != nil {
			return err
		}
		if err := fn(info, nil); err != nil {
			return err
		}
		for i := range v.decorators {
			if err := v.decorators[i](info, nil); err != nil {
				return err
			}
		}
		return nil
	})
}

func callDec() {
	info := Info{}
	var v Visitor = &info
	v = NewDecoratedVisitor(v, NameVisitor, OtherVisitor)
	v.Visit(LoadFile)
	//用一个 DecoratedVisitor 的结构来存放所有的VistorFunc函数
	//NewDecoratedVisitor 可以把所有的 VisitorFunc转给它，构造 DecoratedVisitor 对象。
	//DecoratedVisitor实现了 Visit() 方法，里面就是来做一个for-loop，顺着调用所有的 VisitorFunc
	//于是，我们的代码就可以这样运作了：
}

//上面的代码并不复杂，

//是不是比之前的那个简单？注意，这个DecoratedVisitor 同样可以成为一个Visitor来使用。

//好，上面的这些代码全部存在于 kubectl 的代码中，你看懂了这里面的代码逻辑，相信你也能够看懂 kubectl 的代码了。
