package manifests

import (
	hvpav1alpha1 "github.com/gardener/hvpa-controller/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Delete return all manifests required to to delete full Vali Stack
func Delete(opts *Options) ([]client.Object, error) {
	return []client.Object{
		&networkingv1.NetworkPolicy{ObjectMeta: metav1.ObjectMeta{Name: "allow-vali", Namespace: opts.Namespace}},
		&networkingv1.NetworkPolicy{ObjectMeta: metav1.ObjectMeta{Name: "allow-to-vali", Namespace: opts.Namespace}},
		&hvpav1alpha1.Hvpa{ObjectMeta: metav1.ObjectMeta{Name: "loki", Namespace: opts.Namespace}},
		&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "loki-config", Namespace: opts.Namespace}},
		&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "loki", Namespace: opts.Namespace}},
		&appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: "loki", Namespace: opts.Namespace}},
		&corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: "loki-loki-0", Namespace: opts.Namespace}},
	}, nil

	//return kubernetesutils.DeleteObjects(ctx, k8sClient, resources...)
}
