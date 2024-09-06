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

func NewDaemonSetReconciler(
	config *reconcilers.Config,
	consumerReconciler *v1alpha1.ConsumerReconciler,
	consumerSecretReconciler *v1alpha1.ConsumerSecretReconciler,
	consumerConfigMapReconciler *v1alpha1.ConsumerConfigMapReconciler,
	producerReconciler *v1alpha1.ProducerReconciler,
	producerSecretReconciler *v1alpha1.ProducerSecretReconciler,
	producerConfigMapReconciler *v1alpha1.ProducerConfigMapReconciler) *DaemonSetReconciler {

	return &DaemonSetReconciler{
		Name: "DaemonSetReconciler",
		Setup: func(ctx context.Context, _ ctrl.Manager, builder *builder.Builder) error {
			builder.Watches(new(k8sv1alpha1.DaemonSet), reconcilers.EnqueueTracked(ctx))
			return nil
		},
		Config: *config,
		Reconciler: reconcilers.Sequence[*k8sv1alpha1.DaemonSet]{
			&reconcilers.CastResource[*k8sv1alpha1.DaemonSet, *apiv1alpha1.Tensegrity]{
				Reconciler: consumerReconciler,
			},
			&reconcilers.CastResource[*k8sv1alpha1.DaemonSet, *metav1.PartialObjectMetadata]{
				Reconciler: consumerSecretReconciler,
			},
			&reconcilers.CastResource[*k8sv1alpha1.DaemonSet, *metav1.PartialObjectMetadata]{
				Reconciler: consumerConfigMapReconciler,
			},
			NewDaemonSetChildReconciler(),
			&reconcilers.CastResource[*k8sv1alpha1.DaemonSet, *apiv1alpha1.Tensegrity]{
				Reconciler: producerReconciler,
			},
			&reconcilers.CastResource[*k8sv1alpha1.DaemonSet, *metav1.PartialObjectMetadata]{
				Reconciler: producerSecretReconciler,
			},
			&reconcilers.CastResource[*k8sv1alpha1.DaemonSet, *metav1.PartialObjectMetadata]{
				Reconciler: producerConfigMapReconciler,
			},
		},
	}
}

// +kubebuilder:rbac:groups=k8s.tensegrity.fastforge.io,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=k8s.tensegrity.fastforge.io,resources=daemonsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=k8s.tensegrity.fastforge.io,resources=daemonsets/finalizers,verbs=update

// DaemonSetReconciler reconciles tensegrity api/k8s/v1alpha1.DaemonSet resource,
// and runs sequence of other reconcilers to get desired workload.
type DaemonSetReconciler = reconcilers.ResourceReconciler[*k8sv1alpha1.DaemonSet]

func NewDaemonSetChildReconciler() *DaemonSetChildReconciler {
	r := new(DaemonSetChildReconciler)
	r.daemonSetChildReconciler = daemonSetChildReconciler{
		Name:                       "DaemonSetChildReconciler",
		DesiredChild:               r.DesiredChild,
		MergeBeforeUpdate:          r.MergeBeforeUpdate,
		ReflectChildStatusOnParent: r.ReflectChildStatusOnParent,
	}
	return r
}

// +kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=daemonsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=daemonsets/finalizers,verbs=update

// DaemonSetChildReconciler creates child k8s.io/api/apps/v1.DaemonSet from workload specs,
// adds ConfigMap and Secret if they are present.
type DaemonSetChildReconciler struct {
	daemonSetChildReconciler
}

func (r *DaemonSetChildReconciler) DesiredChild(
	ctx context.Context, resource *k8sv1alpha1.DaemonSet) (*appsv1.DaemonSet, error) {

	child := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        resource.Name,
			Labels:      resource.Labels,
			Namespace:   resource.Namespace,
			Annotations: make(map[string]string),
		},
		Spec: resource.Spec.DaemonSetSpec,
	}

	var envFrom []corev1.EnvFromSource
	if name := v1alpha1.ConsumerSecretNameFromContext(ctx); len(name) > 0 {
		envFrom = append(envFrom, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: name},
			},
		})
	}

	if name := v1alpha1.ConsumerConfigMapNameFromContext(ctx); len(name) > 0 {
		envFrom = append(envFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: name},
			},
		})
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

func (r *DaemonSetChildReconciler) MergeBeforeUpdate(current, desired *appsv1.DaemonSet) {
	current.Labels = desired.Labels
	current.Spec = desired.Spec
}

func (r *DaemonSetChildReconciler) ReflectChildStatusOnParent(
	_ context.Context, _ *k8sv1alpha1.DaemonSet, _ *appsv1.DaemonSet, _ error) {

	return
}

type daemonSetChildReconciler = reconcilers.ChildReconciler[
	*k8sv1alpha1.DaemonSet,
	*appsv1.DaemonSet,
	*appsv1.DaemonSetList,
]
