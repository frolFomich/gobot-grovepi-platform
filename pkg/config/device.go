package config

type DeviceConfig struct {
	Name       string                 `yaml:"name"`
	Driver     string                 `yaml:"driver"`
	Pin        string                 `yaml:"pin"`
	Properties map[string]interface{} `yaml:"config,omitempty"`
}

type OptDeviceConfig func(d *DeviceConfig)

func NewDeviceConfig(options ...OptDeviceConfig) *DeviceConfig {
	dc := &DeviceConfig{
		Driver:     "unknown",
		Pin:        "unknown",
		Properties: map[string]interface{}{},
	}
	for _, opt := range options {
		opt(dc)
	}
	return dc
}

func WithDeviceDriver(dd interface{}) OptDeviceConfig {
	return func(d *DeviceConfig) {
		if d != nil && dd != nil {
			if drv, ok := dd.(string); ok {
				d.Driver = drv
			}
		}
	}
}

func WithDeviceName(n interface{}) OptDeviceConfig {
	return func(d *DeviceConfig) {
		if d != nil && n != nil {
			if name, ok := n.(string); ok {
				d.Name = name
			}
		}
	}
}

func WithDevicePin(p interface{}) OptDeviceConfig {
	return func(d *DeviceConfig) {
		if d != nil && p != nil {
			if pin, ok := p.(string); ok {
				d.Pin = pin
			}
		}
	}
}

func WithDeviceProperty(k string, v interface{}) OptDeviceConfig {
	return func(d *DeviceConfig) {
		if d != nil && k != "" {
			d.Properties[k] = v
		}
	}
}

func WithDeviceProperties(dp interface{}) OptDeviceConfig {
	return func(d *DeviceConfig) {
		if d != nil && dp != nil {
			if m, ok := dp.(map[string]interface{}); ok {
				for k, v := range m {
					d.Properties[k] = v
				}
			}
		}
	}
}
