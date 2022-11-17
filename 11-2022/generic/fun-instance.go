package fun

import "fmt"

type Employee struct {
  Name     string
  Age      int
  Vacation int
  Salary   float32
}

var employees = []Employee{
  {"Hao", 44, 0, 8000.5},
  {"Bob", 34, 10, 5000.5},
  {"Alice", 23, 5, 9000.0},
  {"Jack", 26, 0, 4000.0},
  {"Tom", 48, 9, 7500.75},
  {"Marry", 29, 0, 6000.0},
  {"Mike", 32, 8, 4000.3},
}
//我们想统一下所有员工的薪水，我们就可以使用前面的reduce函数
total_pay := gReduce(employees, 0.0, func(result float32, e Employee) float32 {
  return result + e.Salary
})
fmt.Printf("Total Salary: %0.2f\n", total_pay) // Total Salary: 43502.05

//一般来说，我们用 reduce 函数大多时候基本上是统计求和或是数个数，所以，是不是我们可以定义的更为直接一些？
//比如下面的这个 CountIf()，就比上面的 Reduce 干净了很多。

func gCountIf[T any](arr []T, f func(T) bool) int {
  cnt := 0
  for _, elem := range arr {
    if f(elem) {
      cnt += 1
    }
  }
  return cnt;
}
//我们做求和，我们也可以写一个Sum的泛型。

//处理 T 类型的数据，返回 U类型的结果
//然后，用户只需要给我一个需要统计的 T 的 U 类型的数据就可以了。
type Sumable interface {
  type int, int8, int16, int32, int64,
        uint, uint8, uint16, uint32, uint64,
        float32, float64
}

func gSum[T any, U Sumable](arr []T, f func(T) U) U {
  var sum U
  for _, elem := range arr {
    sum += f(elem)
  }
  return sum
}
//上面的代码我们动用了一个叫 Sumable 的接口，其限定了 U 类型，只能是 Sumable里的那些类型，也就是整型或浮点型，这个支持可以让我们的泛型代码更健壮一些

//统计年龄大于40岁的员工数
old := gCountIf(employees, func (e Employee) bool  {
    return e.Age > 40
})
fmt.Printf("old people(>40): %d\n", old)
// ld people(>40): 2
//把没有休假的员工过滤出来
no_vacation := gFilterIn(employees, func(e Employee) bool {
  return e.Vacation == 0
})
print(no_vacation)
//{Hao 44 0 8000.5} {Jack 26 0 4000} {Marry 29 0 6000}
