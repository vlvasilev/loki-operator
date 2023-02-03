package operations

import (
	"context"
	"fmt"

	gardenextension "github.com/gardener/gardener/extensions/pkg/controller"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1helper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
	"github.com/go-logr/logr"
	monitoring1alpha1 "github.com/vlvasilev/loki-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Processing(
	ctx context.Context,
	c client.Client,
	log logr.Logger,
	valistack *monitoring1alpha1.ValiStack,
	lastOperationType gardencorev1beta1.LastOperationType,
	description string,
) error {
	if c == nil {
		return fmt.Errorf("client is not set. Call Processing() with valid client")
	}

	log.Info(description) //nolint:logcheck

	patch := client.MergeFrom(valistack.DeepCopyObject().(client.Object))
	lastOp := gardenextension.LastOperation(lastOperationType, gardencorev1beta1.LastOperationStateProcessing, 1, description)
	valistack.Status.LastOperation = lastOp

	return c.Status().Patch(ctx, valistack, patch)
}

func Error(
	ctx context.Context,
	c client.Client,
	log logr.Logger,
	valistack *monitoring1alpha1.ValiStack,
	err error,
	lastOperationType gardencorev1beta1.LastOperationType,
	description string,
) error {
	if c == nil {
		return fmt.Errorf("client is not set. Call Error() with valid client")
	}

	var (
		errDescription  = v1beta1helper.FormatLastErrDescription(fmt.Errorf("%s: %v", description, err))
		lastOp, lastErr = gardenextension.ReconcileError(lastOperationType, errDescription, 50, v1beta1helper.ExtractErrorCodes(err)...)
	)

	log.Error(fmt.Errorf(errDescription), "Error") //nolint:logcheck

	patch := client.MergeFrom(valistack.DeepCopyObject().(client.Object))
	valistack.Status.ObservedGeneration = valistack.GetGeneration()
	valistack.Status.LastOperation = lastOp
	valistack.Status.LastError = lastErr

	return c.Status().Patch(ctx, valistack, patch)
}

func Success(
	ctx context.Context,
	c client.Client,
	log logr.Logger,
	valistack *monitoring1alpha1.ValiStack,
	lastOperationType gardencorev1beta1.LastOperationType,
	description string,
) error {
	if c == nil {
		return fmt.Errorf("client is not set. Call InjectClient() first")
	}

	log.Info(description) //nolint:logcheck

	patch := client.MergeFrom(valistack.DeepCopyObject().(client.Object))
	lastOp, lastErr := gardenextension.ReconcileSucceeded(lastOperationType, description)
	valistack.Status.ObservedGeneration = valistack.GetGeneration()
	valistack.Status.LastOperation = lastOp
	valistack.Status.LastError = lastErr

	return c.Status().Patch(ctx, valistack, patch)
}
