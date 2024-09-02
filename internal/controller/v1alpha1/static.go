/*
Copyright 2024 FastForge Inc. support@fastforge.com.

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

package v1alpha1

import (
	"context"
	"reconciler.io/runtime/reconcilers"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"

	apiv1alpha1 "github.com/fastforgeinc/tensegrity/api/v1alpha1"
)

func NewStaticReconciler(config *reconcilers.Config) *StaticReconciler {
	return &StaticReconciler{
		Name: "StaticReconciler",
		Setup: func(ctx context.Context, _ ctrl.Manager, builder *builder.Builder) error {
			builder.Watches(new(apiv1alpha1.Static), reconcilers.EnqueueTracked(ctx))
			return nil
		},
		Config: *config,
		Reconciler: reconcilers.Sequence[*apiv1alpha1.Static]{
			&reconcilers.CastResource[*apiv1alpha1.Static, *apiv1alpha1.Tensegrity]{
				Reconciler: NewProducerReconciler(),
			},
		},
	}
}

//+kubebuilder:rbac:groups=tensegrity.fastforge.io,resources=statics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tensegrity.fastforge.io,resources=statics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tensegrity.fastforge.io,resources=statics/finalizers,verbs=update

// StaticReconciler reconciles tensegrity api/v1alpha1.Static resource,
// and runs sequence of other reconcilers to get desired state.
type StaticReconciler = reconcilers.ResourceReconciler[*apiv1alpha1.Static]
