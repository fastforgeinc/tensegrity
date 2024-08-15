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

func NewStatefulSetReconciler(config *reconcilers.Config) *StatefulSetReconciler {
	return &StatefulSetReconciler{
		Name:   "StatefulSetReconciler",
		Config: *config,
		Setup: func(ctx context.Context, _ ctrl.Manager, builder *builder.Builder) error {
			builder.Watches(new(k8sv1alpha1.StatefulSet), reconcilers.EnqueueTracked(ctx))
			return nil
		},
		Reconciler: reconcilers.Sequence[*k8sv1alpha1.StatefulSet]{
			&reconcilers.CastResource[*k8sv1alpha1.StatefulSet, *apiv1alpha1.Tensegrity]{
				Reconciler: v1alpha1.NewConsumerReconciler(),
			},
			&reconcilers.CastResource[*k8sv1alpha1.StatefulSet, *metav1.PartialObjectMetadata]{
				Reconciler: v1alpha1.NewConfigMapReconciler(),
			},
			&reconcilers.CastResource[*k8sv1alpha1.StatefulSet, *metav1.PartialObjectMetadata]{
				Reconciler: v1alpha1.NewSecretReconciler(),
			},
			NewStatefulSetChildReconciler(),
			&reconcilers.CastResource[*k8sv1alpha1.StatefulSet, *apiv1alpha1.Tensegrity]{
				Reconciler: v1alpha1.NewProducerReconciler(),
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
		MergeBeforeUpdate:          r.MergeBeforeUpdate,
		ReflectChildStatusOnParent: r.ReflectChildStatusOnParent,
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
	if secret := v1alpha1.SecretFromContext(ctx); secret != nil {
		envFrom = append(envFrom, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: secret.Name},
			},
		})
	}

	if configMap := v1alpha1.ConfigMapFromContext(ctx); configMap != nil {
		envFrom = append(envFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: configMap.Name},
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

func (r *StatefulSetChildReconciler) MergeBeforeUpdate(current, desired *appsv1.StatefulSet) {
	current.Labels = desired.Labels
	current.Spec = desired.Spec
}

func (r *StatefulSetChildReconciler) ReflectChildStatusOnParent(
	_ context.Context, _ *k8sv1alpha1.StatefulSet, _ *appsv1.StatefulSet, _ error) {

	return
}

type statefulSetChildReconciler = reconcilers.ChildReconciler[
	*k8sv1alpha1.StatefulSet,
	*appsv1.StatefulSet,
	*appsv1.StatefulSetList,
]
