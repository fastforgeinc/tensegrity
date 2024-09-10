/*
This file is part of the Tensegrity distribution (https://github.com/fastforgeinc/tensegrity)
Copyright (C) 2024 FastForge Inc.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package v1alpha1

import (
	"context"
	"encoding/base64"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reconciler.io/runtime/reconcilers"
)

const consumerSecretReconcilerName = "ConsumerSecretReconciler"
const consumerSecretKeysStashKey reconcilers.StashKey = "tensegrity.fastforge.io/consumerSecretKeys"
const consumerSecretNameStashKey reconcilers.StashKey = "tensegrity.fastforge.io/consumerSecretName"
const consumerSecretVersionStashKey reconcilers.StashKey = "tensegrity.fastforge.io/consumerSecretVersion"

func NewConsumerSecretReconciler() *ConsumerSecretReconciler {
	r := new(ConsumerSecretReconciler)
	r.consumerSecretReconciler = consumerSecretReconciler{
		Name:                       consumerSecretReconcilerName,
		OurChild:                   r.OurChild,
		DesiredChild:               r.DesiredChild,
		MergeBeforeUpdate:          r.MergeBeforeUpdate,
		ReflectChildStatusOnParent: r.ReflectChildStatusOnParent,
	}
	return r
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets/status,verbs=get

// ConsumerSecretReconciler creates a Secret Kubernetes resource
// from consumed keys and binds to the related workload.
type ConsumerSecretReconciler struct {
	consumerSecretReconciler
}

type consumerSecretReconciler = reconcilers.ChildReconciler[
	*metav1.PartialObjectMetadata,
	*corev1.Secret,
	*corev1.SecretList,
]

func (r *ConsumerSecretReconciler) DesiredChild(
	ctx context.Context, resource *metav1.PartialObjectMetadata) (*corev1.Secret, error) {

	keys, ok := reconcilers.RetrieveValue(ctx, consumerSecretKeysStashKey).(map[string]string)
	if !ok || len(keys) == 0 {
		return nil, nil
	}

	encodedData := make(map[string][]byte, len(keys))
	for env, value := range keys {
		encodedValue, _ := base64.StdEncoding.DecodeString(value)
		encodedData[env] = encodedValue
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        reconcilers.RetrieveValue(ctx, consumerSecretNameStashKey).(string),
			Labels:      resource.Labels,
			Namespace:   resource.Namespace,
			Annotations: map[string]string{"reconciler": consumerSecretReconcilerName},
		},
		Type: corev1.SecretTypeOpaque,
		Data: encodedData,
	}
	return secret, nil
}

func (r *ConsumerSecretReconciler) MergeBeforeUpdate(current, desired *corev1.Secret) {
	current.Labels = desired.Labels
	current.Data = desired.Data
}

func (r *ConsumerSecretReconciler) OurChild(_ *metav1.PartialObjectMetadata, child *corev1.Secret) bool {
	return child.Annotations["reconciler"] == consumerSecretReconcilerName
}

func (r *ConsumerSecretReconciler) ReflectChildStatusOnParent(
	ctx context.Context, _ *metav1.PartialObjectMetadata, secret *corev1.Secret, _ error) {
	if secret != nil {
		if version := secret.GetResourceVersion(); len(version) > 0 {
			reconcilers.StashValue(ctx, consumerSecretVersionStashKey, version)
		}
	}
}

func ConsumerSecretNameFromContext(ctx context.Context) string {
	if name, ok := reconcilers.RetrieveValue(ctx, consumerSecretNameStashKey).(string); ok {
		return name
	}
	return ""
}

func ConsumerSecretAnnotationFromContext(ctx context.Context) (string, string) {
	if version, ok := reconcilers.RetrieveValue(ctx, consumerSecretVersionStashKey).(string); ok {
		return string(consumerSecretVersionStashKey), version
	}
	return "", ""
}
