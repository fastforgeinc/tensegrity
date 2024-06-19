package v1alpha1

import (
	"github.com/fastforgeinc/tensegrity/api/v1alpha1"
	"reconciler.io/runtime/reconcilers"
)

type workloadReconciler = reconcilers.SyncReconciler[*v1alpha1.Tensegrity]
