package platform

import (
	"errors"
	"gobot-grovepi-platform/pkg/config"
	driver "gobot-grovepi-platform/pkg/gobot-driver"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GrovePi struct {
	adaptor       *raspi.Adaptor
	robot         *gobot.Robot
	devicesByPin  map[string]gobot.Device
	devicesByName map[string]gobot.Device
	work          func()
}

const (
	GrovePiLEDDriverName              = "GroveLedDriver"
	GrovePiRotarySensorDriverName     = "GroveRotaryDriver"
	GrovePiButtonDriverName           = "GroveButtonDriver"
	GrovePiBuzzerDriverName           = "GroveBuzzerDriver"
	GrovePiSoundSensorDriverName      = "GroveSoundSensorDriver"
	GrovePiLightSensorDriverName      = "GroveLightSensorDriver"
	GrovePiRGBLCDPanelDriverName      = "GroveLcdDriver"
	GrovePiDHTSensorDriverName        = "GroveTemperatureAndHumidityDriver"
	GrovePiUltrasonicRangerDriverName = "GroveUltrasonicRangerDriver"

	SamplingIntervalPropertyName = "samplingInterval"

	RobotDefaultName = "gobot-grovepi-platform"
)

var (
	platformOnce     sync.Once
	platformInstance *GrovePi

	deviceFactories = map[string]func(*driver.GrovePiDriver, *config.DeviceConfig, *raspi.Adaptor) (gobot.Device, error){
		GrovePiLEDDriverName:              newLed,
		GrovePiRotarySensorDriverName:     newRotary,
		GrovePiButtonDriverName:           newButton,
		GrovePiBuzzerDriverName:           newBuzzer,
		GrovePiSoundSensorDriverName:      newSoundSensor,
		GrovePiLightSensorDriverName:      newLightSensor,
		GrovePiRGBLCDPanelDriverName:      newLcdPanel,
		GrovePiDHTSensorDriverName:        newDHT,
		GrovePiUltrasonicRangerDriverName: newUltrasonicRanger,
	}

	ErrorAlreadyInitialized = errors.New("already initialized")
	ErrorNotInitialized     = errors.New("not initialized yet")
	ErrorDriverNotSupported = errors.New("driver not supported")
	ErrorPinAlreadyInUse    = errors.New("pin already in use")
	ErrorInvalidI2CAddress  = errors.New("invalid I2C address")
	ErrorNameAlreadyInUse   = errors.New("name already in use")
)

func GetPlatform() *GrovePi {
	platformOnce.Do(func() {
		platformInstance = &GrovePi{
			adaptor:       raspi.NewAdaptor(),
			devicesByPin:  map[string]gobot.Device{},
			devicesByName: map[string]gobot.Device{},
			work:          func() {},
		}
	})
	return platformInstance
}

func (p *GrovePi) Init(conf *config.GrovePiConfig, w ...func()) error {

	if p.robot != nil {
		return ErrorAlreadyInitialized
	}

	gp := driver.NewGrovePiDriver(p.adaptor, i2c.WithBus(conf.Bus), i2c.WithAddress(conf.Address))
	devices, err := p.createDevices(gp, conf.Devices...)

	ds := make([]gobot.Device, 0)
	ds = append(ds, gp)
	ds = append(ds, devices...)

	if err != nil {
		return err
	}

	if w != nil && len(w) > 0 {
		p.work = w[0]
	}

	p.robot = gobot.NewRobot(RobotDefaultName,
		[]gobot.Connection{p.adaptor, gp},
		ds,
		p.work)
	return nil
}

func (p *GrovePi) Run() error {
	if p.robot == nil {
		return ErrorNotInitialized
	}

	mbot := gobot.NewMaster()
	mbot.AddRobot(p.robot)

	a := api.NewAPI(mbot)
	a.Debug()
	a.Start()

	err := mbot.Start()
	if err != nil {
		return err
	}
	return nil
}

func (p *GrovePi) createDevices(gp *driver.GrovePiDriver, conf ...*config.DeviceConfig) (devices []gobot.Device, err error) {
	if gp == nil || conf == nil || len(conf) <= 0 {
		return nil, nil
	}
	devices = make([]gobot.Device, 0)

	for _, cfg := range conf {
		if _, occupied := p.devicesByPin[cfg.Pin]; occupied {
			return nil, ErrorPinAlreadyInUse
		}
		if _, inUse := p.devicesByName[cfg.Name]; inUse {
			return nil, ErrorNameAlreadyInUse
		}
		createDevice, found := deviceFactories[cfg.Driver]
		if !found {
			return nil, ErrorDriverNotSupported
		}
		d, err := createDevice(gp, cfg, p.adaptor)
		if err != nil {
			return nil, err
		}
		if d != nil {
			d.SetName(cfg.Name)
			p.devicesByName[cfg.Name] = d
			p.devicesByPin[cfg.Pin] = d
			devices = append(devices, d)
		}
	}
	return devices, nil
}

//-------------------------------------------------------------------------------------------------------------
func newButton(gp *driver.GrovePiDriver, cfg *config.DeviceConfig, _ *raspi.Adaptor) (gobot.Device, error) {
	if gp == nil {
		return nil, ErrorNotInitialized
	}
	if durationStr, found := cfg.Properties[SamplingIntervalPropertyName]; found {
		duration, err := time.ParseDuration(durationStr.(string))
		if err != nil {
			return nil, err
		}
		return gpio.NewGroveButtonDriver(gp, cfg.Pin, duration), nil
	}
	return gpio.NewGroveButtonDriver(gp, cfg.Pin), nil
}

func newBuzzer(gp *driver.GrovePiDriver, cfg *config.DeviceConfig, _ *raspi.Adaptor) (gobot.Device, error) {
	if gp == nil {
		return nil, ErrorNotInitialized
	}
	return driver.NewGroveBuzzerDriver(gp, cfg.Pin), nil
}

func newDHT(gp *driver.GrovePiDriver, cfg *config.DeviceConfig, _ *raspi.Adaptor) (gobot.Device, error) {
	if gp == nil {
		return nil, ErrorNotInitialized
	}
	if durationStr, found := cfg.Properties[SamplingIntervalPropertyName]; found {
		duration, err := time.ParseDuration(durationStr.(string))
		if err != nil {
			return nil, err
		}
		return driver.NewGroveTemperatureAndHumidityDriver(gp, cfg.Pin, duration), nil
	}
	return driver.NewGroveTemperatureAndHumidityDriver(gp, cfg.Pin), nil
}

func newLed(gp *driver.GrovePiDriver, cfg *config.DeviceConfig, _ *raspi.Adaptor) (gobot.Device, error) {
	if gp == nil {
		return nil, ErrorNotInitialized
	}
	return gpio.NewGroveLedDriver(gp, cfg.Pin), nil
}

func newLightSensor(gp *driver.GrovePiDriver, cfg *config.DeviceConfig, _ *raspi.Adaptor) (gobot.Device, error) {
	if gp == nil {
		return nil, ErrorNotInitialized
	}
	if durationStr, found := cfg.Properties[SamplingIntervalPropertyName]; found {
		duration, err := time.ParseDuration(durationStr.(string))
		if err != nil {
			return nil, err
		}
		return aio.NewGroveLightSensorDriver(gp, cfg.Pin, duration), nil
	}
	return aio.NewGroveLightSensorDriver(gp, cfg.Pin), nil
}

func newRotary(gp *driver.GrovePiDriver, cfg *config.DeviceConfig, _ *raspi.Adaptor) (gobot.Device, error) {
	if gp == nil {
		return nil, ErrorNotInitialized
	}
	if durationStr, found := cfg.Properties[SamplingIntervalPropertyName]; found {
		duration, err := time.ParseDuration(durationStr.(string))
		if err != nil {
			return nil, err
		}
		return aio.NewGroveRotaryDriver(gp, cfg.Pin, duration), nil
	}
	return aio.NewGroveRotaryDriver(gp, cfg.Pin), nil
}

func newSoundSensor(gp *driver.GrovePiDriver, cfg *config.DeviceConfig, _ *raspi.Adaptor) (gobot.Device, error) {
	if gp == nil {
		return nil, ErrorNotInitialized
	}
	if durationStr, found := cfg.Properties[SamplingIntervalPropertyName]; found {
		duration, err := time.ParseDuration(durationStr.(string))
		if err != nil {
			return nil, err
		}
		return aio.NewGroveSoundSensorDriver(gp, cfg.Pin, duration), nil
	}
	return aio.NewGroveSoundSensorDriver(gp, cfg.Pin), nil
}

func newUltrasonicRanger(gp *driver.GrovePiDriver, cfg *config.DeviceConfig, _ *raspi.Adaptor) (gobot.Device, error) {
	if gp == nil {
		return nil, ErrorNotInitialized
	}
	if durationStr, found := cfg.Properties[SamplingIntervalPropertyName]; found {
		duration, err := time.ParseDuration(durationStr.(string))
		if err != nil {
			return nil, err
		}
		return aio.NewGroveLightSensorDriver(gp, cfg.Pin, duration), nil
	}
	return driver.NewGroveUltrasonicRangerDriver(gp, cfg.Pin), nil
}

func newLcdPanel(gp *driver.GrovePiDriver, cfg *config.DeviceConfig, raspi *raspi.Adaptor) (gobot.Device, error) {
	if gp == nil {
		return nil, ErrorNotInitialized
	}

	const i2cPrefix = "i2c-"

	addr := cfg.Pin
	if strings.HasPrefix(cfg.Pin, i2cPrefix) {
		addr = strings.TrimPrefix(cfg.Pin, i2cPrefix)
	}
	if len(addr) != 1 {
		return nil, ErrorInvalidI2CAddress
	}

	address, err := strconv.Atoi(addr)
	if err != nil {
		return nil, err
	}

	return i2c.NewGroveLcdDriver(raspi, i2c.WithAddress(address)), nil
}
