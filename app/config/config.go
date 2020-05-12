package config

import (
	"gopkg.daemonl.com/envconf"
)

// Configuration contains the environment variables
type Configuration struct {
	Port int `env:"PORT" default:"8080"`
}

// Load will return the Configuration of the environment
func Load() (Configuration, error) {
	config := Configuration{}

	if err := envconf.Parse(&config); err != nil {
		return Configuration{}, err
	}

	return config, nil
}
