package manifests

import "k8s.io/apimachinery/pkg/util/intstr"

var (
	lokiPortInt = intstr.FromInt(3100)
)

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
