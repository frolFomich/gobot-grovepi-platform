// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	driver "gobot-grovepi-platform/pkg/gobot-driver"
	"time"

	"gobot.io/x/gobot"
	old "gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	gp := driver.NewGrovePiDriver(r)
	lcd := old.NewGroveLcdDriver(r, old.WithAddress(1))

	work := func() {
		temp := float32(0)
		humi := float32(0)
		gobot.Every(1*time.Second, func() {

			curr := time.Now()
			t, h, err := gp.ReadDHT(7)
			if err == nil {
				if h >= 0 {
					temp = t
					humi = h
				}
			} else {
				fmt.Printf("Error: %v\n", err)
			}

			lcd.Clear()
			lcd.SetRGB(0, 0, 0)

			lcd.SetPosition(0)
			lcd.Write(fmt.Sprintf("Time %d:%d:%d", curr.Hour(), curr.Minute(), curr.Second()))

			lcd.SetPosition(16)
			lcd.Write(fmt.Sprintf("T %v; H %v", temp, humi))

		})
	}

	robot := gobot.NewRobot("rgbLcdBot",
		[]gobot.Connection{r},
		[]gobot.Device{gp, lcd},
		work,
	)

	robot.Start()
}
