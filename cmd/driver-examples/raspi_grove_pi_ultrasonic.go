// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"gobot-grovpi-platform/pkg/gobot-driver"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	gp := driver.NewGrovePiDriver(r)
	ur := driver.NewGroveUltrasonicRangerDriver(gp, "D6")

	work := func() {

		ur.On(aio.Data, func(data interface{}) {
			fmt.Printf("Dist: %v\n", data)
		})
	}

	robot := gobot.NewRobot("rgbLcdBot",
		[]gobot.Connection{r},
		[]gobot.Device{gp, ur},
		work,
	)

	robot.Start()
}
