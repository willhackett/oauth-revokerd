package provider

type K8sEnvironment struct {
	KubeconfigPath string `env:"K8S_KUBECONFIG_PATH" default:""`
	Namespace      string `env:"K8S_NAMESPACE" default:"default"`
	LabelSelector  string `env:"K8S_LABEL_SELECTOR" default:""`
	FieldSelector  string `env:"K8S_FIELD_SELECTOR" default:""`
	HostNetwork    string `env:"K8S_HOST_NETWORK" default:""`
}

func K8sMapping(env K8sEnvironment) map[string]string {
	return map[string]string{
		"kubeconfig":     env.KubeconfigPath,
		"namespace":      env.Namespace,
		"label_selector": env.LabelSelector,
		"field_selector": env.FieldSelector,
		"host_network":   env.HostNetwork,
	}
}
