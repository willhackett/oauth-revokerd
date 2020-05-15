package provider

type AzureEnvironment struct {
	TenantID        string `env:"AZURE_TENANT_ID" default:""`
	ClientID        string `env:"AZURE_CLIENT_ID" default:""`
	SubscriptionID  string `env:"AZURE_SUBSCRIPTION_ID" default:""`
	SecretAccessKey string `env:"AZURE_SECRET_ACCESS_KEY" default:""`
	TagName         string `env:"AZURE_TAG_NAME" default:""`
	TagValue        string `env:"AZURE_TAG_VALUE" default:""`
	ResourceGroup   string `env:"AZURE_RESOURCE_GROUP" default:""`
	VmScaleSet      string `env:"AZURE_VM_SCALE_SET" default:""`
}

func AzureMapping(env AzureEnvironment) map[string]string {
	return map[string]string{
		"tenant_id":         env.TenantID,
		"client_id":         env.ClientID,
		"subscription_id":   env.SubscriptionID,
		"secret_access_key": env.SecretAccessKey,
		"tag_name":          env.TagName,
		"tag_value":         env.TagValue,
		"resource_group":    env.ResourceGroup,
		"vm_scale_set":      env.VmScaleSet,
	}
}
