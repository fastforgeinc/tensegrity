package v1alpha1

import (
	"context"
	"encoding/base64"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reconciler.io/runtime/reconcilers"
)

const consumerSecretReconcilerName = "ConsumerSecretReconciler"
const consumerSecretKeysStashKey reconcilers.StashKey = "tensegrity.fastforge.io/consumer/secretKeys"
const consumerSecretNameStashKey reconcilers.StashKey = "tensegrity.fastforge.io/consumer/secretName"

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
	_ context.Context, _ *metav1.PartialObjectMetadata, _ *corev1.Secret, _ error) {
	return
}

func ConsumerSecretNameFromContext(ctx context.Context) string {
	if name, ok := reconcilers.RetrieveValue(ctx, consumerSecretNameStashKey).(string); ok {
		return name
	}
	return ""
}
