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
	k8sv1alpha1 "github.com/fastforgeinc/tensegrity/api/k8s/v1alpha1"
	apiv1alpha1 "github.com/fastforgeinc/tensegrity/api/v1alpha1"
	"github.com/fastforgeinc/tensegrity/internal/controller/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reconciler.io/runtime/reconcilers"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

func NewDaemonSetReconciler(config *reconcilers.Config) *DaemonSetReconciler {
	return &DaemonSetReconciler{
		Name: "DaemonSetReconciler",
		Setup: func(ctx context.Context, _ ctrl.Manager, builder *builder.Builder) error {
			builder.Watches(new(k8sv1alpha1.DaemonSet), reconcilers.EnqueueTracked(ctx))
			return nil
		},
		Config: *config,
		Reconciler: reconcilers.Sequence[*k8sv1alpha1.DaemonSet]{
			&reconcilers.CastResource[*k8sv1alpha1.DaemonSet, *apiv1alpha1.Tensegrity]{
				Reconciler: v1alpha1.NewConsumerReconciler(),
			},
			&reconcilers.CastResource[*k8sv1alpha1.DaemonSet, *metav1.PartialObjectMetadata]{
				Reconciler: v1alpha1.NewConfigMapReconciler(),
			},
			&reconcilers.CastResource[*k8sv1alpha1.DaemonSet, *metav1.PartialObjectMetadata]{
				Reconciler: v1alpha1.NewSecretReconciler(),
			},
			NewDaemonSetChildReconciler(),
			&reconcilers.CastResource[*k8sv1alpha1.DaemonSet, *apiv1alpha1.Tensegrity]{
				Reconciler: v1alpha1.NewProducerReconciler(),
			},
		},
	}
}

// DaemonSetReconciler reconciles tensegrity v1alpha1.DaemonSet resource,
// and runs sequence of other reconcilers to get desired workload.
// +kubebuilder:rbac:groups=tensegrity.fastforge.io,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tensegrity.fastforge.io,resources=daemonsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=tensegrity.fastforge.io,resources=daemonsets/finalizers,verbs=update
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

// DaemonSetChildReconciler creates child daemonSet from workload specs,
// add ConfigMap and Secret if they are present.
// +kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=daemonsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=daemonsets/finalizers,verbs=update
type DaemonSetChildReconciler struct {
	daemonSetChildReconciler
}

func (r *DaemonSetChildReconciler) DesiredChild(
	ctx context.Context, resource *k8sv1alpha1.DaemonSet) (*appsv1.DaemonSet, error) {

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
		for i, container := range resource.Spec.DaemonSetSpec.Template.Spec.InitContainers {
			resource.Spec.DaemonSetSpec.Template.Spec.InitContainers[i].EnvFrom = append(
				container.EnvFrom, envFrom...)
		}
		for i, container := range resource.Spec.DaemonSetSpec.Template.Spec.Containers {
			resource.Spec.DaemonSetSpec.Template.Spec.Containers[i].EnvFrom = append(
				container.EnvFrom, envFrom...)
		}
	}

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        resource.Name,
			Labels:      resource.Labels,
			Namespace:   resource.Namespace,
			Annotations: make(map[string]string),
		},
		Spec: resource.Spec.DaemonSetSpec,
	}, nil
}

func (r *DaemonSetChildReconciler) MergeBeforeUpdate(current, desired *appsv1.DaemonSet) {
	current.Labels = desired.Labels
	current.Spec = desired.Spec
}

func (r *DaemonSetChildReconciler) ReflectChildStatusOnParent(
	ctx context.Context, parent *k8sv1alpha1.DaemonSet, child *appsv1.DaemonSet, err error) {

	// TODO: add status of configuration
	return
}

type daemonSetChildReconciler = reconcilers.ChildReconciler[
	*k8sv1alpha1.DaemonSet,
	*appsv1.DaemonSet,
	*appsv1.DaemonSetList,
]