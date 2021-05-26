package config

import (
	"errors"
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
)

type AppConfig struct {
	Version  string           `yaml:"version"`
	Services []*ServiceConfig `yaml:"services,omitempty"`
	Platform *GrovePiConfig   `yaml:"platform"`
}

type OptAppConfig func(ac *AppConfig)

const (
	CurrentAppConfigVersion = "0.0.1"
)

var (
	ErrorInvalidConfigFile = errors.New("invalid config file reference")
)

func NewAppConfig(options ...OptAppConfig) *AppConfig {
	ac := &AppConfig{
		Version:  CurrentAppConfigVersion,
		Services: []*ServiceConfig{},
		Platform: NewGrovePiConfig(),
	}
	for _, opt := range options {
		opt(ac)
	}
	return ac
}

func LoadFromFile(conf string) (*AppConfig, error) {
	if conf == "" {
		return nil, ErrorInvalidConfigFile
	}

	bytes, err := ioutil.ReadFile(conf)
	if err != nil {
		return nil, err
	}

	ac := &AppConfig{}

	err = yaml.Unmarshal(bytes, ac)
	if err != nil {
		return nil, err
	}
	return ac, nil
}

func (a *AppConfig) ToYaml() ([]byte, error) {

	bytes, err := yaml.Marshal(a)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func WithServiceConfig(s interface{}) OptAppConfig {
	return func(ac *AppConfig) {
		if ac != nil {
			if svc, ok := s.(*ServiceConfig); ok {
				ac.Services = append(ac.Services, svc)
			}
		}
	}
}

func WithPlatformConfig(p interface{}) OptAppConfig {
	return func(ac *AppConfig) {
		if ac != nil {
			if platform, ok := p.(*GrovePiConfig); ok {
				ac.Platform = platform
			}
		}
	}
}
