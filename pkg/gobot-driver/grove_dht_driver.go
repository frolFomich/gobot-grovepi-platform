package gobot_driver

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"time"
)

// GroveTemperatureAndHumidityDriver represents a Grove Temperature and Humidity Sensor
// Temperature is reported in Celsius degrees
type GroveTemperatureAndHumidityDriver struct {
	name     string
	halt     chan bool
	pin      string
	temp     float32
	humid    float32
	interval time.Duration
	grovepi  *GrovePiDriver
	gobot.Eventer
	gobot.Commander
}

const (
	Temperature = "temperature"
	Humidity    = "humidity"
)

// NewGroveTemperatureAndHumidityDriver creates new instance of GroveUltrasonicRangerDriver
// Params:
//   c Connector - adaptor to use with this driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewGroveTemperatureAndHumidityDriver(gp *GrovePiDriver, pin string, i ...time.Duration) *GroveTemperatureAndHumidityDriver {
	drv := &GroveTemperatureAndHumidityDriver{
		name:      gobot.DefaultName("TemperatureAndHumiditySensor"),
		halt:      make(chan bool),
		pin:       pin,
		grovepi:   gp,
		interval:  600 * time.Millisecond,
		Eventer:   gobot.NewEventer(),
		Commander: gobot.NewCommander(),
	}

	if i != nil && len(i) > 0 {
		drv.interval = i[0]
	}

	drv.AddEvent(Temperature)
	drv.AddEvent(Humidity)
	drv.AddEvent(aio.Error)

	drv.AddCommand("ReadTemperature", func(params map[string]interface{}) interface{} {
		return map[string]interface{}{
			"temperature": drv.Temperature()}
	})

	drv.AddCommand("ReadHumidity", func(params map[string]interface{}) interface{} {
		return map[string]interface{}{
			"humidity": drv.Humidity()}
	})

	return drv
}

// Name returns the Name for the Driver
func (d *GroveTemperatureAndHumidityDriver) Name() string { return d.name }

// SetName sets the Name for the Driver
func (d *GroveTemperatureAndHumidityDriver) SetName(n string) { d.name = n }

// Connection returns the connection for the Driver
func (d *GroveTemperatureAndHumidityDriver) Connection() gobot.Connection {
	return d.grovepi.Connection()
}

// Start initialized the GrovePi
func (d *GroveTemperatureAndHumidityDriver) Start() (err error) {

	d.temp = 0
	d.humid = 0

	go func() {
		for {

			newT, newH, err := d.Read()
			//fmt.Printf("Sample: t: %f, h: %f\n", newT, newH)

			if err != nil {
				d.Publish(aio.Error, err)
			} else {
				if newT != d.temp || d.temp == 0 {
					d.temp = newT
					d.Publish(Temperature, d.temp)
				}
				if newH != d.humid || d.humid == 0 {
					d.humid = newH
					d.Publish(Humidity, d.humid)
				}
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
func (d *GroveTemperatureAndHumidityDriver) Halt() (err error) {
	d.halt <- true
	return
}

func (d *GroveTemperatureAndHumidityDriver) Temperature() float32 {
	return d.temp
}

func (d *GroveTemperatureAndHumidityDriver) Humidity() float32 {
	return d.humid
}

//Read performs a read on temperature and humidity sensor.
func (d *GroveTemperatureAndHumidityDriver) Read() (float32, float32, error) {
	return d.grovepi.ReadDHT(d.pin)
}
