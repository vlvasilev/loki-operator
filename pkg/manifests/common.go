package manifests

func lokiConfigMapName(suffix string) string {
	return "vali-config-" + suffix
}

func commonLabels() map[string]string {
	return map[string]string{
		"gardener.cloud/role": "logging",
		"role":                "logging",
		"app":                 "loki",
	}
}
