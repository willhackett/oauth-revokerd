package config

import (
	"gopkg.daemonl.com/envconf"
)

type Configuration struct {
	Port int `env:"PORT" default:"8000"`
}

func Load() (Configuration, error) {
	config := Configuration{}

	if err := envconf.Parse(&config); err != nil {
		return Configuration{}, err
	}

	return config, nil
}
