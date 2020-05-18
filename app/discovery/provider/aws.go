package provider

// AWSEnvironment contains environment variables for AWS
type AWSEnvironment struct {
	Region          string `env:"AWS_REGION" default:""`
	TagKey          string `env:"AWS_TAG_KEY" default:"default"`
	TagValue        string `env:"AWS_TAG_VALUE" default:""`
	AddrType        string `env:"AWS_ADDRESS_TYPE" default:"private_v4"`
	AccessKeyID     string `env:"AWS_HOST_NETWORK" default:""`
	SecretAccessKey string `env:"AWS_HOST_NETWORK" default:""`
}

// AwsMapping creates the values needed by orlic
func AwsMapping(env AWSEnvironment) map[string]string {
	return map[string]string{
		"region":            env.Region,
		"tag_key":           env.TagKey,
		"tag_value":         env.TagValue,
		"addr_type":         env.AddrType,
		"access_key_id":     env.AccessKeyID,
		"secret_access_key": env.SecretAccessKey,
	}
}
