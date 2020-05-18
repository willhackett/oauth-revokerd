package provider

// MDNSEnvironment contains the env vars to configure mDNS
type MDNSEnvironment struct {
	Service string `env:"MDNS_SERVICE_NAME" default:"oauth-revokerd"`
	Domain  string `env:"MDNS_DOMAIN" default:"local"`
	Timeout string `env:"MDNS_TIMEOUT" default:"5s"`
	V6      string `env:"MDNS_V6" default:"true"`
	V4      string `env:"MDNS_V4" default:"true"`
}

// MDNSMapping contains the field mappings for the environment
func MDNSMapping(env MDNSEnvironment) map[string]string {
	return map[string]string{
		"service": env.Service,
		"domain":  env.Domain,
		"timeout": env.Timeout,
		"v6":      env.V6,
		"v4":      env.V4,
	}
}
