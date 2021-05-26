package config

type GrovePiConfig struct {
	Bus     int             `yaml:"bus,omitempty"`
	Address int             `yaml:"address,omitempty"`
	Devices []*DeviceConfig `yaml:"devices,omitempty"`
}

type OptGrovePiConfig func(c *GrovePiConfig)

func NewGrovePiConfig(options ...OptGrovePiConfig) *GrovePiConfig {
	gpc := &GrovePiConfig{
		Devices: []*DeviceConfig{},
	}
	for _, opt := range options {
		opt(gpc)
	}
	return gpc
}

func WithGrovePiDeviceConfig(d interface{}) OptGrovePiConfig {
	return func(c *GrovePiConfig) {
		if c != nil && d != nil {
			if dvc, ok := d.(*DeviceConfig); ok {
				c.Devices = append(c.Devices, dvc)
			}
		}
	}
}

func WithGrovePiBus(b interface{}) OptGrovePiConfig {
	return func(c *GrovePiConfig) {
		if c != nil && b != nil {
			if bus, ok := b.(int); ok {
				c.Bus = bus
			}
		}
	}
}

func WithGrovePiAddress(a interface{}) OptGrovePiConfig {
	return func(c *GrovePiConfig) {
		if c != nil && a != nil {
			if adr, ok := a.(int); ok {
				c.Address = adr
			}
		}
	}
}
