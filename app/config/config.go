package config

import (
	"github.com/willhackett/oauth-revokerd/app/discovery"
	"gopkg.daemonl.com/envconf"
)

// Configuration contains the environment variables
type Configuration struct {
	Port int `env:"PORT" default:"8080"`
	// MulticastPort specified the port in which peers
	// should search for eachother on.
	MulticastPort int `env:"MULTICAST_PORT" default:"8484"`
	// MulticastAddress specifies the multicast address.
	// You should be able to use any of 224.0.0.0/4 or ff00::/8.
	// By default it uses the Simple Service Discovery Protocol
	// address (239.255.255.250 for IPv4).
	//
	// Specify ff02::c to use with IPv6
	MulticastAddress string `env:"MULTICAST_ADDRESS" default:"239.255.255.250"`
	// MulticastProtocol specified the IP version to use
	MulticastProtocol discovery.IPv
}

// Load will return the Configuration of the environment
func Load() (Configuration, error) {
	config := Configuration{}

	if err := envconf.Parse(&config); err != nil {
		return Configuration{}, err
	}

	return config, nil
}
