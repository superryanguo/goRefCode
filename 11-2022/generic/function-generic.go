package fun

import "fmt"

//这个 map函数中我使用了两个类型 – T1 和 T2 ，
//T1 – 是需要处理数据的类型
//T2 – 是处理后的数据类型
//T1 和 T2 可以一样，也可以不一样。
//我们还有一个函数参数 –  func(T1) T2 意味着，进入的是 T1 类型的，出来的是 T2 类型的。
func gMap[T1 any, T2 any] (arr []T1, f func(T1) T2) []T2 {
  result := make([]T2, len(arr))
  for i, elem := range arr {
    result[i] = f(elem)
  }
  return result
}
nums := []int {0,1,2,3,4,5,6,7,8,9}
squares := gMap(nums, func (elem int) int {
  return elem * elem
})
print(squares)  //0 1 4 9 16 25 36 49 64 81

strs := []string{"Hao", "Chen", "MegaEase"}
upstrs := gMap(strs, func(s string) string  {
  return strings.ToUpper(s)
})
print(upstrs) // HAO CHEN MEGAEASE


dict := []string{"零", "壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"}
strs =  gMap(nums, func (elem int) string  {
  return  dict[elem]
})
print(strs) // 零 壹 贰 叁 肆 伍 陆 柒 捌 玖
nums := []int {0,1,2,3,4,5,6,7,8,9}


//reduce函数是把一堆数据合成一个
func gReduce[T1 any, T2 any] (arr []T1, init T2, f func(T2, T1) T2) T2 {
  result := init
  for _, elem := range arr {
    result = f(result, elem)
  }
  return result
}

sum := gReduce(nums, 0, func (result, elem int) int  {
    return result + elem
})
fmt.Printf("Sum = %d \n", sum)

//filter函数主要是用来做过滤的，把数据中一些符合条件（filter in）或是不符合条件（filter out）的数据过滤出来
//用户需要提从一个 bool 的函数，我们会把数据传给用户，然后用户只需要告诉我行还是不行，于是我们就会返回一个过滤好的数组给用户
func gFilter[T any] (arr []T, in bool, f func(T) bool) []T {
  result := []T{}
  for _, elem := range arr {
    choose := f(elem)
    if (in && choose) || (!in && !choose) {
      result = append(result, elem)
    }
  }
  return result
}

func gFilterIn[T any] (arr []T, f func(T) bool) []T {
  return gFilter(arr, true, f)
}

func gFilterOut[T any] (arr []T, f func(T) bool) []T {
  return gFilter(arr, false, f)
}
