package manifests

import (
	"fmt"

	monitoring1alpha1 "github.com/vlvasilev/loki-operator/api/v1alpha1"

	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/pointer"
)

// Options is a set of configuration values to use when building manifests such as resource sizes, etc.
// Most of this should be provided - either directly or indirectly - by the user.
type Options struct {
	Name              string
	Namespace         string
	Vali              ValiOptions
	Curator           CuratorOptions
	ConfigSHA         string
	PriorityClassName *string
}

type ValiOptions struct {
	Spec     *monitoring1alpha1.ValiSpec
	Image    string
	Resorces corev1.ResourceRequirements
	PVCSize  *resource.Quantity
}

var defaultValiStorage = resource.MustParse("30Gi")
var defaultValiOptions = ValiOptions{
	Spec: &monitoring1alpha1.ValiSpec{
		AuthEnabled: false,
		Replicas:    1,
	},
	Image: "ghcr.io/credativ/vali:v2.2.3",
	Resorces: corev1.ResourceRequirements{
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse("1"),
			corev1.ResourceMemory: resource.MustParse("1Gi"),
		},
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse("1"),
			corev1.ResourceMemory: resource.MustParse("1Gi"),
		},
	},
	PVCSize: &defaultValiStorage,
}

type CuratorOptions struct {
	//TODO: implement me
	Image string
}

var defaultCuratorOptions = CuratorOptions{
	Image: "eu.gcr.io/gardener-project/gardener/loki-curator:v0.48.0",
}

// ApplyDefaultSettings manipulates the options to conform to
// build specifications
func ApplyDefaultSettings(opts *Options) error {
	if opts.PriorityClassName == nil {
		opts.PriorityClassName = pointer.String("")
	}

	if err := ApplyDefaultValiSettings(opts); err != nil {
		return err
	}

	if err := ApplyDefaultCuratorSettings(opts); err != nil {
		return err
	}
	//TODO:(vlvasilev) Add here Kube-RBAC-Proxy and Telegraf default settings
	return nil
}

// ApplyDefaultValiSettings manipulates the Vali options to conform to
// build specifications
func ApplyDefaultValiSettings(opts *Options) error {
	if err := mergo.Merge(&opts.Vali, &defaultValiOptions, mergo.WithOverride); err != nil {
		return fmt.Errorf("failed merging Vali user options in namespace %s: %w", opts.Namespace, err)
	}

	return nil
}

// ApplyDefaultCuratorSettings manipulates the Curator options to conform to
// build specifications
func ApplyDefaultCuratorSettings(opts *Options) error {
	if err := mergo.Merge(&opts.Curator, &defaultCuratorOptions, mergo.WithOverride); err != nil {
		return fmt.Errorf("failed merging Curator user options in namespace %s: %w", opts.Namespace, err)
	}

	return nil
}
