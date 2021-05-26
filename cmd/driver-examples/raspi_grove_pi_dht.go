// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"gobot-grovpi-platform/pkg/gobot-driver"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	gp := driver.NewGrovePiDriver(r)
	dht := driver.NewGroveTemperatureAndHumidityDriver(gp, "D7")
	red := gpio.NewLedDriver(gp, "D3")
	green := gpio.NewLedDriver(gp, "D4")
	blue := gpio.NewLedDriver(gp, "D5")
	lcd := i2c.NewGroveLcdDriver(r, i2c.WithAddress(1))

	work := func() {

		prevT := float32(0)
		prevH := float32(0)

		lcd.Clear()
		lcd.SetRGB(0, 0, 0)
		red.Off()
		green.Off()
		blue.Off()

		dht.On(driver.Temperature, func(d interface{}) {
			fmt.Printf("T: %v\n", d)
			red.Off()
			green.Off()
			blue.Off()

			if prevT == 0 {
				green.On()
			} else {
				diff := d.(float32) - prevT
				if diff > 0 {
					red.On()
				} else if diff < 0 {
					blue.On()
				} else {
					green.On()
				}
			}
			prevT = d.(float32)
		})

		dht.On(driver.Humidity, func(d interface{}) {
			fmt.Printf("H: %v\n", d)
			lcd.Clear()
			lcd.SetRGB(0, 0, 0)
			if prevT == 0 {
				lcd.SetRGB(0, 255, 0)
			} else {
				diff := d.(float32) - prevH
				if diff > 0 {
					lcd.SetRGB(255, 0, 0)
				} else if diff < 0 {
					lcd.SetRGB(0, 0, 255)
				} else {
					lcd.SetRGB(0, 255, 0)
				}
			}
			prevH = d.(float32)
		})
	}

	robot := gobot.NewRobot("rgbLcdBot",
		[]gobot.Connection{r},
		//		[]gobot.Device{gp, us},
		[]gobot.Device{gp, dht, red, green, blue, lcd},
		work,
	)

	robot.Start()
}
