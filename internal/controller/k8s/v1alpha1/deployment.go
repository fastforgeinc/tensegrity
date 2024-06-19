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

func NewDeploymentReconciler(config *reconcilers.Config) *DeploymentReconciler {
	return &DeploymentReconciler{
		Name: "DeploymentReconciler",
		Setup: func(ctx context.Context, _ ctrl.Manager, builder *builder.Builder) error {
			builder.Watches(new(k8sv1alpha1.Deployment), reconcilers.EnqueueTracked(ctx))
			return nil
		},
		Config: *config,
		Reconciler: reconcilers.Sequence[*k8sv1alpha1.Deployment]{
			&reconcilers.CastResource[*k8sv1alpha1.Deployment, *apiv1alpha1.Tensegrity]{
				Reconciler: v1alpha1.NewConsumerReconciler(),
			},
			&reconcilers.CastResource[*k8sv1alpha1.Deployment, *metav1.PartialObjectMetadata]{
				Reconciler: v1alpha1.NewConfigMapReconciler(),
			},
			&reconcilers.CastResource[*k8sv1alpha1.Deployment, *metav1.PartialObjectMetadata]{
				Reconciler: v1alpha1.NewSecretReconciler(),
			},
			NewDeploymentChildReconciler(),
			&reconcilers.CastResource[*k8sv1alpha1.Deployment, *apiv1alpha1.Tensegrity]{
				Reconciler: v1alpha1.NewProducerReconciler(),
			},
		},
	}
}

// DeploymentReconciler reconciles tensegrity v1alpha1.Deployment resource,
// and runs sequence of other reconcilers to get desired workload.
// +kubebuilder:rbac:groups=tensegrity.fastforge.io,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tensegrity.fastforge.io,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=tensegrity.fastforge.io,resources=deployments/finalizers,verbs=update
type DeploymentReconciler = reconcilers.ResourceReconciler[*k8sv1alpha1.Deployment]

func NewDeploymentChildReconciler() *DeploymentChildReconciler {
	r := new(DeploymentChildReconciler)
	r.deploymentChildReconciler = deploymentChildReconciler{
		Name:                       "DeploymentChildReconciler",
		DesiredChild:               r.DesiredChild,
		MergeBeforeUpdate:          r.MergeBeforeUpdate,
		ReflectChildStatusOnParent: r.ReflectChildStatusOnParent,
	}
	return r
}

// DeploymentChildReconciler creates child deployment from workload specs,
// add ConfigMap and Secret if they are present.
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update
type DeploymentChildReconciler struct {
	deploymentChildReconciler
}

func (r *DeploymentChildReconciler) DesiredChild(
	ctx context.Context, resource *k8sv1alpha1.Deployment) (*appsv1.Deployment, error) {

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
		for i, container := range resource.Spec.DeploymentSpec.Template.Spec.InitContainers {
			resource.Spec.DeploymentSpec.Template.Spec.InitContainers[i].EnvFrom = append(
				container.EnvFrom, envFrom...)
		}
		for i, container := range resource.Spec.DeploymentSpec.Template.Spec.Containers {
			resource.Spec.DeploymentSpec.Template.Spec.Containers[i].EnvFrom = append(
				container.EnvFrom, envFrom...)
		}
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        resource.Name,
			Labels:      resource.Labels,
			Namespace:   resource.Namespace,
			Annotations: make(map[string]string),
		},
		Spec: resource.Spec.DeploymentSpec,
	}, nil
}

func (r *DeploymentChildReconciler) MergeBeforeUpdate(current, desired *appsv1.Deployment) {
	current.Labels = desired.Labels
	current.Spec = desired.Spec
}

func (r *DeploymentChildReconciler) ReflectChildStatusOnParent(
	ctx context.Context, parent *k8sv1alpha1.Deployment, child *appsv1.Deployment, err error) {

	// TODO: add status of configuration
	return
}

type deploymentChildReconciler = reconcilers.ChildReconciler[
	*k8sv1alpha1.Deployment,
	*appsv1.Deployment,
	*appsv1.DeploymentList,
]
