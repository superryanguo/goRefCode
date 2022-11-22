//观察者模式 (Observer Pattern)，定义对象间的一种一对多依赖关系，使得每当一个对象状态发生改变时，其相关依赖对象皆得到通知，依赖对象在收到通知后，可自行调用自身的处理程序，实现想要干的事情，比如更新自己的状态。

//发布者对观察者唯一了解的是它实现了某个接口（观察者接口）。这种松散耦合的设计最大限度地减少了对象之间的相互依赖，因此使我们能够构建灵活的系统。
//观察者模式也经常被叫做发布 - 订阅（Publish/Subscribe）模式、上面说的定义对象间的一种一对多依赖关系，一 - 指的是发布变更的主体对象，多 - 指的是订阅变更通知的订阅者对象。

//发布的状态变更信息会被包装到一个对象里，这个对象被称为事件，事件一般用英语过去式的语态来命名，比如用户注册时，用户模块在用户创建好后发布一个事件 UserCreated 或者 UserWasCreated 都行，这样从名字上就能看出，这是一个已经发生过的事件。

//事件发布给订阅者的过程，其实就是遍历一下已经注册的事件订阅者，逐个去调用订阅者实现的观察者接口方法，比如叫 handleEvent 之类的方法，这个方法的参数一般就是当前的事件对象。

//至于很多人会好奇的，事件的处理是不是异步的？主要看我们的需求是什么，一般情况下是同步的，即发布事件后，触发事件的方法会阻塞等到全部订阅者返回后再继续，当然也可以让订阅者的处理异步执行，完全看我们的需求。

//大部分场景下其实是同步执行的，单体架构会在一个数据库事务里持久化因为主体状态变更，而需要更改的所有实体类。

//微服务架构下常见的做法是有一个事件存储，订阅者接到事件通知后，会把事件先存到事件存储里，这两步也需要在一个事务里完成才能保证最终一致性，后面会再有其他线程把事件从事件存储里搞到消息设施里，发给其他服务，从而在微服务架构下实现各个位于不同服务的实体间的最终一致性。

//所以观察者模式，从程序效率上看，大多数情况下没啥提升，更多的是达到一种程序结构上的解耦，让代码不至于那么难维护。
package main

import "fmt"

// Subject 接口，它相当于是发布者的定义
type Subject interface {
	Subscribe(observer Observer)
	Notify(msg string)
}

// Observer 观察者接口
type Observer interface {
	Update(msg string)
}

// Subject 实现
type SubjectImpl struct {
	observers []Observer
}

// Subscribe 添加观察者（订阅者）
func (sub *SubjectImpl) Subscribe(observer Observer) {
	sub.observers = append(sub.observers, observer)
}

// Notify 发布通知
func (sub *SubjectImpl) Notify(msg string) {
	for _, o := range sub.observers {
		o.Update(msg)
	}
}

// Observer1 Observer1
type Observer1 struct{}

// Update 实现观察者接口
func (Observer1) Update(msg string) {
	fmt.Printf("Observer1: %s\n", msg)
}

// Observer2 Observer2
type Observer2 struct{}

// Update 实现观察者接口
func (Observer2) Update(msg string) {
	fmt.Printf("Observer2: %s\n", msg)
}

func main() {
	sub := &SubjectImpl{}
	sub.Subscribe(&Observer1{})
	sub.Subscribe(&Observer2{})
	sub.Notify("Hello")
}
//+++++++++++++++++++++++++++++++++++++
//下面我们实现一个支持以下功能的事件总线

//异步不阻塞
//支持任意参数值
// 代码来自https://lailin.xyz/post/observer.html
package eventbus

import (
 "fmt"
 "reflect"
 "sync"
)

// Bus Bus
type Bus interface {
 Subscribe(topic string, handler interface{}) error
 Publish(topic string, args ...interface{})
}

// AsyncEventBus 异步事件总线
type AsyncEventBus struct {
 handlers map[string][]reflect.Value
 lock     sync.Mutex
}

// NewAsyncEventBus new
func NewAsyncEventBus() *AsyncEventBus {
 return &AsyncEventBus{
  handlers: map[string][]reflect.Value{},
  lock:     sync.Mutex{},
 }
}

// Subscribe 订阅
func (bus *AsyncEventBus) Subscribe(topic string, f interface{}) error {
 bus.lock.Lock()
 defer bus.lock.Unlock()

 v := reflect.ValueOf(f)
 if v.Type().Kind() != reflect.Func {
  return fmt.Errorf("handler is not a function")
 }

 handler, ok := bus.handlers[topic]
 if !ok {
  handler = []reflect.Value{}
 }
 handler = append(handler, v)
 bus.handlers[topic] = handler

 return nil
}

// Publish 发布
// 这里异步执行，并且不会等待返回结果
func (bus *AsyncEventBus) Publish(topic string, args ...interface{}) {
 handlers, ok := bus.handlers[topic]
 if !ok {
  fmt.Println("not found handlers in topic:", topic)
  return
 }

 params := make([]reflect.Value, len(args))
 for i, arg := range args {
  params[i] = reflect.ValueOf(arg)
 }

 for i := range handlers {
  go handlers[i].Call(params)
 }
}
package eventbus

import (
 "fmt"
 "testing"
 "time"
)

func sub1(msg1, msg2 string) {
 time.Sleep(1 * time.Microsecond)
 fmt.Printf("sub1, %s %s\n", msg1, msg2)
}

func sub2(msg1, msg2 string) {
 fmt.Printf("sub2, %s %s\n", msg1, msg2)
}
func TestAsyncEventBus_Publish(t *testing.T) {
 bus := NewAsyncEventBus()
 bus.Subscribe("topic:1", sub1)
 bus.Subscribe("topic:1", sub2)
 bus.Publish("topic:1", "test1", "test2")
 bus.Publish("topic:1", "testA", "testB")
 time.Sleep(1 * time.Second)
}
