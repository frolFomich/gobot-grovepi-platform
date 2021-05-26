package gobot_driver

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

// GroveBuzzerDriver represents a buzzer
// with a Grove connector
type GroveBuzzerDriver struct {
	*gpio.BuzzerDriver
	gobot.Commander
}

// NewGroveBuzzerDriver return a new GroveBuzzerDriver given a DigitalWriter and pin.
func NewGroveBuzzerDriver(a gpio.DigitalWriter, pin string) *GroveBuzzerDriver {
	l := &GroveBuzzerDriver{
		BuzzerDriver: gpio.NewBuzzerDriver(a, pin),
		Commander:    gobot.NewCommander(),
	}

	l.AddCommand("Read", func(params map[string]interface{}) interface{} {
		isOn := l.State()
		str := "Off"
		if isOn {
			str = "On"
		}
		return map[string]interface{}{"state": str}
	})

	l.AddCommand("Tone", func(params map[string]interface{}) interface{} {
		hz := params["tone"].(float64)
		duration := params["duration"].(float64)
		return l.Tone(hz, duration)
	})

	l.AddCommand("Toggle", func(params map[string]interface{}) interface{} {
		return l.Toggle()
	})

	l.AddCommand("On", func(params map[string]interface{}) interface{} {
		return l.On()
	})

	l.AddCommand("Off", func(params map[string]interface{}) interface{} {
		return l.Off()
	})

	return l
}
