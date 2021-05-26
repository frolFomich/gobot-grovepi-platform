package gobot_driver

import (
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
)

const grovePiAddress = 0x04

// Commands format
const (
	CommandReadDigital    = 1
	CommandWriteDigital   = 2
	CommandReadAnalog     = 3
	CommandWriteAnalog    = 4
	CommandPinMode        = 5
	CommandReadUltrasonic = 7
	CommandReadDHT        = 40
)

// GrovePiDriver is a driver for the GrovePi+ for I²C bus interface.
// https://www.dexterindustries.com/grovepi/
//
// To use this driver with the GrovePi, it must be running the 1.3.0+ firmware.
// https://forum.dexterindustries.com/t/pre-release-of-grovepis-firmware-v1-3-0-open-to-testers/5119
//
type GrovePiDriver struct {
	name        string
	digitalPins map[int]string
	analogPins  map[int]string
	mutex       *sync.Mutex
	connector   i2c.Connector
	connection  i2c.Connection
	i2c.Config
}

// NewGrovePiDriver creates a new driver with specified i2c interface
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewGrovePiDriver(a i2c.Connector, options ...func(i2c.Config)) *GrovePiDriver {
	d := &GrovePiDriver{
		name:        gobot.DefaultName("GrovePi"),
		digitalPins: make(map[int]string),
		analogPins:  make(map[int]string),
		mutex:       &sync.Mutex{},
		connector:   a,
		Config:      i2c.NewConfig(),
	}

	for _, option := range options {
		option(d)
	}

	// TODO: add commands for API
	return d
}

// Name returns the Name for the Driver
func (d *GrovePiDriver) Name() string { return d.name }

// SetName sets the Name for the Driver
func (d *GrovePiDriver) SetName(n string) { d.name = n }

// Connection returns the connection for the Driver
func (d *GrovePiDriver) Connection() gobot.Connection { return d.connector.(gobot.Connection) }

// Start initialized the GrovePi
func (d *GrovePiDriver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(grovePiAddress)

	d.connection, err = d.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}

	return
}

// Halt returns true if devices is halted successfully
func (d *GrovePiDriver) Halt() (err error) { return }

// Connect is here to implement the Adaptor interface.
func (d *GrovePiDriver) Connect() (err error) {
	return
}

// Finalize is here to implement the Adaptor interface.
func (d *GrovePiDriver) Finalize() (err error) {
	return
}

// AnalogRead returns value from analog pin implementing the AnalogReader interface.
func (d *GrovePiDriver) AnalogRead(pin string) (value int, err error) {
	pin = getPin(pin)

	var pinNum int
	pinNum, err = strconv.Atoi(pin)
	if err != nil {
		return
	}

	value, err = d.readAnalog(byte(pinNum))

	return
}

// DigitalRead performs a read on a digital pin.
func (d *GrovePiDriver) DigitalRead(pin string) (val int, err error) {
	pin = getPin(pin)

	var pinNum int
	pinNum, err = strconv.Atoi(pin)
	if err != nil {
		return
	}

	if dir, ok := d.digitalPins[pinNum]; !ok || dir != "input" {
		err := d.PinMode(byte(pinNum), "input")
		if err != nil {
			return 0, err
		}
		d.digitalPins[pinNum] = "input"
	}

	val, err = d.readDigital(byte(pinNum))

	return
}

// UltrasonicRead performs a read on an ultrasonic pin.
func (d *GrovePiDriver) UltrasonicRead(pin string) (val int, err error) {
	pin = getPin(pin)

	var pinNum int
	pinNum, err = strconv.Atoi(pin)
	if err != nil {
		return
	}

	if dir, ok := d.digitalPins[pinNum]; !ok || dir != "input" {
		err := d.PinMode(byte(pinNum), "input")
		if err != nil {
			return 0, err
		}
		d.digitalPins[pinNum] = "input"
	}

	val, err = d.readUltrasonic(byte(pinNum))

	return
}

// DigitalWrite writes a value to a specific digital pin implementing the DigitalWriter interface.
func (d *GrovePiDriver) DigitalWrite(pin string, val byte) (err error) {
	pin = getPin(pin)

	var pinNum int
	pinNum, err = strconv.Atoi(pin)
	if err != nil {
		return
	}

	if dir, ok := d.digitalPins[pinNum]; !ok || dir != "output" {
		err := d.PinMode(byte(pinNum), "output")
		if err != nil {
			return err
		}
		d.digitalPins[pinNum] = "output"
	}

	err = d.writeDigital(byte(pinNum), val)

	return
}

