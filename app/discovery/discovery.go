package discovery

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/go-discover"
	"github.com/willhackett/oauth-revokerd/app/discovery/provider"
)

// CloudDiscovery is the discovery provider type
type CloudDiscovery struct {
	config   *Config
	log      *log.Logger
	discover *discover.Discover
}

// Config is the structure of config for the CloudProvider
type Config struct {
	Provider string
	Args     string
}

func (c *CloudDiscovery) checkErrors() error {
	if c.config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if c.log == nil {
		return fmt.Errorf("logger cannot be nil")
	}

	_, ok := discover.Providers[c.config.Provider]
	if !ok {
		return fmt.Errorf("invalid provider: %s", c.config.Provider)
	}

	return nil
}

// Initialize is used in interfacing with Olric
func (c *CloudDiscovery) Initialize() error {
	if err := c.checkErrors(); err != nil {
		return err
	}

	m := map[string]discover.Provider{}

	provider, _ := discover.Providers[c.config.Provider]
	m[c.config.Provider] = provider

	opt := discover.WithProviders(m)
	d, err := discover.New(opt)
	if err != nil {
		return fmt.Errorf("discover.New returned an error: %s", err)
	}
	c.discover = d
	c.log.Printf("[INFO] Service discovery plugin is enabled, provider: %s", c.config.Provider)
	return nil
}

// SetLogger is used in interfacing with Olric
func (c *CloudDiscovery) SetLogger(l *log.Logger) {
	c.log = l
}

// SetConfig is used in interfacing with Olric
func (c *CloudDiscovery) SetConfig(cfg map[string]interface{}) error {
	provider, ok := cfg["provider"].(string)
	if !ok {
		return errors.New("Provider has not been supplied or is invalid")
	}
	args, ok := cfg["args"].(string)
	if !ok {
		args = ""
	}

	c.config = &Config{
		Provider: provider,
		Args:     args,
	}
	return nil
}

func (c *CloudDiscovery) getArgs() string {
	result := fmt.Sprintf("provider=%s", c.config.Provider)

	return result + provider.GenerateArgs(c.config.Args)
}

// DiscoverPeers is used in interfacing with Olric
func (c *CloudDiscovery) DiscoverPeers() ([]string, error) {
	peers, err := c.discover.Addrs(c.getArgs(), c.log)
	if err != nil {
		return nil, err
	}
	if len(peers) == 0 {
		return nil, fmt.Errorf("no peer found")
	}
	return peers, nil
}

// Register is unused but required to make this work
func (c *CloudDiscovery) Register() error { return nil }

// Deregister is unused but required to make this work
func (c *CloudDiscovery) Deregister() error { return nil }

// Close is unused but required to make this work
func (c *CloudDiscovery) Close() error { return nil }
