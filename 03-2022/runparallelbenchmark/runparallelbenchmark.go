//go test -bench=. -cpu=4 -count=2 -benchtime=10s
//Benchmark 执行的总时间一定是大于 -benchtime 设置的时间的。
//~/go1.17/src/testing/benchmark.go
package main

import "fmt"

func main() {
	fmt.Println("vim-go")
}

//-benchtime t
//Run enough iterations of each benchmark to take t, specified
//as a time.Duration (for example, -benchtime 1h30s).
//The default is 1 second (1s).
//The special syntax Nx means to run the benchmark N times
//(for example, -benchtime 100x).

//-count n
//Run each test and benchmark n times (default 1).
//If -cpu is set, run n times for each GOMAXPROCS value.
//Examples are always run once.

//RunParallel runs a benchmark in parallel. It creates multiple goroutines and distributes b.N iterations among them. The number of goroutines defaults to GOMAXPROCS. To increase parallelism for non-CPU-bound benchmarks, call SetParallelism before RunParallel. RunParallel is usually used with the go test -cpu flag.
func BenchmarkPooledObject(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			object := pool.Get().(*MyObject)
			Consume(object)
			// 用完了放回对象池
			object.Reset()
			pool.Put(object)
		}
	})
}

func BenchmarkNewObject(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			object := &MyObject{
				Name: "hello",
				Age:  2,
			}
			Consume(object)
		}
	})
}

func (b *B) runN(n int) {
	benchmarkLock.Lock()
	defer benchmarkLock.Unlock()
	defer b.runCleanup(normalPanic)
	// 注意看这里，帮我们GC了
	runtime.GC()
	b.raceErrors = -race.Errors()
	b.N = n
	b.parallelism = 1
	// 重置计时器
	b.ResetTimer()
	// 开始计时
	b.StartTimer()
	// 执行 benchmark 方法
	b.benchFunc(b)
	// 停止计时
	b.StopTimer()
	b.previousN = n
	b.previousDuration = b.duration
	b.raceErrors += race.Errors()
	if b.raceErrors > 0 {
		b.Errorf("race detected during execution of benchmark")
	}
}

func (b *B) launch() {
   //...
 // 标注①
 if b.benchTime.n > 0 {
  // We already ran a single iteration in run1.
  // If -benchtime=1x was requested, use that result.
  if b.benchTime.n > 1 {
   b.runN(b.benchTime.n)
  }
 } else {
  d := b.benchTime.d
   // 标注②
  for n := int64(1); !b.failed && b.duration < d && n < 1e9; {
   last := n
   goalns := d.Nanoseconds()
   prevIters := int64(b.N)
   prevns := b.duration.Nanoseconds()
   if prevns <= 0 {
    prevns = 1
   }
    // 标注③
   n = goalns * prevIters / prevns
   // Run more iterations than we think we'll need (1.2x).
   // 标注④
   n += n / 5
   // Don't grow too fast in case we had timing errors previously.
   // 标注⑤
   n = min(n, 100*last)
   // Be sure to run at least one more than last time.
   // 标注⑥
   n = max(n, last+1)
   // Don't run more than 1e9 times. (This also keeps n in int range on 32 bit platforms.)
   // 标注⑦
   n = min(n, 1e9)
   // 标注⑧
   b.runN(int(n))
  }
 }
 b.result = BenchmarkResult{b.N, b.duration, b.bytes, b.netAllocs, b.netBytes, b.extra}
//}
//核心都标了序号，这里来解释下：

//标注①：Go 的 Benchmark 执行两种传参，执行次数和执行时间限制，我用的是执行时间，也可以用 -benchtime=1000x来表示需要测试1000次。

//标注②：这里是当设置了执行时间限制时，判断时间是否足够的条件，可以看到除了时间的判断外，还有 n < 1e9 的限制，也就是最多执行次数是 1e9，也就是 1000000000，这解释了上面的一个困惑，为啥执行时间还比设置的 benchtime 小。因为 Go 限制了最大执行次数为 1e9，并不是设置多少就是多少，还有个上限。

//标注③到⑧: Go 是如何知道 n 取多少时，时间刚好符合我们设置的 benchtime？答案是试探！

//n 从1 开始试探，执行1次后，根据执行时间来估算 n。n = goalns * prevIters / prevns，这就是估算公式，goalns 是设置的执行时间（单位纳秒），prevIters 是上次执行次数，prevns 是上一次执行时间（纳秒）

//根据上次执行的时间和目标设定的执行总时间，计算出需要执行的次数，大概是这样吧：

//目标执行次数 = 执行目标时间 / (上次执行时间 / 上次执行次数)

//化简下得到：

//目标执行次数 = 执行目标时间 * 上次执行次数 / 上次执行时间，这不就是上面那个公式~

//目标执行次数 n 的计算，源码中还做了一些其他处理：

//标注④：让实际执行次数大概是目标执行次数的1.2倍，万一达不到目标时间不是有点尴尬？索性多跑一会
//标注⑤：也不能让 n 增长的太快了，设置个最大增长幅度为100倍，当 n 增长太快时，被测试方法一定是执行时间很短，误差可能较大，缓慢增长好测出真实的水平
//标注⑥：n 不能原地踏步，怎么也得+1
//标注⑦：n 得设置个 1e9 的上限，这是为了在32位系统上不要溢出
//Go Benchmark 的执行原理大致摸清了，但我们要的答案还未浮出水面。

//接着我对 Benchmark 进行了断点调试。

//首先是 -benchtime=10s

//发现 n 的试探增长是 1，100，10000，1000000，100000000，1000000000，最终 n 是 1000000000

//这说明我们的执行方法耗时很短，执行次数达到了上限。

//再看-benchtime=150s，开始还挺正常：

//n 增长是 1，100，10000，1000000，100000000，但后一个出现了问题：
//n 居然变成了负数！显然这是溢出了。

//n = goalns * prevIters / prevns 这个公式，在目标执行时间（goalns）很大，测试方法执行时间（prevns）很短时，会导致 n 溢出！

//溢出有什么后果呢？

//后面的 n = min(n, 100*last) 永远等于 100000000 了，但还有 n = max(n, last+1) 保证，所以 n 还是在增加，不过很慢，每次都只 +1，所以后续试探的 n 序列为 100000001，100000002，100000003....

//这就导致了 n 很难达到 1e9 的上限，而且总的执行耗时也很难达到设定的预期时间，所以测试程序会一直跑~直到超时！
