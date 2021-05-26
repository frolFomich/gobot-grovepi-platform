// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	board := raspi.NewAdaptor()
	gp := i2c.NewGrovePiDriver(board)
	sensor := aio.NewGroveRotaryDriver(gp, "A1")
	lcd := i2c.NewGroveLcdDriver(board, i2c.WithAddress(1))

	work := func() {
		sensor.On(aio.Data, func(data interface{}) {
			fmt.Println("sensor", data)
			lcd.Clear()
			lcd.SetRGB(0, 0, 0)
			lcd.SetPosition(0)
			lcd.Write(fmt.Sprintf("%v", data))
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{board},
		[]gobot.Device{gp, sensor, lcd},
		work,
	)

	robot.Start()
}
