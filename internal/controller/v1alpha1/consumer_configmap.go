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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reconciler.io/runtime/reconcilers"
)

const consumerConfigMapReconcilerName = "ConsumerConfigMapReconciler"
const consumerConfigMapKeysStashKey reconcilers.StashKey = "tensegrity.fastforge.io/consumerConfigMapKeys"
const consumerConfigMapNameStashKey reconcilers.StashKey = "tensegrity.fastforge.io/consumerConfigMapName"
const consumerConfigMapVersionStashKey reconcilers.StashKey = "tensegrity.fastforge.io/consumerConfigMapVersion"

func NewConsumerConfigMapReconciler() *ConsumerConfigMapReconciler {
	r := new(ConsumerConfigMapReconciler)
	r.consumerConfigMapReconciler = consumerConfigMapReconciler{
		Name:                       consumerConfigMapReconcilerName,
		OurChild:                   r.OurChild,
		DesiredChild:               r.DesiredChild,
		ReflectChildStatusOnParent: r.ReflectChildStatusOnParent,
		ChildObjectManager: &reconcilers.UpdatingObjectManager[*corev1.ConfigMap]{
			MergeBeforeUpdate: r.MergeBeforeUpdate,
		},
	}
	return r
}

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps/status,verbs=get

// ConsumerConfigMapReconciler creates a ConfigMap Kubernetes resource
// from consumed keys and binds to the related workload.
type ConsumerConfigMapReconciler struct {
	consumerConfigMapReconciler
}

type consumerConfigMapReconciler = reconcilers.ChildReconciler[
	*metav1.PartialObjectMetadata,
	*corev1.ConfigMap,
	*corev1.ConfigMapList,
]

func (r *ConsumerConfigMapReconciler) DesiredChild(
	ctx context.Context, resource *metav1.PartialObjectMetadata) (*corev1.ConfigMap, error) {

	keys, ok := reconcilers.RetrieveValue(ctx, consumerConfigMapKeysStashKey).(map[string]string)
	if !ok || len(keys) == 0 {
		return nil, nil
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        reconcilers.RetrieveValue(ctx, consumerConfigMapNameStashKey).(string),
			Labels:      resource.Labels,
			Namespace:   resource.Namespace,
			Annotations: map[string]string{"reconciler": consumerConfigMapReconcilerName},
		},
		Data: keys,
	}

	return configMap, nil
}

func (r *ConsumerConfigMapReconciler) MergeBeforeUpdate(current, desired *corev1.ConfigMap) {
	current.Labels = desired.Labels
	current.Data = desired.Data
	current.BinaryData = desired.BinaryData
}

func (r *ConsumerConfigMapReconciler) OurChild(_ *metav1.PartialObjectMetadata, child *corev1.ConfigMap) bool {
	return child.Annotations["reconciler"] == consumerConfigMapReconcilerName
}

func (r *ConsumerConfigMapReconciler) ReflectChildStatusOnParent(
	ctx context.Context, _ *metav1.PartialObjectMetadata, configMap *corev1.ConfigMap, _ error) {
	if configMap != nil {
		if version := configMap.GetResourceVersion(); len(version) > 0 {
			reconcilers.StashValue(ctx, consumerConfigMapVersionStashKey, version)
		}
	}
}

func ConsumerConfigMapNameFromContext(ctx context.Context) string {
	if name, ok := reconcilers.RetrieveValue(ctx, consumerConfigMapNameStashKey).(string); ok {
		return name
	}
	return ""
}

func ConsumerConfigMapAnnotationFromContext(ctx context.Context) (string, string) {
	if version, ok := reconcilers.RetrieveValue(ctx, consumerConfigMapVersionStashKey).(string); ok {
		return string(consumerConfigMapVersionStashKey), version
	}
	return "", ""
}
