package manifests

import (
	gardenerv1beta1const "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	allowValiDes   = "Allows Ingress to the Vali API pods labeled with 'networking.gardener.cloud/to-vali=allowed'."
	allowToValiDes = "Allows Egress from pods labeled with 'networking.gardener.cloud/to-vali=allowed' to the Vali API"
)

var (
	protocolTCP = corev1.ProtocolTCP
)

// TODO: (vlvasilev) Make the NetworkPolicies to be configured solely by the ValiStack specs
func buildNetworkPolicies(opts *Options) []client.Object {
	return []client.Object{
		&networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "allow-vali",
				Namespace: opts.Namespace,
				Annotations: map[string]string{
					gardenerv1beta1const.GardenerDescription: allowToValiDes,
				},
			},
			Spec: networkingv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: commonLabels(),
				},
				Ingress: []networkingv1.NetworkPolicyIngressRule{
					{
						From: []networkingv1.NetworkPolicyPeer{
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"gardener.cloud/role": "logging",
										"app":                 "fluent-bit",
										"role":                "logging",
									},
								},
								NamespaceSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"role": "garden",
									},
								},
							},
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"app":  "aggregate-prometheus",
										"role": "monitoring",
									},
								},
								NamespaceSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"role": "garden",
									},
								},
							},
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"networking.gardener.cloud/to-vali": "allowed",
									},
								},
							},
						},
						Ports: []networkingv1.NetworkPolicyPort{
							{
								Protocol: &protocolTCP,
								Port:     &lokiPortInt,
							},
						},
					},
				},
				PolicyTypes: []networkingv1.PolicyType{
					networkingv1.PolicyTypeIngress,
				},
			},
		},
		&networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "allow-to-vali",
				Namespace: opts.Namespace,
				Annotations: map[string]string{
					gardenerv1beta1const.GardenerDescription: allowValiDes,
				},
			},
			Spec: networkingv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{
						"networking.gardener.cloud/to-vali": "allowed",
					},
				},
				Egress: []networkingv1.NetworkPolicyEgressRule{
					{
						To: []networkingv1.NetworkPolicyPeer{
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: commonLabels(),
								},
							},
						},
						Ports: []networkingv1.NetworkPolicyPort{
							{
								Protocol: &protocolTCP,
								Port:     &lokiPortInt,
							},
						},
					},
				},
				PolicyTypes: []networkingv1.PolicyType{
					networkingv1.PolicyTypeEgress,
				},
			},
		},
	}
}
