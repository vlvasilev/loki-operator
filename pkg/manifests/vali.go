package manifests

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

var (
	volumeFileSystemMode = corev1.PersistentVolumeFilesystem
)

func getValiContainer(opt *ValiOptions) corev1.Container {
	return corev1.Container{
		Name:  "vali",
		Image: opt.Image,
		Args: []string{
			"-config.file=/etc/vali/vali.yaml",
			//TODO: (vlvasilev) add extra args
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "config",
				MountPath: "/etc/vali",
			},
			{
				Name:      "data",
				MountPath: "/data",
			},
		},
		Ports: []corev1.ContainerPort{
			{
				Name:          "metrics",
				ContainerPort: 3100,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		LivenessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/ready",
					Port:   intstr.FromString("metrics"),
					Scheme: corev1.URISchemeHTTP,
				},
			},
			InitialDelaySeconds: 120,
			FailureThreshold:    5,
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/ready",
					Port:   intstr.FromString("metrics"),
					Scheme: corev1.URISchemeHTTP,
				},
			},
			InitialDelaySeconds: 80,
			FailureThreshold:    7,
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("200m"),
				corev1.ResourceMemory: resource.MustParse("300Mi"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("3Gi"),
			},
		},
		SecurityContext: &corev1.SecurityContext{
			RunAsGroup:             pointer.Int64(10001),
			RunAsUser:              pointer.Int64(10001),
			RunAsNonRoot:           pointer.Bool(true),
			ReadOnlyRootFilesystem: pointer.Bool(true),
		},
	}
}

func getValiVolumes(opt *Options) []corev1.Volume {
	return []corev1.Volume{
		{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: lokiConfigMapName(opt.ConfigSHA),
					},
				},
			},
		},
	}
}

func getValiVolumeClaimTemplate(opts *ValiOptions) corev1.PersistentVolumeClaim {
	return corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{

			Name: "data",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: *opts.PVCSize,
				},
			},
			VolumeMode: &volumeFileSystemMode,
		},
	}
}

func getEmptyValiStatefulSet(opt *Options) (*appsv1.StatefulSet, error) {
	if opt.PriorityClassName == nil {
		return nil, fmt.Errorf("the PriorityClassName name is nill. Please, call ApplyDefaultSettings before getEmptyValiStatefulSet")
	}
	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "vali",
			Namespace: opt.Namespace,
			Labels:    commonLabels(),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &opt.Vali.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: commonLabels(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: commonLabels(),
					// TODO: (vlvasilev) move this login in a build/compile function
					// because we need the monitor the telegraf config map as well.
					Annotations: map[string]string{
						"reference.resources.gardener.cloud/configmap-vali": lokiConfigMapName(opt.ConfigSHA),
					},
				},
				Spec: corev1.PodSpec{
					AutomountServiceAccountToken: pointer.Bool(false),
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup:      pointer.Int64(10001),
						RunAsGroup:   pointer.Int64(10001),
						RunAsUser:    pointer.Int64(10001),
						RunAsNonRoot: pointer.Bool(true),
					},
					PriorityClassName: *opt.PriorityClassName,
				},
			},
		},
	}, nil
}

func getCuratorContainer(opt *CuratorOptions) corev1.Container {
	return corev1.Container{
		Name:  "curator",
		Image: opt.Image,
		Args: []string{
			"-config=/etc/vali/curator.yaml",
		},
		Ports: []corev1.ContainerPort{
			{
				Name:          "metrics",
				ContainerPort: 2718,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("10m"),
				corev1.ResourceMemory: resource.MustParse("12Mi"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("700Mi"),
			},
		},
		SecurityContext: &corev1.SecurityContext{
			RunAsGroup:             pointer.Int64(10001),
			RunAsUser:              pointer.Int64(10001),
			RunAsNonRoot:           pointer.Bool(true),
			ReadOnlyRootFilesystem: pointer.Bool(true),
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "config",
				MountPath: "/etc/vali",
			},
			{
				Name:      "data",
				MountPath: "/data",
			},
		},
	}
}
