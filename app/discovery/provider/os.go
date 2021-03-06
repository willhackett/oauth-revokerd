package provider

// OSEnvironment contains the environment vars for OpenStack config
type OSEnvironment struct {
	AuthURL   string `env:"OS_AUTH_URL" default:""`
	ProjectID string `env:"OS_PROJECT_ID" default:""`
	TagKey    string `env:"OS_TAG_KEY" default:""`
	TagValue  string `env:"OS_TAG_VALUE" default:""`
	UserName  string `env:"OS_USER_NAME" default:""`
	Password  string `env:"OS_PASSWORD" default:""`
	Token     string `env:"OS_TOKEN" default:""`
	Insecure  string `env:"OS_INSECURE" default:""`
}

// OSMapping contains the mappings for the env vars
func OSMapping(env OSEnvironment) map[string]string {
	return map[string]string{
		"auth_url":   env.AuthURL,
		"project_id": env.ProjectID,
		"tag_key":    env.TagKey,
		"tag_value":  env.TagValue,
		"user_name":  env.UserName,
		"password":   env.Password,
		"token":      env.Token,
		"insecure":   env.Insecure,
	}
}
