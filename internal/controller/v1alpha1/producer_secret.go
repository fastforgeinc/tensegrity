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
	"encoding/base64"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reconciler.io/runtime/reconcilers"
)

const producerSecretReconcilerName = "ProducerSecretReconciler"
const producerSecretKeysStashKey reconcilers.StashKey = "tensegrity.fastforge.io/producer/secretKeys"
const producerSecretNameStashKey reconcilers.StashKey = "tensegrity.fastforge.io/producer/secretName"

func NewProducerSecretReconciler() *ProducerSecretReconciler {
	r := new(ProducerSecretReconciler)
	r.producerSecretReconciler = producerSecretReconciler{
		Name:                       producerSecretReconcilerName,
		OurChild:                   r.OurChild,
		DesiredChild:               r.DesiredChild,
		MergeBeforeUpdate:          r.MergeBeforeUpdate,
		ReflectChildStatusOnParent: r.ReflectChildStatusOnParent,
	}
	return r
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets/status,verbs=get

// ProducerSecretReconciler creates a Secret Kubernetes resource from produced keys.
type ProducerSecretReconciler struct {
	producerSecretReconciler
}

type producerSecretReconciler = reconcilers.ChildReconciler[
	*metav1.PartialObjectMetadata,
	*corev1.Secret,
	*corev1.SecretList,
]

func (r *ProducerSecretReconciler) DesiredChild(
	ctx context.Context, resource *metav1.PartialObjectMetadata) (*corev1.Secret, error) {

	keys, ok := reconcilers.RetrieveValue(ctx, producerSecretKeysStashKey).(map[string]string)
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
			Name:        reconcilers.RetrieveValue(ctx, producerSecretNameStashKey).(string),
			Labels:      resource.Labels,
			Namespace:   resource.Namespace,
			Annotations: map[string]string{"reconciler": producerSecretReconcilerName},
		},
		Type: corev1.SecretTypeOpaque,
		Data: encodedData,
	}
	return secret, nil
}

func (r *ProducerSecretReconciler) OurChild(_ *metav1.PartialObjectMetadata, child *corev1.Secret) bool {
	return child.Annotations["reconciler"] == producerSecretReconcilerName
}

func (r *ProducerSecretReconciler) MergeBeforeUpdate(current, desired *corev1.Secret) {
	current.Labels = desired.Labels
	current.Data = desired.Data
}

func (r *ProducerSecretReconciler) ReflectChildStatusOnParent(
	_ context.Context, _ *metav1.PartialObjectMetadata, _ *corev1.Secret, _ error) {
	return
}
