package config

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestNewAppConfig(t *testing.T) {
	ac := NewAppConfig(
		WithPlatformConfig(
			NewGrovePiConfig(
				WithGrovePiBus(0x04),
				WithGrovePiAddress(0x01),
				WithGrovePiDeviceConfig(
					NewDeviceConfig(
						WithDeviceName("redLed"),
						WithDeviceDriver("GrovePiLEDDriver"),
						WithDevicePin("D3"),
						WithDeviceProperty("color", "red"))),
				WithGrovePiDeviceConfig(
					NewDeviceConfig(
						WithDeviceName("greenLed"),
						WithDeviceDriver("GrovePiLEDDriver"),
						WithDevicePin("D4"),
						WithDeviceProperty("color", "green"))))))

	bytes, err := ac.ToYaml()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bytes))
}

func TestLoadAppConfigFromFile(t *testing.T) {
	wdir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	fname := strings.Join([]string{wdir, "../../config/app.yaml"}, string(os.PathSeparator))
	ac, err := LoadFromFile(fname)
	if err != nil {
		t.Error(err)
	}
	bytes, err := ac.ToYaml()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bytes))
}
