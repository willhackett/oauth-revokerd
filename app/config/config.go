package config

import (
	"gopkg.daemonl.com/envconf"
)

// Configuration contains the environment variables
type Configuration struct {
	// Port is the port the API will run on
	Port int `env:"PORT" default:"8080"`
	// CachePort is the port the cache service runs on
	CachePort int `env:"CACHE_PORT" default:"8444"`
}

// Load will return the Configuration of the environment
func Load() (Configuration, error) {
	config := Configuration{}

	if err := envconf.Parse(&config); err != nil {
		return Configuration{}, err
	}

	return config, nil
}
