package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rpio.Close()

	rasPi := raspi.NewAdaptor()
	sensor := i2c.NewGroveUltrasonicDriver(rasPi, "D2")
	led := gpio.NewLedDriver(rasPi, "11")
	buzzer := gpio.NewBuzzerDriver(rasPi, "13")

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{rasPi},
		[]gobot.Device{sensor, led, buzzer},
		work,
	)

	robot.Start()
}

func work() {
	for {
		distance, _ := sensor.Distance()
		fmt.Println("Distance:", distance)

		if distance > 50 {
			led.On()
			buzzer.On()
		} else {
			led.Off()
			buzzer.Off()
		}

		time.Sleep(500 * time.Millisecond)
	}
}

//- raspberry pi（我们使用的是raspberry pi 3）
//- 一个机器人底盘
//- 一个摄像头
//- 一个超声波传感器
//- 一个LED灯
//- 一个蜂鸣器
//第一步：连接硬件设备

//将超声波传感器连接到树莓派的GPIO引脚中，并将LED灯和蜂鸣器连接到树莓派的GPIO引脚中。此外，还需要将摄像头连接到树莓派的USB端口。

//此程序可以实现以下功能：

//- 通过超声波传感器测量距离，并在控制台上输出距离值。
//- 如果距离大于50厘米，则LED灯和蜂鸣器同时亮起，否则灭掉。

//第三步：运行程序

//在命令行中，运行以下命令以编译和运行程序：

//```bash
//go run main.go
//```

//通过摄像头可以观察到LED和蜂鸣器的反应。
