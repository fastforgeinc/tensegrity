/*
This file is part of the Tensegrity distribution (https://github.com/fastforgeinc/tensegrity)
Copyright (C) 2024 FastForge, Inc.

Tensegrity is free software: you can redistribute it and/or modify it
under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License,
or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with
this program. If not, see http://www.gnu.org/licenses/.
*/

package v1alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reconciler.io/runtime/reconcilers"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"

	k8sv1alpha1 "github.com/fastforgeinc/tensegrity/api/k8s/v1alpha1"
	apiv1alpha1 "github.com/fastforgeinc/tensegrity/api/v1alpha1"
)

func NewStaticReconciler(
	config *reconcilers.Config,
	validationReconciler *ValidationReconciler,
	producerReconciler *ProducerReconciler,
	producerSecretReconciler *ProducerSecretReconciler,
	producerConfigMapReconciler *ProducerConfigMapReconciler) *StaticReconciler {

	return &StaticReconciler{
		Name: "StaticReconciler",
		Setup: func(ctx context.Context, _ ctrl.Manager, builder *builder.Builder) error {
			builder.Watches(new(apiv1alpha1.Static), reconcilers.EnqueueTracked(ctx))
			builder.Watches(new(k8sv1alpha1.DaemonSet), reconcilers.EnqueueTracked(ctx))
			builder.Watches(new(k8sv1alpha1.Deployment), reconcilers.EnqueueTracked(ctx))
			builder.Watches(new(k8sv1alpha1.StatefulSet), reconcilers.EnqueueTracked(ctx))
			return nil
		},
		Config: *config,
		Reconciler: reconcilers.Sequence[*apiv1alpha1.Static]{
			&reconcilers.CastResource[*apiv1alpha1.Static, *apiv1alpha1.Tensegrity]{
				Reconciler: validationReconciler,
			},
			&reconcilers.CastResource[*apiv1alpha1.Static, *apiv1alpha1.Tensegrity]{
				Reconciler: producerReconciler,
			},
			&reconcilers.CastResource[*apiv1alpha1.Static, *metav1.PartialObjectMetadata]{
				Reconciler: producerSecretReconciler,
			},
			&reconcilers.CastResource[*apiv1alpha1.Static, *metav1.PartialObjectMetadata]{
				Reconciler: producerConfigMapReconciler,
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
