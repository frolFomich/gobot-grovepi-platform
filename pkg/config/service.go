package config

type ServiceConfig struct {
	Name       string                 `yaml:"name"`
	Properties map[string]interface{} `yaml:"config,omitempty"`
}
