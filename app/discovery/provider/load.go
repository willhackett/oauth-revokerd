package provider

import (
	"fmt"

	"gopkg.daemonl.com/envconf"
)

func getParams(provider string) map[string]string {
	switch provider {
	case "aws":
		config := AWSEnvironment{}
		envconf.Parse(&config)
		return AwsMapping(config)
	case "azure":
		config := AzureEnvironment{}
		envconf.Parse(&config)
		return AzureMapping(config)

	case "digitalocean":
		config := DigitalOceanEnvironment{}
		envconf.Parse(&config)
		return DigitalOceanMapping(config)

	case "gce":
		config := GCEEnvironment{}
		envconf.Parse(&config)
		return GCEMapping(config)

	case "k8s":
		config := K8sEnvironment{}
		envconf.Parse(&config)
		return K8sMapping(config)

	case "mdns":
		config := MDNSEnvironment{}
		envconf.Parse(&config)
		return MDNSMapping(config)

	case "os":
		config := OSEnvironment{}
		envconf.Parse(&config)
		return OSMapping(config)
	default:
		return map[string]string{}
	}
}

// GenerateArgs creates the string for the arguments
func GenerateArgs(provider string) string {
	args := ""

	params := getParams(provider)

	for key, value := range params {
		if value != "" {
			args = args + fmt.Sprintf("%s=%s", key, value) + " "
		}
	}

	return args
}
