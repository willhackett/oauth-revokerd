package provider

// DigitalOceanEnvironment contains environment variables for DO
type DigitalOceanEnvironment struct {
	Region   string `env:"DO_REGION" default:""`
	TagName  string `env:"DO_TAG_NAME" default:""`
	APIToken string `env:"DO_API_TOKEN" default:""`
}

// DigitalOceanMapping contains the mappings for the env vars
func DigitalOceanMapping(env DigitalOceanEnvironment) map[string]string {
	return map[string]string{
		"region":    env.Region,
		"tag_name":  env.TagName,
		"api_token": env.APIToken,
	}
}
