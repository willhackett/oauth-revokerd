package provider

type DigitalOceanEnvironment struct {
	Region   string `env:"DO_REGION" default:""`
	TagName  string `env:"DO_TAG_NAME" default:""`
	APIToken string `env:"DO_API_TOKEN" default:""`
}

func DigitalOceanMapping(env DigitalOceanEnvironment) map[string]string {
	return map[string]string{
		"region":    env.Region,
		"tag_name":  env.TagName,
		"api_token": env.APIToken,
	}
}
