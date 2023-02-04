package manifests

import (
	hvpav1alpha1 "github.com/gardener/hvpa-controller/api/v1alpha1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpa_api "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"
	"k8s.io/utils/pointer"
)

var (
	updateModeAuto                        = hvpav1alpha1.UpdateModeAuto
	containerControlledValuesRequestsOnly = vpa_api.ContainerControlledValuesRequestsOnly
	containerScalingModeOff               = vpa_api.ContainerScalingModeOff
)

func buildHVPA(opts *Options) *hvpav1alpha1.Hvpa {
	hvpa := &hvpav1alpha1.Hvpa{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "vali",
			Namespace: opts.Namespace,
			Labels:    commonLabels(),
		},
		Spec: hvpav1alpha1.HvpaSpec{
			Replicas:              &opts.Vali.Spec.Replicas,
			MaintenanceTimeWindow: &opts.HVPA.MaintenanceTimeWindow,
			Hpa: hvpav1alpha1.HpaSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"role": "vali-hpa",
					},
				},
				Deploy: false,
				Template: hvpav1alpha1.HpaTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"role": "vali-hpa",
						},
					},
					Spec: hvpav1alpha1.HpaTemplateSpec{
						MaxReplicas: opts.Vali.Spec.Replicas,
						MinReplicas: &opts.Vali.Spec.Replicas,
						Metrics: []autoscalingv2beta1.MetricSpec{
							{
								Type: autoscalingv2beta1.ResourceMetricSourceType,
								Resource: &autoscalingv2beta1.ResourceMetricSource{
									Name:                     corev1.ResourceCPU,
									TargetAverageUtilization: pointer.Int32(80),
								},
							},
							{
								Type: autoscalingv2beta1.ResourceMetricSourceType,
								Resource: &autoscalingv2beta1.ResourceMetricSource{
									Name:                     corev1.ResourceMemory,
									TargetAverageUtilization: pointer.Int32(80),
								},
							},
						},
					},
				},
			},
			Vpa: hvpav1alpha1.VpaSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"role": "vali-hpa",
					},
				},
				Deploy: true,
				ScaleUp: hvpav1alpha1.ScaleType{
					UpdatePolicy: hvpav1alpha1.UpdatePolicy{
						UpdateMode: &updateModeAuto,
					},
					StabilizationDuration: pointer.String("5m"),
					MinChange: hvpav1alpha1.ScaleParams{
						CPU: hvpav1alpha1.ChangeParams{
							Value:      pointer.String("100m"),
							Percentage: pointer.Int32(80),
						},
						Memory: hvpav1alpha1.ChangeParams{
							Value:      pointer.String("300M"),
							Percentage: pointer.Int32(80),
						},
					},
				},
				ScaleDown: hvpav1alpha1.ScaleType{
					UpdatePolicy: hvpav1alpha1.UpdatePolicy{
						UpdateMode: &updateModeAuto,
					},
					StabilizationDuration: pointer.String("168h"),
					MinChange: hvpav1alpha1.ScaleParams{
						CPU: hvpav1alpha1.ChangeParams{
							Value:      pointer.String("200m"),
							Percentage: pointer.Int32(80),
						},
						Memory: hvpav1alpha1.ChangeParams{
							Value:      pointer.String("500M"),
							Percentage: pointer.Int32(80),
						},
					},
				},
				LimitsRequestsGapScaleParams: hvpav1alpha1.ScaleParams{
					CPU: hvpav1alpha1.ChangeParams{
						Value:      pointer.String("300m"),
						Percentage: pointer.Int32(40),
					},
					Memory: hvpav1alpha1.ChangeParams{
						Value:      pointer.String("1000M"),
						Percentage: pointer.Int32(40),
					},
				},
				Template: hvpav1alpha1.VpaTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"role": "vali-vpa",
						},
					},
					Spec: hvpav1alpha1.VpaTemplateSpec{
						ResourcePolicy: &vpa_api.PodResourcePolicy{
							ContainerPolicies: []vpa_api.ContainerResourcePolicy{
								{
									ContainerName:    "vali",
									ControlledValues: &containerControlledValuesRequestsOnly,
									MinAllowed: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse("200m"),
										corev1.ResourceMemory: resource.MustParse("300M"),
									},
									MaxAllowed: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse("800m"),
										corev1.ResourceMemory: resource.MustParse("3Gi"),
									},
								},
								{
									ContainerName: "curator",
									Mode:          &containerScalingModeOff,
								},
							},
						},
					},
				},
			},
			WeightBasedScalingIntervals: []hvpav1alpha1.WeightBasedScalingInterval{
				{
					VpaWeight:         hvpav1alpha1.VpaOnly,
					StartReplicaCount: opts.Vali.Spec.Replicas,
					LastReplicaCount:  opts.Vali.Spec.Replicas,
				},
			},
			TargetRef: &autoscalingv2beta1.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "something",
				Name:       "something",
			},
		},
	}

	if opts.HVPA.MaintenanceTimeWindow.Begin != "" {
		hvpa.Spec.Vpa.ScaleDown.UpdatePolicy.UpdateMode = pointer.String("MaintenanceWindow")
	}
	return hvpa
}