// WriteAnalog writes PWM aka analog to the GrovePi. Not yet working.
func (d *GrovePiDriver) WriteAnalog(pin byte, val byte) error {
	buf := []byte{CommandWriteAnalog, pin, val, 0}
	_, err := d.connection.Write(buf)

	time.Sleep(2 * time.Millisecond)

	data := make([]byte, 1)
	_, err = d.connection.Read(data)

	return err
}

// PinMode sets the pin mode to input or output.
func (d *GrovePiDriver) PinMode(pin byte, mode string) error {
	var b []byte
	if mode == "output" {
		b = []byte{CommandPinMode, pin, 1, 0}
	} else {
		b = []byte{CommandPinMode, pin, 0, 0}
	}
	_, err := d.connection.Write(b)

	time.Sleep(2 * time.Millisecond)

	_, err = d.connection.ReadByte()

	return err
}

//readUltrasonic read ultrasonic ranger data
func (d *GrovePiDriver) readUltrasonic(pin byte) (int, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	err := d.PinMode(pin, "input")
	if err != nil {
		return 0, err
	}

	buf := []byte{CommandReadUltrasonic, pin, 0, 0}
	_, err = d.connection.Write(buf)

	time.Sleep(300 * time.Millisecond)

	_, err = d.connection.ReadByte()

	time.Sleep(100 * time.Millisecond)

	raw := make([]byte, 3)

	_, err = d.connection.Read(raw)

	if err != nil {
		return -1, err
	}

	return int(raw[1])*255 + int(raw[2]), nil
}

func (d *GrovePiDriver) ReadDHT(pin string) (float32, float32, error) {

	pin = getPin(pin)
	var pinNum int
	pinNum, err := strconv.Atoi(pin)
	if err != nil {
		return 0, 0, err
	}
	rawdata, err := d.readDHTRawData(byte(pinNum))
	if err != nil {
		return 0, 0, err
	}
	temperatureData := rawdata[1:5]

	tInt := int32(temperatureData[0]) | int32(temperatureData[1])<<8 | int32(temperatureData[2])<<16 | int32(temperatureData[3])<<24
	t := (*(*float32)(unsafe.Pointer(&tInt)))

	humidityData := rawdata[5:9]
	humInt := int32(humidityData[0]) | int32(humidityData[1])<<8 | int32(humidityData[2])<<16 | int32(humidityData[3])<<24
	h := (*(*float32)(unsafe.Pointer(&humInt)))

	return t, h, nil
}

func (d *GrovePiDriver) readDHTRawData(pin byte) ([]byte, error) {

	d.mutex.Lock()
	defer d.mutex.Unlock()

	err := d.PinMode(pin, "input")
	if err != nil {
		return nil, err
	}

	b := []byte{CommandReadDHT, pin, 0, 0}
	_, err = d.connection.Write(b)
	if err != nil {
		return nil, err
	}

	time.Sleep(600 * time.Millisecond)

	_, err = d.connection.ReadByte()
	if err != nil {
		return nil, err
	}
	time.Sleep(100 * time.Millisecond)

	raw := make([]byte, 9)
	_, err = d.connection.Read(raw)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func getPin(pin string) string {
	if len(pin) > 1 {
		if strings.ToUpper(pin[0:1]) == "A" || strings.ToUpper(pin[0:1]) == "D" {
			return pin[1:]
		}
	}

	return pin
}

// readAnalog reads analog value from the GrovePi.
func (d *GrovePiDriver) readAnalog(pin byte) (int, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	b := []byte{CommandReadAnalog, pin, 0, 0}
	_, err := d.connection.Write(b)
	if err != nil {
		return 0, err
	}

	time.Sleep(2 * time.Millisecond)

	data := make([]byte, 3)
	_, err = d.connection.Read(data)
	if err != nil || data[0] != CommandReadAnalog {
		return -1, err
	}

	v1 := int(data[1])
	v2 := int(data[2])
	return ((v1 * 256) + v2), nil
}

// readDigital reads digitally from the GrovePi.
func (d *GrovePiDriver) readDigital(pin byte) (val int, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	buf := []byte{CommandReadDigital, pin, 0, 0}
	_, err = d.connection.Write(buf)
	if err != nil {
		return
	}

	time.Sleep(2 * time.Millisecond)

	data := make([]byte, 2)
	_, err = d.connection.Read(data)
	if err != nil || data[0] != CommandReadDigital {
		return 0, err
	}

	return int(data[1]), err
}

// writeDigital writes digitally to the GrovePi.
func (d *GrovePiDriver) writeDigital(pin byte, val byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	buf := []byte{CommandWriteDigital, pin, val, 0}
	_, err := d.connection.Write(buf)

	time.Sleep(2 * time.Millisecond)

	_, err = d.connection.ReadByte()

	return err
}
