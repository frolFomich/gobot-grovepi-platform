package gobot_driver

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"time"
)

// GroveUltrasonicRangerDriver represents a Grove ultrasonic ranger Sensor
// Distance is reported in centimeters
type GroveUltrasonicRangerDriver struct {
	name     string
	halt     chan bool
	pin      string
	distance int
	interval time.Duration
	grovepi  *GrovePiDriver
	gobot.Eventer
	gobot.Commander
}

// NewGroveUltrasonicRangerDriver creates new instance of GroveUltrasonicRangerDriver
// Params:
//   c Connector - adaptor to use with this driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewGroveUltrasonicRangerDriver(gp *GrovePiDriver, pin string, i ...time.Duration) *GroveUltrasonicRangerDriver {
	drv := &GroveUltrasonicRangerDriver{
		name:      gobot.DefaultName("GroveUltrasonicRanger"),
		halt:      make(chan bool),
		pin:       pin,
		grovepi:   gp,
		interval:  10 * time.Millisecond,
		Eventer:   gobot.NewEventer(),
		Commander: gobot.NewCommander(),
	}

	if i != nil && len(i) > 0 {
		drv.interval = i[0]
	}

	drv.AddEvent(aio.Data)
	drv.AddEvent(aio.Error)

	drv.AddCommand("Read", func(params map[string]interface{}) interface{} {
		val, err := drv.Read()
		return map[string]interface{}{"val": val, "err": err}
	})

	return drv
}

// Name returns the Name for the Driver
func (d *GroveUltrasonicRangerDriver) Name() string { return d.name }

// SetName sets the Name for the Driver
func (d *GroveUltrasonicRangerDriver) SetName(n string) { d.name = n }

// Connection returns the connection for the Driver
func (d *GroveUltrasonicRangerDriver) Connection() gobot.Connection {
	return d.grovepi.Connection()
}

// Start initialized the GrovePi
func (d *GroveUltrasonicRangerDriver) Start() (err error) {

	d.distance = 0

	go func() {
		for {
			newValue, err := d.Read()

			if err != nil {
				d.Publish(aio.Error, err)
			} else if newValue != d.distance && d.distance != -1 {
				d.distance = newValue
				d.Publish(aio.Data, d.distance)
			}

			select {
			case <-time.After(d.interval):
			case <-d.halt:
				return
			}
		}
	}()

	return
}

// Halt returns true if devices is halted successfully
func (d *GroveUltrasonicRangerDriver) Halt() (err error) {
	d.halt <- true
	return
}

func (d *GroveUltrasonicRangerDriver) Distance() int {
	return d.distance
}

//Read performs a read on an ultrasonic ranger sensor.
func (d *GroveUltrasonicRangerDriver) Read() (val int, err error) {
	return d.grovepi.UltrasonicRead(d.pin)
}
