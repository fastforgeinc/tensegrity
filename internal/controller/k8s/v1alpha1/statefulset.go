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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"

	"reconciler.io/runtime/reconcilers"

	k8sv1alpha1 "github.com/fastforgeinc/tensegrity/api/k8s/v1alpha1"
	apiv1alpha1 "github.com/fastforgeinc/tensegrity/api/v1alpha1"
	"github.com/fastforgeinc/tensegrity/internal/controller/v1alpha1"
)

func NewStatefulSetReconciler(
	config *reconcilers.Config,
	validationReconciler *v1alpha1.ValidationReconciler,
	consumerReconciler *v1alpha1.ConsumerReconciler,
	consumerSecretReconciler *v1alpha1.ConsumerSecretReconciler,
	consumerConfigMapReconciler *v1alpha1.ConsumerConfigMapReconciler,
	producerReconciler *v1alpha1.ProducerReconciler,
	producerSecretReconciler *v1alpha1.ProducerSecretReconciler,
	producerConfigMapReconciler *v1alpha1.ProducerConfigMapReconciler) *StatefulSetReconciler {

	return &StatefulSetReconciler{
		Name:   "StatefulSetReconciler",
		Config: *config,
		Setup: func(ctx context.Context, _ ctrl.Manager, builder *builder.Builder) error {
			builder.Watches(new(apiv1alpha1.Static), reconcilers.EnqueueTracked(ctx))
			builder.Watches(new(k8sv1alpha1.DaemonSet), reconcilers.EnqueueTracked(ctx))
			builder.Watches(new(k8sv1alpha1.Deployment), reconcilers.EnqueueTracked(ctx))
			builder.Watches(new(k8sv1alpha1.StatefulSet), reconcilers.EnqueueTracked(ctx))
			return nil
		},
		Reconciler: reconcilers.Sequence[*k8sv1alpha1.StatefulSet]{
			&reconcilers.CastResource[*k8sv1alpha1.StatefulSet, *apiv1alpha1.Tensegrity]{
				Reconciler: validationReconciler,
			},
			&reconcilers.CastResource[*k8sv1alpha1.StatefulSet, *apiv1alpha1.Tensegrity]{
				Reconciler: consumerReconciler,
			},
			&reconcilers.CastResource[*k8sv1alpha1.StatefulSet, *metav1.PartialObjectMetadata]{
				Reconciler: consumerSecretReconciler,
			},
			&reconcilers.CastResource[*k8sv1alpha1.StatefulSet, *metav1.PartialObjectMetadata]{
				Reconciler: consumerConfigMapReconciler,
			},
			NewStatefulSetChildReconciler(),
			&reconcilers.CastResource[*k8sv1alpha1.StatefulSet, *apiv1alpha1.Tensegrity]{
				Reconciler: producerReconciler,
			},
			&reconcilers.CastResource[*k8sv1alpha1.StatefulSet, *metav1.PartialObjectMetadata]{
				Reconciler: producerSecretReconciler,
			},
			&reconcilers.CastResource[*k8sv1alpha1.StatefulSet, *metav1.PartialObjectMetadata]{
				Reconciler: producerConfigMapReconciler,
			},
		},
	}
}

// +kubebuilder:rbac:groups=k8s.tensegrity.fastforge.io,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=k8s.tensegrity.fastforge.io,resources=statefulsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=k8s.tensegrity.fastforge.io,resources=statefulsets/finalizers,verbs=update

// StatefulSetReconciler reconciles tensegrity api/k8s/v1alpha1.Deployment resource,
// and runs sequence of other reconcilers to get desired workload.
type StatefulSetReconciler = reconcilers.ResourceReconciler[*k8sv1alpha1.StatefulSet]

func NewStatefulSetChildReconciler() *StatefulSetChildReconciler {
	r := new(StatefulSetChildReconciler)
	r.statefulSetChildReconciler = statefulSetChildReconciler{
		Name:                       "StatefulSetChildReconciler",
		DesiredChild:               r.DesiredChild,
		ReflectChildStatusOnParent: r.ReflectChildStatusOnParent,
		ChildObjectManager: &reconcilers.UpdatingObjectManager[*appsv1.StatefulSet]{
			MergeBeforeUpdate: r.MergeBeforeUpdate,
		},
	}
	return r
}

// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulsets/finalizers,verbs=update

// StatefulSetChildReconciler creates child k8s.io/api/apps/v1.StatefulSet set from workload specs,
// adds ConfigMap and Secret if they are present.
type StatefulSetChildReconciler struct {
	statefulSetChildReconciler
}

func (r *StatefulSetChildReconciler) DesiredChild(
	ctx context.Context, resource *k8sv1alpha1.StatefulSet) (*appsv1.StatefulSet, error) {

	child := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        resource.Name,
			Labels:      resource.Labels,
			Namespace:   resource.Namespace,
			Annotations: make(map[string]string),
		},
		Spec: resource.Spec.StatefulSetSpec,
	}

	var envFrom []corev1.EnvFromSource
	if name := v1alpha1.ConsumerSecretNameFromContext(ctx); len(name) > 0 {
		envFrom = append(envFrom, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: name},
			},
		})

		if key, value := v1alpha1.ConsumerSecretAnnotationFromContext(ctx); len(key) > 0 && len(value) > 0 {
			child.Annotations[key] = value
			if child.Spec.Template.Annotations == nil {
				child.Spec.Template.Annotations = make(map[string]string)
			}
			child.Spec.Template.Annotations[key] = value
		}
	}

	if name := v1alpha1.ConsumerConfigMapNameFromContext(ctx); len(name) > 0 {
		envFrom = append(envFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: name},
			},
		})

		if key, value := v1alpha1.ConsumerConfigMapAnnotationFromContext(ctx); len(key) > 0 && len(value) > 0 {
			child.Annotations[key] = value
			if child.Spec.Template.Annotations == nil {
				child.Spec.Template.Annotations = make(map[string]string)
			}
			child.Spec.Template.Annotations[key] = value
		}
	}

	if len(envFrom) > 0 {
		for i, container := range child.Spec.Template.Spec.InitContainers {
			child.Spec.Template.Spec.InitContainers[i].EnvFrom = append(
				container.EnvFrom, envFrom...)
		}
		for i, container := range child.Spec.Template.Spec.Containers {
			child.Spec.Template.Spec.Containers[i].EnvFrom = append(
				container.EnvFrom, envFrom...)
		}
	}

	return child, nil
}

func (r *StatefulSetChildReconciler) MergeBeforeUpdate(current, desired *appsv1.StatefulSet) {
	current.Annotations = reconcilers.MergeMaps(current.Annotations, desired.Annotations)
	current.Labels = desired.Labels
}

func (r *StatefulSetChildReconciler) ReflectChildStatusOnParent(
	_ context.Context, _ *k8sv1alpha1.StatefulSet, _ *appsv1.StatefulSet, _ error) {
}

type statefulSetChildReconciler = reconcilers.ChildReconciler[
	*k8sv1alpha1.StatefulSet,
	*appsv1.StatefulSet,
	*appsv1.StatefulSetList,
]
