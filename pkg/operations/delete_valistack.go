package operations

import (
	"context"
	"fmt"
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/controllerutils"
	kubernetesutils "github.com/gardener/gardener/pkg/utils/kubernetes"
	"github.com/go-logr/logr"
	monitoring1alpha1 "github.com/vlvasilev/loki-operator/api/v1alpha1"
	manifests "github.com/vlvasilev/loki-operator/pkg/manifests"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Delete(ctx context.Context, valistack *monitoring1alpha1.ValiStack, c client.Client, logger logr.Logger) (reconcile.Result, error) {
	ll := logger.WithValues("valistack", types.NamespacedName{Namespace: valistack.Namespace, Name: valistack.Name}, "event", "delete")

	if !controllerutil.ContainsFinalizer(valistack, FinalizerName) {
		ll.Info("Deleting ValiStack causes a no-op as there is no finalizer")
		return reconcile.Result{}, nil
	}

	if err := Processing(ctx, c, ll, valistack, gardencorev1beta1.LastOperationTypeDelete, "Deleting the ValiStack"); err != nil {
		return reconcile.Result{}, err
	}
	ll.Info("Starting the deletion of valistack")

	// Here we will translate the valiv1alpha1.ValiStack options into manifest options
	// A lot this options are redundant
	opts := manifests.Options{
		Name:      valistack.Name,
		Namespace: valistack.Namespace,
		Vali: manifests.ValiOptions{
			Spec: valistack.Spec.Vali,
		},
		HVPA:              valistack.Spec.HVPA.DeepCopy(),
		PriorityClassName: valistack.Spec.PriorityClassName,
	}
	ll.Info("begin deleting manifests")

	objects, err := manifests.Delete(&opts)
	if err != nil {
		ll.Error(err, "failed to delete manifests")
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, err
	}
	ll.Info("manifests to delete", "count", len(objects))

	if err := kubernetesutils.DeleteObjects(ctx, c, objects...); err != nil {
		ll.Error(err, "failed to delete manifests")
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, err
	}

	if err := Success(ctx, c, logger, valistack, gardencorev1beta1.LastOperationTypeDelete, "Successfully deleted ValiStack"); err != nil {
		return reconcile.Result{}, err
	}

	if controllerutil.ContainsFinalizer(valistack, FinalizerName) {
		ll.Info("Removing finalizer")
		if err := controllerutils.RemoveFinalizers(ctx, c, valistack, FinalizerName); err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
		}
	}

	return reconcile.Result{}, nil
}
