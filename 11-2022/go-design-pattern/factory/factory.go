package main

import "fmt"

//简单工厂
//Go 语言没有构造函数一说，所以一般会定义 NewXXX 函数来初始化相关类。NewXXX 函数返回接口时就是简单工厂模式。
// Printer 简单工厂要返回的接口类型
type Printer interface {
	Print(name string) string
}

func NewPrinter(lang string) Printer {
	switch lang {
	case "cn":
		return new(CnPrinter)
	case "en":
		return new(EnPrinter)
	default:
		return new(CnPrinter)
	}
}

// CnPrinter 是 Printer 接口的实现，它说中文
type CnPrinter struct{}

func (*CnPrinter) Print(name string) string {
	return fmt.Sprintf("你好, %s", name)
}

// EnPrinter 是 Printer 接口的实现，它说中文
type EnPrinter struct{}

func (*EnPrinter) Print(name string) string {
	return fmt.Sprintf("Hello, %s", name)
}

//示例源码运行Demo
func test() {
	printer := NewPrinter("en")
	fmt.Println(printer.Print("Bob"))
}

//目前这个场景里我们先提供两个语种的打印机，他们都是 Printer 接口的具体实现类型。
//简单工厂模式主要包含3个角色。

//简单工厂：是简单工厂模式的核心，负责实现创建所有实例的内部逻辑。工厂类的创建产品类的方法可以被外界直接调用，创建所需的产品对象。
//抽象产品：是简单工厂创建的所有对象的抽象父类/接口，负责描述所有实例的行为。
//具体产品：是简单工厂模式的创建目标。
//简单工厂的优点是，简单，缺点嘛，如果具体产品扩产，就必须修改工厂内部，增加Case，一旦产品过多就会导致简单工厂过于臃肿，为了解决这个问题，才有了下一级别的工厂模式--工厂方法。

//工厂方法
//工厂方法模式（Factory Method Pattern）又叫作多态性工厂模式，指的是定义一个创建对象的接口，但由实现这个接口的工厂类来决定实例化哪个产品类，工厂方法把类的实例化推迟到子类中进行。

//在工厂方法模式中，不再由单一的工厂类生产产品，而是由工厂类的子类实现具体产品的创建。因此，当增加一个产品时，只需增加一个相应的工厂类的子类, 以解决简单工厂生产太多产品时导致其内部代码臃肿（switch … case分支过多）的问题。

//下面举个简单的例子来理解工厂方法的设计思想，考虑有这样一个生产计算器的工厂，每类计算器产品都由一个子工厂负责生产。

//注意：Go中没有继承，所以这里说的工厂子类，其实是直接实现工厂接口的具体工厂类。

// OperatorFactory 工厂接口，由具体工厂类来实现
type OperatorFactory interface {
	Create() MathOperator
}

// MathOperator 实际产品实现的接口--表示数学运算器应该有哪些行为
type MathOperator interface {
	SetOperandA(int)
	SetOperandB(int)
	ComputeResult() int
}

//现在我们假定程序可以生产两类计算器，加法计算器和乘法计算器，也就是在工厂方法模式中，存在两个子类工厂。
//这两个子类工厂分别用来生产加法和乘法计算器，

//注意：这里为了理解，例子都很简单，真实场景下每个子类工厂创建产品实例的时候是可以放进去复杂逻辑的，不是简单的New一下。

// BaseOperator 是所有 Operator 的基类
// 封装公用方法，因为Go不支持继承，具体Operator类
// 只能组合它来实现类似继承的行为表现。
type BaseOperator struct {
	operandA, operandB int
}

func (o *BaseOperator) SetOperandA(operand int) {
	o.operandA = operand
}

func (o *BaseOperator) SetOperandB(operand int) {
	o.operandB = operand
}

//PlusOperatorFactory 是 PlusOperator 加法运算器的工厂类
type PlusOperatorFactory struct{}

func (pf *PlusOperatorFactory) Create() MathOperator {
	return &PlusOperator{
		BaseOperator: &BaseOperator{},
	}
}

//PlusOperator 实际的产品类--加法运算器
type PlusOperator struct {
	*BaseOperator
}

//ComputeResult 计算并获取结果
func (p *PlusOperator) ComputeResult() int {
	return p.operandA + p.operandB
}

// MultiOperatorFactory 是乘法运算器产品的工厂
type MultiOperatorFactory struct{}

func (mf *MultiOperatorFactory) Create() MathOperator {
	return &MultiOperator{
		BaseOperator: &BaseOperator{},
	}
}

// MultiOperator 实际的产品类--乘法运算器
type MultiOperator struct {
	*BaseOperator
}

func (m *MultiOperator) ComputeResult() int {
	return m.operandA * m.operandB
}

