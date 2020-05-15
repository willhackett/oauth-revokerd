package provider

type GCEEnvironment struct {
	ProjectName     string `env:"GCE_PROJECT_NAME" default:""`
	TagValue        string `env:"GCE_TAG_VALUE" default:""`
	ZonePattern     string `env:"GCE_ZONE_PATTERN" default:""`
	CredentialsFile string `env:"GCE_CREDENTIALS_FILE" default:""`
}

func GCEMapping(env GCEEnvironment) map[string]string {
	return map[string]string{
		"project_name":     env.ProjectName,
		"tag_value":        env.TagValue,
		"zone_pattern":     env.ZonePattern,
		"credentials_file": env.CredentialsFile,
	}
}
