//其实并不需要 isChanClose 函数 !!!
//上面实现的 isChanClose 是可以判断出 channel 是否 close，但是适用场景优先，因为可能等你 isChanClose 判断的时候返回值 false，你以为 channel 还是正常的，但是下一刻 channel 被关闭了，这个时候往里面“写”数据就又会 panic ，如下：

if isChanClose( c ) {
    // 关闭的场景，exit
    return
}
// 未关闭的场景，继续执行（可能还是会 panic）
c <- x
因为判断之后还是有时间窗，所以 isChanClose 的适用还是有限，那么是否有更好的办法？

我们换一个思路，你其实并不是一定要判断 channel 是否 close，真正的目的是：安全的使用 channel，避免使用到已经关闭的 closed channel，从而导致 panic 。

这个问题的本质上是保证一个事件的时序，官方推荐通过 context 来配合使用，我们可以通过一个 ctx 变量来指明 close 事件，而不是直接去判断 channel 的一个状态。举个栗子：

select {
case <-ctx.Done():
    // ... exit
    return
case v, ok := <-c:
    // do something....
default:
    // do default ....
}
ctx.Done() 事件发生之后，我们就明确不去读 channel 的数据。

或者

select {
case <-ctx.Done():
    // ... exit
    return
default:
    // push
    c <- x
}
ctx.Done() 事件发生之后，我们就明确不写数据到 channel ，或者不从 channel 里读数据，那么保证这个时序即可。就一定不会有问题。

我们只需要确保一点：

触发时序保证：一定要先触发 ctx.Done() 事件，再去做 close channel 的操作，保证这个时序的才能保证 select 判断的时候没有问题；
只有这个时序，才能保证在获悉到 Done 事件的时候，一切还是安全的；
条件判断顺序：select 的 case 先判断 ctx.Done() 事件，这个很重要哦，否则很有可能先执行了 chan 的操作从而导致 panic 问题；

Image
怎么优雅关闭 chan ？

Image


方法一：panic-recover
关闭一个 channel 直接调用 close 即可，但是关闭一个已经关闭的 channel 会导致 panic，怎么办？panic-recover 配合使用即可。

func SafeClose(ch chan int) (closed bool) {
 defer func() {
  if recover() != nil {
   closed = false
  }
 }()
 // 如果 ch 是一个已经关闭的，会 panic 的，然后被 recover 捕捉到；
 close(ch)
 return true
}
这并不优雅。

方法二：sync.Once
可以使用 sync.Once 来确保 close 只执行一次。

type ChanMgr struct {
 C    chan int
 once sync.Once
}
func NewChanMgr() *ChanMgr {
 return &ChanMgr{C: make(chan int)}
}
func (cm *ChanMgr) SafeClose() {
 cm.once.Do(func() { close(cm.C) })
}
这看着还可以。

方法三：事件同步来解决
对于关闭 channel 这个我们有两个简要的原则：

永远不要尝试在读端关闭 channel ；
永远只允许一个 goroutine（比如，只用来执行关闭操作的一个 goroutine ）执行关闭操作；
可以使用 sync.WaitGroup 来同步这个关闭事件，遵守以上的原则，举几个例子：

第一个例子：一个 sender

package main

import "sync"

func main() {
 // channel 初始化
 c := make(chan int, 10)
 // 用来 recevivers 同步事件的
 wg := sync.WaitGroup{}

 // sender（写端）
 go func() {
  // 入队
  c <- 1
  // ...
  // 满足某些情况，则 close channel
  close(c)
 }()

 // receivers （读端）
 for i := 0; i < 10; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   // ... 处理 channel 里的数据
   for v := range c {
    _ = v
   }
  }()
 }
 // 等待所有的 receivers 完成；
 wg.Wait()
}
这里例子里面，我们在 sender 的 goroutine 关闭 channel，因为只有一个 sender，所以关闭自然是安全的。receiver 使用 WaitGroup 来同步事件，receiver 的 for 循环只有在 channel close 之后才会退出，主协程的 wg.Wait() 语句只有所有的 receivers 都完成才会返回。所以，事件的顺序是：

写端入队一个整形元素
关闭 channel
所有的读端安全退出
主协程返回
一切都是安全的。

第二个例子：多个 sender

package main

import (
 "context"
 "sync"
 "time"
)

func main() {
 // channel 初始化
 c := make(chan int, 10)
 // 用来 recevivers 同步事件的
 wg := sync.WaitGroup{}
 // 上下文
 ctx, cancel := context.WithCancel(context.TODO())

 // 专门关闭的协程
 go func() {
  time.Sleep(2 * time.Second)
  cancel()
  // ... 某种条件下，关闭 channel
  close(c)
 }()

 // senders（写端）
 for i := 0; i < 10; i++ {
  go func(ctx context.Context, id int) {
   select {
   case <-ctx.Done():
    return
   case c <- id: // 入队
    // ...
   }
  }(ctx, i)
 }

 // receivers（读端）
 for i := 0; i < 10; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   // ... 处理 channel 里的数据
   for v := range c {
    _ = v
   }
  }()
 }
 // 等待所有的 receivers 完成；
 wg.Wait()
}
这个例子我们看到有多个 sender 和 receiver ，这种情况我们还是要保证一点：close(ch) 操作的只能有一个人，我们单独抽出来一个 goroutine 来做这个事情，并且使用 context 来做事件同步，事件发生顺序是：

10 个写端协程（sender）运行，投递元素；
10 个读端协程（receiver）运行，读取元素；
2 分钟超时之后，单独协程执行 close(channel) 操作；
主协程返回；
一切都是安全的。

Image
总结

Image


channel 并没有直接提供判断是否 close 的接口，官方推荐使用 context 和 select 语法配合使用，事件通知的方式，达到优雅判断 channel 关闭的效果；
channel 关闭姿势也有讲究，永远不要尝试在读端关闭，永远保持一个关闭入口处，使用 sync.WaitGroup 和 context 实现事件同步，达到优雅关闭效果；
