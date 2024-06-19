package v1alpha1

import (
	"context"
	"encoding/base64"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reconciler.io/runtime/reconcilers"
)

const secretStashKey reconcilers.StashKey = "tensegrity.fastforge.io/secret"
const secretKeysStashKey reconcilers.StashKey = "tensegrity.fastforge.io/secretKeys"
const secretNameStashKey reconcilers.StashKey = "tensegrity.fastforge.io/secretName"

func NewSecretReconciler() *SecretReconciler {
	r := new(SecretReconciler)
	r.secretReconciler = secretReconciler{
		Name:                       "ConfigMapReconciler",
		DesiredChild:               r.DesiredChild,
		MergeBeforeUpdate:          r.MergeBeforeUpdate,
		ReflectChildStatusOnParent: r.ReflectChildStatusOnParent,
	}
	return r
}

// SecretReconciler creates a Secret Kubernetes resource
// from consumed keys and binds to the related workload.
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets/status,verbs=get
type SecretReconciler struct {
	secretReconciler
}

type secretReconciler = reconcilers.ChildReconciler[
	*metav1.PartialObjectMetadata,
	*corev1.Secret,
	*corev1.SecretList,
]

func (r *SecretReconciler) DesiredChild(
	ctx context.Context, resource *metav1.PartialObjectMetadata) (*corev1.Secret, error) {

	keys, ok := reconcilers.RetrieveValue(ctx, secretKeysStashKey).(map[string]string)
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
			Name:        resource.Name + "-tensegrity",
			Labels:      resource.Labels,
			Namespace:   resource.Namespace,
			Annotations: make(map[string]string),
		},
		Type: corev1.SecretTypeOpaque,
		Data: encodedData,
	}
	if name, ok := reconcilers.RetrieveValue(ctx, secretNameStashKey).(string); ok && len(name) > 0 {
		secret.ObjectMeta.Name = name
	}
	reconcilers.StashValue(ctx, secretStashKey, secret)
	return secret, nil
}

func (r *SecretReconciler) MergeBeforeUpdate(current, desired *corev1.Secret) {
	current.Labels = desired.Labels
	current.Data = desired.Data
}

func (r *SecretReconciler) ReflectChildStatusOnParent(
	_ context.Context, _ *metav1.PartialObjectMetadata, _ *corev1.Secret, _ error) {
	return
}

func SecretFromContext(ctx context.Context) *corev1.Secret {
	if secret, ok := reconcilers.RetrieveValue(ctx, secretStashKey).(*corev1.Secret); ok {
		return secret
	}
	return nil
}
