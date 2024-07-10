package v1alpha1

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reconciler.io/runtime/reconcilers"
)

const configMapStashKey reconcilers.StashKey = "tensegrity.fastforge.io/configMap"
const configMapKeysStashKey reconcilers.StashKey = "tensegrity.fastforge.io/configMapKeys"
const configMapNameStashKey reconcilers.StashKey = "tensegrity.fastforge.io/configMapName"

func NewConfigMapReconciler() *ConfigMapReconciler {
	r := new(ConfigMapReconciler)
	r.configMapReconciler = configMapReconciler{
		Name:                       "ConfigMapReconciler",
		DesiredChild:               r.DesiredChild,
		MergeBeforeUpdate:          r.MergeBeforeUpdate,
		ReflectChildStatusOnParent: r.ReflectChildStatusOnParent,
	}
	return r
}

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps/status,verbs=get

// ConfigMapReconciler creates a ConfigMap Kubernetes resource
// from consumed keys and binds to the related workload.
type ConfigMapReconciler struct {
	configMapReconciler
}

type configMapReconciler = reconcilers.ChildReconciler[
	*metav1.PartialObjectMetadata,
	*corev1.ConfigMap,
	*corev1.ConfigMapList,
]

func (r *ConfigMapReconciler) DesiredChild(
	ctx context.Context, resource *metav1.PartialObjectMetadata) (*corev1.ConfigMap, error) {

	keys, ok := reconcilers.RetrieveValue(ctx, configMapKeysStashKey).(map[string]string)
	if !ok || len(keys) == 0 {
		return nil, nil
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        reconcilers.RetrieveValue(ctx, configMapNameStashKey).(string),
			Labels:      resource.Labels,
			Namespace:   resource.Namespace,
			Annotations: make(map[string]string),
		},
		Data: keys,
	}

	reconcilers.StashValue(ctx, configMapStashKey, configMap)
	return configMap, nil
}

func (r *ConfigMapReconciler) MergeBeforeUpdate(current, desired *corev1.ConfigMap) {
	current.Labels = desired.Labels
	current.Data = desired.Data
	current.BinaryData = desired.BinaryData
}

func (r *ConfigMapReconciler) ReflectChildStatusOnParent(
	_ context.Context, _ *metav1.PartialObjectMetadata, _ *corev1.ConfigMap, _ error) {
	return
}

func ConfigMapFromContext(ctx context.Context) *corev1.ConfigMap {
	if configMap, ok := reconcilers.RetrieveValue(ctx, configMapStashKey).(*corev1.ConfigMap); ok {
		return configMap
	}
	return nil
}
