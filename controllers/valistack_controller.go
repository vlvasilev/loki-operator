/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1helper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	monitoring1alpha1 "github.com/vlvasilev/loki-operator/api/v1alpha1"
	"github.com/vlvasilev/loki-operator/pkg/operations"
)

// ValiStackReconciler reconciles a ValiStack object
type ValiStackReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=monitoring.gardener.cloud.gardener.cloud,resources=valistacks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.gardener.cloud.gardener.cloud,resources=valistacks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.gardener.cloud.gardener.cloud,resources=valistacks/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=pods;nodes;services;endpoints;configmaps;serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments;statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings;clusterroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;create;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ValiStack object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ValiStackReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	var stack monitoring1alpha1.ValiStack
	if err := r.Get(ctx, req.NamespacedName, &stack); err != nil {
		if apierrors.IsNotFound(err) {
			// maybe the user deleted it before we could react? Either way this isn't an issue
			r.Log.Error(err, "could not find the requested loki stack", "name", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, fmt.Errorf("failed to lookup valistack %s: %w", req.NamespacedName, err)
	}

	operationType := v1beta1helper.ComputeOperationType(stack.ObjectMeta, stack.Status.LastOperation)
	switch {
	case shouldSkipOperation(operationType, &stack):
		return reconcile.Result{}, nil
	case operationType == gardencorev1beta1.LastOperationTypeMigrate:
		return operations.Delete(ctx, &stack, r.Client, r.Log)
	case stack.DeletionTimestamp != nil:
		return operations.Delete(ctx, &stack, r.Client, r.Log)
	case operationType == gardencorev1beta1.LastOperationTypeRestore:
		return operations.Reconcile(ctx, &stack, r.Client, r.Log, operationType)
	default:
		return operations.Reconcile(ctx, &stack, r.Client, r.Log, operationType)
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ValiStackReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoring1alpha1.ValiStack{}).
		Complete(r)
}

// ShouldSkipOperation checks if the current operation should be skipped depending on the lastOperation of the ValiSack object.
func shouldSkipOperation(operationType gardencorev1beta1.LastOperationType, valistack *monitoring1alpha1.ValiStack) bool {
	return operationType != gardencorev1beta1.LastOperationTypeMigrate && operationType != gardencorev1beta1.LastOperationTypeRestore && isMigrated(valistack)
}

// IsMigrated checks if an ValiSack object has been migrated
func isMigrated(valistack *monitoring1alpha1.ValiStack) bool {
	return valistack.Status.LastOperation != nil &&
		valistack.Status.LastOperation.Type == gardencorev1beta1.LastOperationTypeMigrate &&
		valistack.Status.LastOperation.State == gardencorev1beta1.LastOperationStateSucceeded
}
