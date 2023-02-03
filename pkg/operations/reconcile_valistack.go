package operations

import (
	"context"
	"fmt"
	"time"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/controllerutils"
	"github.com/go-logr/logr"

	monitoring1alpha1 "github.com/vlvasilev/loki-operator/api/v1alpha1"
	manifests "github.com/vlvasilev/loki-operator/pkg/manifests"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// FinalizerName is the worker controller finalizer.
	FinalizerName = "extensions.gardener.cloud/worker"
)

func Reconcile(ctx context.Context, valistack *monitoring1alpha1.ValiStack, k client.Client, logger logr.Logger, operationType v1beta1.LastOperationType) (reconcile.Result, error) {
	ll := logger.WithValues("valistack", types.NamespacedName{Namespace: valistack.Namespace, Name: valistack.Name}, "event", "reconcile")

	if !controllerutil.ContainsFinalizer(valistack, FinalizerName) {
		ll.Info("Adding finalizer")
		if err := controllerutils.AddFinalizers(ctx, k, valistack, FinalizerName); err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
	}

	if err := Processing(ctx, k, ll, valistack, operationType, "Reconciling the ValiStack"); err != nil {
		return reconcile.Result{}, err
	}
	ll.Info("Starting the reconciliation of valistack")

	// Here we will translate the valiv1alpha1.ValiStack options into manifest options
	opts := manifests.Options{
		Name:      valistack.Name,
		Namespace: valistack.Namespace,
		Vali: manifests.ValiOptions{
			Spec: valistack.Spec.Vali,
		},
		PriorityClassName: valistack.Spec.PriorityClassName,
	}
	ll.Info("begin building manifests")

	if err := manifests.ApplyDefaultSettings(&opts); err != nil {
		ll.Error(err, "failed to conform options to build settings")
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, err
	}

	objects, err := manifests.BuildAll(&opts)
	if err != nil {
		ll.Error(err, "failed to build manifests")
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, err
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
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, fmt.Errorf("failed to configure valistack resources %s", valistack.Namespace)
	}

	if err := Success(ctx, k, logger, valistack, operationType, "Successfully reconciled ValiStack"); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
