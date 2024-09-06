package v1alpha1

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reconciler.io/runtime/reconcilers"
)

const consumerConfigMapReconcilerName = "ConsumerConfigMapReconciler"
const consumerConfigMapKeysStashKey reconcilers.StashKey = "tensegrity.fastforge.io/consumer/configMapKeys"
const consumerConfigMapNameStashKey reconcilers.StashKey = "tensegrity.fastforge.io/consumer/configMapName"

func NewConsumerConfigMapReconciler() *ConsumerConfigMapReconciler {
	r := new(ConsumerConfigMapReconciler)
	r.consumerConfigMapReconciler = consumerConfigMapReconciler{
		Name:                       consumerConfigMapReconcilerName,
		OurChild:                   r.OurChild,
		DesiredChild:               r.DesiredChild,
		MergeBeforeUpdate:          r.MergeBeforeUpdate,
		ReflectChildStatusOnParent: r.ReflectChildStatusOnParent,
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
	_ context.Context, _ *metav1.PartialObjectMetadata, _ *corev1.ConfigMap, _ error) {
	return
}

func ConsumerConfigMapNameFromContext(ctx context.Context) string {
	if name, ok := reconcilers.RetrieveValue(ctx, consumerConfigMapNameStashKey).(string); ok {
		return name
	}
	return ""
}
