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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reconciler.io/runtime/reconcilers"
)

const producerConfigMapReconcilerName = "ProducerConfigMapReconciler"
const producerConfigMapKeysStashKey reconcilers.StashKey = "tensegrity.fastforge.io/producer/configMapKeys"
const producerConfigMapNameStashKey reconcilers.StashKey = "tensegrity.fastforge.io/producer/configMapName"

func NewProducerConfigMapReconciler() *ProducerConfigMapReconciler {
	r := new(ProducerConfigMapReconciler)
	r.producerConfigMapReconciler = producerConfigMapReconciler{
		Name:                       producerConfigMapReconcilerName,
		OurChild:                   r.OurChild,
		DesiredChild:               r.DesiredChild,
		MergeBeforeUpdate:          r.MergeBeforeUpdate,
		ReflectChildStatusOnParent: r.ReflectChildStatusOnParent,
	}
	return r
}

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps/status,verbs=get

// ProducerConfigMapReconciler creates a ConfigMap Kubernetes resource from produced keys.
type ProducerConfigMapReconciler struct {
	producerConfigMapReconciler
}

type producerConfigMapReconciler = reconcilers.ChildReconciler[
	*metav1.PartialObjectMetadata,
	*corev1.ConfigMap,
	*corev1.ConfigMapList,
]

func (r *ProducerConfigMapReconciler) DesiredChild(
	ctx context.Context, resource *metav1.PartialObjectMetadata) (*corev1.ConfigMap, error) {

	keys, ok := reconcilers.RetrieveValue(ctx, producerConfigMapKeysStashKey).(map[string]string)
	if !ok || len(keys) == 0 {
		return nil, nil
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        reconcilers.RetrieveValue(ctx, producerConfigMapNameStashKey).(string),
			Labels:      resource.Labels,
			Namespace:   resource.Namespace,
			Annotations: map[string]string{"reconciler": producerConfigMapReconcilerName},
		},
		Data: keys,
	}

	return configMap, nil
}

func (r *ProducerConfigMapReconciler) MergeBeforeUpdate(current, desired *corev1.ConfigMap) {
	current.Labels = desired.Labels
	current.Data = desired.Data
	current.BinaryData = desired.BinaryData
}

func (r *ProducerConfigMapReconciler) OurChild(_ *metav1.PartialObjectMetadata, child *corev1.ConfigMap) bool {
	return child.Annotations["reconciler"] == producerConfigMapReconcilerName
}

func (r *ProducerConfigMapReconciler) ReflectChildStatusOnParent(
	_ context.Context, _ *metav1.PartialObjectMetadata, _ *corev1.ConfigMap, _ error) {
	return
}
