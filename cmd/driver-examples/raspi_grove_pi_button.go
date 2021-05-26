// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	driver "gobot-grovepi-platform/pkg/gobot-driver"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	gp := driver.NewGrovePiDriver(r)
	button := gpio.NewButtonDriver(gp, "D2")
	led := gpio.NewLedDriver(gp, "D3")

	work := func() {
		button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("button pressed")
			led.On()
		})

		button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("button released")
			led.Off()
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{r},
		[]gobot.Device{gp, button, led},
		work,
	)

	robot.Start()
}
