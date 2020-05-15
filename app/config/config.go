package config

import (
	"gopkg.daemonl.com/envconf"
)

// Configuration contains the environment variables
type Configuration struct {
	// Port defines the port the HTTP server runs on
	Port int `env:"PORT" default:"8080"`
	// DiscoveryProvider allows you to specify a provider for Discovery
	DiscoveryProvider string `env:"DISCOVERY_PROVIDER" default:""`
	// MemberlistConfig
	MemberlistConfig string `env:"MEMBERLIST_CONFIG" default:"local"`
}

// Load will return the Configuration of the environment
func Load() (Configuration, error) {
	config := Configuration{}

	if err := envconf.Parse(&config); err != nil {
		return Configuration{}, err
	}

	return config, nil
}
