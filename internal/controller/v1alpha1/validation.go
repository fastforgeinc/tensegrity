package v1alpha1

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/fastforgeinc/tensegrity/api/v1alpha1"
)

func NewValidationReconciler() *ValidationReconciler {
	r := new(ValidationReconciler)
	r.workloadReconciler = workloadReconciler{
		Name: "ValidationReconciler",
		Sync: r.Sync,
	}
	return r
}

type ValidationReconciler struct {
	workloadReconciler
}

func (r *ValidationReconciler) Sync(_ context.Context, resource *v1alpha1.Tensegrity) error {
	if errs := resource.Spec.Validate(); len(errs) != 0 {
		aggrErr := errs.ToAggregate()
		message := fmt.Sprintf(v1alpha1.SpecInvalidMessage, aggrErr.Error())
		condition := v1alpha1.NewTensegrityCondition(
			v1alpha1.TensegrityInvalid, corev1.ConditionTrue,
			v1alpha1.SpecInvalidReason, message)
		v1alpha1.SetTensegrityCondition(&resource.Status, *condition)
		return aggrErr
	}
	return nil
}
