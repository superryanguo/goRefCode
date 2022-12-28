package main

import "fmt"

//现在我们只要实例化一个Car的实例，在实例上面调用Drive()方法就能让车开起来，不过如果我们的驾驶员现在还是个未成年，那么在地球的大部分国家都是不允许开车的，如果在开车时要加一个驾驶员的年龄限制，我们该怎么办呢？

//给Car结构体加一个Age字段显然是不合理的，因为我们要表示的驾驶员的年龄而不是车的车龄。同理驾驶员年龄的判断我们也不应该加在 Car 实现的 Drive() 方法里， 这样会导致每个实现 Vehicle 接口的类型都要在自己的 Drive() 方法里加上类似的判断。

//这个时候通常的做法是，加一个表示驾驶员的类型 Driver。

type Car struct{}

type Vehicle interface {
	Drive()
}

type Driver struct {
	Age int
}

func (c *Car) Drive() {
	fmt.Println("Car is being driven")
}

type CarProxy struct {
	vehicle Vehicle
	driver  *Driver
}

func NewCarProxy(driver *Driver) *CarProxy {
	return &CarProxy{&Car{}, driver}
}

//样的话我们接可以通过，用包装类型代理vehicle属性的 Drive() 行为时，给它加上驾驶员的年龄限制。
func (c *CarProxy) Drive() {
	if c.driver.Age >= 16 {
		c.vehicle.Drive()
	} else {
		fmt.Println("Driver too young!")
	}
}

//现在我们通过代理模式给 Car 类型的 Drive() 行为扩充了检查驾驶员的行为，下面我们执行一下程序试试效果。
func main() {
	car := NewCarProxy(&Driver{12})
	car.Drive() // 输出 Driver too young!
	car2 := NewCarProxy(&Driver{22})
	car2.Drive() // 输出 Car is being driven
}
