package operations

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	monitoring1alpha1 "github.com/vlvasilev/loki-operator/api/v1alpha1"
	manifests "github.com/vlvasilev/loki-operator/pkg/manifests"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Reconcile(ctx context.Context, valistack monitoring1alpha1.ValiStack, k client.Client, s *runtime.Scheme, logger logr.Logger) error {
	ll := logger.WithValues("valistack", types.NamespacedName{Namespace: valistack.Namespace, Name: valistack.Name}, "event", "reconcile")

	// Here we will translate the lokiv1beta1.LokiStack options into manifest options
	opts := manifests.Options{
		Name:      valistack.Name,
		Namespace: valistack.Namespace,
		Vali: manifests.ValiOptions{
			Spec: valistack.Spec.Vali,
		},
		PriorityClassName: valistack.Spec.PriorityClassName,
	}
	ll.Info("begin building manifests")

	if optErr := manifests.ApplyDefaultSettings(&opts); optErr != nil {
		ll.Error(optErr, "failed to conform options to build settings")
		return optErr
	}

	objects, err := manifests.BuildAll(&opts)
	if err != nil {
		ll.Error(err, "failed to build manifests")
		return err
	}
	ll.Info("manifests built", "count", len(objects))

	var errCount int32

	for _, obj := range objects {
		l := ll.WithValues("object_name", obj.GetName(), "object_kind", obj.GetObjectKind())

		desired := obj.DeepCopyObject().(client.Object)
		mutateFn := manifests.MutateFuncFor(obj, desired)

		op, err := ctrl.CreateOrUpdate(ctx, k, obj, mutateFn)
		if err != nil {
			l.Error(err, "failed to configure resource")
			errCount++
			continue
		}

		l.Info(fmt.Sprintf("Resource has been %s", op))
	}

	if errCount > 0 {
		return fmt.Errorf("failed to configure valistack resources %s", valistack.Namespace)
	}

	return nil
}
