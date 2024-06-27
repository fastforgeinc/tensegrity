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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"

	"reconciler.io/runtime/reconcilers"

	rolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"

	argov1alpha1 "github.com/fastforgeinc/tensegrity/api/argo/v1alpha1"
	apiv1alpha1 "github.com/fastforgeinc/tensegrity/api/v1alpha1"
	"github.com/fastforgeinc/tensegrity/internal/controller/v1alpha1"
)

func NewRolloutReconciler(config *reconcilers.Config) *RolloutReconciler {
	return &RolloutReconciler{
		Name: "RolloutReconciler",
		Setup: func(ctx context.Context, _ ctrl.Manager, builder *builder.Builder) error {
			builder.Watches(new(argov1alpha1.Rollout), reconcilers.EnqueueTracked(ctx))
			return nil
		},
		Config: *config,
		Reconciler: reconcilers.Sequence[*argov1alpha1.Rollout]{
			&reconcilers.CastResource[*argov1alpha1.Rollout, *apiv1alpha1.Tensegrity]{
				Reconciler: v1alpha1.NewConsumerReconciler(),
			},
			&reconcilers.CastResource[*argov1alpha1.Rollout, *metav1.PartialObjectMetadata]{
				Reconciler: v1alpha1.NewConfigMapReconciler(),
			},
			&reconcilers.CastResource[*argov1alpha1.Rollout, *metav1.PartialObjectMetadata]{
				Reconciler: v1alpha1.NewSecretReconciler(),
			},
			NewRolloutChildReconciler(),
			&reconcilers.CastResource[*argov1alpha1.Rollout, *apiv1alpha1.Tensegrity]{
				Reconciler: v1alpha1.NewProducerReconciler(),
			},
		},
	}
}

// +kubebuilder:rbac:groups=argo.tensegrity.fastforge.io,resources=rollouts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=argo.tensegrity.fastforge.io,resources=rollouts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=argo.tensegrity.fastforge.io,resources=rollouts/finalizers,verbs=update

// RolloutReconciler reconciles tensegrity api/argo/v1alpha1.Deployment resource,
// and runs sequence of other reconcilers to get desired workload.
type RolloutReconciler = reconcilers.ResourceReconciler[*argov1alpha1.Rollout]

func NewRolloutChildReconciler() *RolloutChildReconciler {
	r := new(RolloutChildReconciler)
	r.rolloutChildReconciler = rolloutChildReconciler{
		Name:                       "RolloutChildReconciler",
		DesiredChild:               r.DesiredChild,
		MergeBeforeUpdate:          r.MergeBeforeUpdate,
		ReflectChildStatusOnParent: r.ReflectChildStatusOnParent,
	}
	return r
}

// +kubebuilder:rbac:groups=argoproj.io,resources=rollouts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=argoproj.io,resources=rollouts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=argoproj.io,resources=rollouts/finalizers,verbs=update

// RolloutChildReconciler creates child github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1.Rollout
// from workload specs, adds ConfigMap and Secret if they are present.
type RolloutChildReconciler struct {
	rolloutChildReconciler
}

func (r *RolloutChildReconciler) DesiredChild(
	ctx context.Context, resource *argov1alpha1.Rollout) (*rolloutsv1alpha1.Rollout, error) {

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
		for i, container := range resource.Spec.RolloutSpec.Template.Spec.InitContainers {
			resource.Spec.RolloutSpec.Template.Spec.InitContainers[i].EnvFrom = append(
				container.EnvFrom, envFrom...)
		}
		for i, container := range resource.Spec.RolloutSpec.Template.Spec.Containers {
			resource.Spec.RolloutSpec.Template.Spec.Containers[i].EnvFrom = append(
				container.EnvFrom, envFrom...)
		}
	}

	return &rolloutsv1alpha1.Rollout{
		ObjectMeta: metav1.ObjectMeta{
			Name:        resource.Name,
			Labels:      resource.Labels,
			Namespace:   resource.Namespace,
			Annotations: make(map[string]string),
		},
		Spec: resource.Spec.RolloutSpec,
	}, nil
}

func (r *RolloutChildReconciler) MergeBeforeUpdate(current, desired *rolloutsv1alpha1.Rollout) {
	current.Labels = desired.Labels
	current.Spec = desired.Spec
}

func (r *RolloutChildReconciler) ReflectChildStatusOnParent(
	ctx context.Context, parent *argov1alpha1.Rollout, child *rolloutsv1alpha1.Rollout, err error) {

	// TODO: add status of configuration
	return
}

type rolloutChildReconciler = reconcilers.ChildReconciler[
	*argov1alpha1.Rollout,
	*rolloutsv1alpha1.Rollout,
	*rolloutsv1alpha1.RolloutList,
]