// 测试运行
//示例源码运行Demo
func main() {
	var factory OperatorFactory
	var mathOp MathOperator
	factory = &PlusOperatorFactory{}
	mathOp = factory.Create()
	mathOp.SetOperandB(3)
	mathOp.SetOperandA(2)
	fmt.Printf("Plus operation reuslt: %d\n", mathOp.ComputeResult())

	factory = &MultiOperatorFactory{}
	mathOp = factory.Create()
	mathOp.SetOperandB(3)
	mathOp.SetOperandA(2)
	fmt.Printf("Multiple operation reuslt: %d\n", mathOp.ComputeResult())
}

//抽象工厂模式：用于创建一系列相关的或者相互依赖的对象。

//为了更清晰地理解工厂方法模式和抽象工厂模式的区别，我们举一个品牌产品生态的例子。

//比如智能家居领域多家公司，现在有华为和小米，他们的工厂除了生产我们熟知的手机外，还会生产电视、空调这种家电设备。

//假如我们有幸作为他们工厂智能化管理软件的供应商，在软件系统里要对工厂进行抽象，这个时候就不能再用工厂方法这种设计模式了，因为工厂方法只能用来生产一种产品。
//目前抽象工厂有两个实际工厂类一个是华为的工厂，一个是小米的工厂，他们用来实际生产自家的产品设备。

// 示例源码运行Demo
type AbstractFactory interface {
	CreateTelevision() ITelevision
	CreateAirConditioner() IAirConditioner
}

type ITelevision interface {
	Watch()
}

type IAirConditioner interface {
	SetTemperature(int)
}

type HuaWeiFactory struct{}

func (hf *HuaWeiFactory) CreateTelevision() ITelevision {
	return &HuaWeiTV{}
}
func (hf *HuaWeiFactory) CreateAirConditioner() IAirConditioner {
	return &HuaWeiAirConditioner{}
}

type HuaWeiTV struct{}

func (ht *HuaWeiTV) Watch() {
	fmt.Println("Watch HuaWei TV")
}

type HuaWeiAirConditioner struct{}

func (ha *HuaWeiAirConditioner) SetTemperature(temp int) {
	fmt.Printf("HuaWei AirConditioner set temperature to %d ℃\n", temp)
}

type MiFactory struct{}

func (mf *MiFactory) CreateTelevision() ITelevision {
	return &MiTV{}
}
func (mf *MiFactory) CreateAirConditioner() IAirConditioner {
	return &MiAirConditioner{}
}

type MiTV struct{}

func (mt *MiTV) Watch() {
	fmt.Println("Watch HuaWei TV")
}

type MiAirConditioner struct{}

func (ma *MiAirConditioner) SetTemperature(temp int) {
	fmt.Printf("Mi AirConditioner set temperature to %d ℃\n", temp)
}

func main2() {
	var factory AbstractFactory
	var tv ITelevision
	var air IAirConditioner

	factory = &HuaWeiFactory{}
	tv = factory.CreateTelevision()
	air = factory.CreateAirConditioner()
	tv.Watch()
	air.SetTemperature(25)

	factory = &MiFactory{}
	tv = factory.CreateTelevision()
	air = factory.CreateAirConditioner()
	tv.Watch()
	air.SetTemperature(26)
}

//同样抽象工厂也具备工厂方法把产品的创建推迟到工厂子类去做的特性，假如未来加入了 VIVO 的产品，我们就可以通过再创建 VIVO 工厂子类来扩展。

//对于抽象工厂我们可以总结以下几点：

//当系统所提供的工厂所需生产的具体产品并不是一个简单的对象，而是多个位于不同产品等级结构中属于不同类型的具体产品时需要使用抽象工厂模式。
//抽象工厂模式是所有形式的工厂模式中最为抽象和最具一般性的一种形态。
//抽象工厂模式与工厂方法模式最大的区别在于，工厂方法模式针对的是一个产品等级结构，而抽象工厂模式则需要面对多个产品等级结构，一个工厂等级结构可以负责多个不同产品等级结构中的产品对象的创建 。
//当一个工厂等级结构可以创建出分属于不同产品等级结构的一个产品族中的所有对象时，抽象工厂模式比工厂方法模式更为简单、有效率。
//抽象工厂模式的优点

//当需要产品族时，抽象工厂可以保证客户端始终只使用同一个产品的产品族。
//抽象工厂增强了程序的可扩展性，对于新产品族的增加，只需实现一个新的具体工厂即可，不需要对已有代码进行修改，符合开闭原则。
//抽象工厂模式的缺点

//规定了所有可能被创建的产品集合，产品族中扩展新的产品困难，需要修改抽象工厂的接口。
//增加了系统的抽象性和理解难度。
