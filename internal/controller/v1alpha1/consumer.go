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
	"fmt"
	"strings"

	"github.com/fastforgeinc/tensegrity/api/v1alpha1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/ptr"
	"reconciler.io/runtime/reconcilers"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type consumedDelegate struct {
	v1alpha1.ConsumesSpec
	Delegate corev1.ObjectReference
}

func NewConsumerReconciler() *ConsumerReconciler {
	r := new(ConsumerReconciler)
	r.workloadReconciler = workloadReconciler{
		Name:  "ConsumerReconciler",
		Sync:  r.Sync,
		Setup: r.Setup,
	}
	return r
}

// +kubebuilder:rbac:groups="*",resources=*,verbs=get;list;watch
// +kubebuilder:rbac:groups="*",resources=*/status,verbs=get
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete

// ConsumerReconciler consumes keys from delegate workloads and puts into stash,
// for further processing by ConsumerConfigMapReconciler and ConsumerSecretReconciler.
type ConsumerReconciler struct {
	workloadReconciler
}

func (r *ConsumerReconciler) Setup(ctx context.Context, _ ctrl.Manager, builder *builder.Builder) error {
	builder.Watches(new(corev1.Secret), reconcilers.EnqueueTracked(ctx))
	builder.Watches(new(corev1.ConfigMap), reconcilers.EnqueueTracked(ctx))
	builder.Watches(new(corev1.Namespace), reconcilers.EnqueueTracked(ctx))
	return nil
}

func (r *ConsumerReconciler) Sync(ctx context.Context, resource *v1alpha1.Tensegrity) (err error) {
	resource.Status.Consumed = nil
	resource.Status.ConsumedKeys = nil
	v1alpha1.RemoveTensegrityCondition(&resource.Status, v1alpha1.TensegrityConsumed)

	if len(resource.Spec.Consumes) == 0 {
		return nil
	}

	keys, sensitiveKeys, err := r.getKeys(ctx, resource)
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		reconcilers.StashValue(ctx, consumerConfigMapKeysStashKey, keys)
		reconcilers.StashValue(ctx, consumerConfigMapNameStashKey, resource.Spec.ConsumesConfigMapName)
		resource.Status.ConsumedConfigMapName = resource.Spec.ConsumesConfigMapName
	}

	if len(sensitiveKeys) > 0 {
		reconcilers.StashValue(ctx, consumerSecretKeysStashKey, sensitiveKeys)
		reconcilers.StashValue(ctx, consumerSecretNameStashKey, resource.Spec.ConsumesSecretName)
		resource.Status.ConsumedSecretName = resource.Spec.ConsumesSecretName
	}
	return nil
}

func (r *ConsumerReconciler) getKeys(
	ctx context.Context, resource *v1alpha1.Tensegrity) (keys, sensitiveKeys map[string]string, err error) {

	keys = make(map[string]string)
	sensitiveKeys = make(map[string]string)
	consumesByRef := make(map[corev1.ObjectReference]v1alpha1.ConsumesSpec)
	consumedByRef := make(map[corev1.ObjectReference]consumedDelegate)
	for _, consumes := range resource.Spec.Consumes {
		consumesByRef[consumes.ObjectReference] = consumes
	}
	for _, delegate := range resource.Spec.Delegates {
		if len(consumesByRef) == 0 {
			break
		}
		switch delegate.Kind {
		case "Namespace":
			if err = r.getKeysFromNamespace(
				ctx, delegate, consumesByRef, consumedByRef, keys, sensitiveKeys); err != nil {

				return nil, nil, err
			}
		default:
			return nil, nil, fmt.Errorf("unsupported delegate kind: %s", delegate.Kind)
		}
	}

	for consumedRef, consumed := range consumedByRef {
		r.updateKeyStatus(resource, ptr.To(consumed.Delegate), consumedRef, consumed.ConsumesSpec, nil)
	}

	err = errors.Errorf("consumed key by reference is not found")
	for consumesRef, consumes := range consumesByRef {
		r.updateKeyStatus(resource, nil, consumesRef, consumes, err)
	}

	r.updateStatus(resource)
	return keys, sensitiveKeys, nil
}

func (r *ConsumerReconciler) getKeysFromNamespace(
	ctx context.Context, delegate corev1.ObjectReference,
	consumesByRef map[corev1.ObjectReference]v1alpha1.ConsumesSpec,
	consumedByRef map[corev1.ObjectReference]consumedDelegate,
	keys, sensitiveKeys map[string]string) error {

	config := reconcilers.RetrieveConfigOrDie(ctx)
	namespace := new(corev1.Namespace)
	namespace.SetName(delegate.Name)
	err := config.TrackAndGet(ctx, client.ObjectKeyFromObject(namespace), namespace)
	if k8serrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

ConsumesByRefLoop:
	for consumesRef, consumes := range consumesByRef {
		tensegrity := v1alpha1.TensegrityFromRef(consumesRef)
		tensegrity.SetNamespace(namespace.Name)
		err = config.TrackAndGet(ctx, client.ObjectKeyFromObject(tensegrity), tensegrity)
		if k8serrors.IsNotFound(err) {
			continue
		} else if err != nil {
			return err
		}

		configMap := new(corev1.ConfigMap)
		if len(tensegrity.Status.ProducedConfigMapName) > 0 {
			configMap.SetName(tensegrity.Status.ProducedConfigMapName)
			configMap.SetNamespace(namespace.Name)
			err = config.TrackAndGet(ctx, client.ObjectKeyFromObject(configMap), configMap)
			if k8serrors.IsNotFound(err) {
				continue
			} else if err != nil {
				return err
			}
		}

		secret := new(corev1.Secret)
		if len(tensegrity.Status.ProducedSecretName) > 0 {
			secret.SetName(tensegrity.Status.ProducedSecretName)
			secret.SetNamespace(namespace.Name)
			err = config.TrackAndGet(ctx, client.ObjectKeyFromObject(secret), secret)
			if k8serrors.IsNotFound(err) {
				continue
			} else if err != nil {
				return err
			}
		}

		localKeys := make(map[string]string, len(configMap.Data))
		localSensitiveKeys := make(map[string]string, len(secret.Data))
		for env, key := range consumes.Maps {
			if v, ok := configMap.Data[key]; ok {
				localKeys[env] = v
				continue
			}
			if v, ok := secret.Data[key]; ok {
				localSensitiveKeys[env] = base64.StdEncoding.EncodeToString(v)
				continue
			}
			continue ConsumesByRefLoop
		}

		for env, v := range localKeys {
			keys[env] = v
		}
		for env, v := range localSensitiveKeys {
			sensitiveKeys[env] = v
		}
		consumedByRef[consumesRef] = consumedDelegate{
			ConsumesSpec: consumes,
			Delegate:     delegate,
		}
		delete(consumesByRef, consumesRef)
	}

	return nil
}

func (r *ConsumerReconciler) updateKeyStatus(
	resource *v1alpha1.Tensegrity, delegate *corev1.ObjectReference, ref corev1.ObjectReference,
	consumes v1alpha1.ConsumesSpec, err error) {

	for env, key := range consumes.Maps {
		consumedKeyStatus := v1alpha1.ConsumedKeyStatus{
			ObjectReference: ref,
			Delegate:        delegate,
			Status:          v1alpha1.ConsumedSuccess,
			Key:             key,
			Env:             env,
		}
		if err != nil {
			consumedKeyStatus.Status = v1alpha1.ConsumedFailure
			consumedKeyStatus.Reason = ptr.To(err.Error())
		}
		resource.Status.ConsumedKeys = append(resource.Status.ConsumedKeys, consumedKeyStatus)
	}
}

func (r *ConsumerReconciler) updateStatus(resource *v1alpha1.Tensegrity) {
	resource.Status.Consumed = ptr.To(v1alpha1.ConsumedSuccess)
	condition := v1alpha1.NewTensegrityCondition(
		v1alpha1.TensegrityConsumed, corev1.ConditionTrue, v1alpha1.KeysConsumedReason, v1alpha1.KeysConsumedMessage)

	var failedEnvs []string
	for _, consumed := range resource.Status.ConsumedKeys {
		if consumed.Status == v1alpha1.ConsumedFailure {
			failedEnvs = append(failedEnvs, consumed.Env)
		}
	}

	if len(failedEnvs) > 0 {
		message := fmt.Sprintf(v1alpha1.KeysNotConsumedMessage, strings.Join(failedEnvs, ", "))
		condition = v1alpha1.NewTensegrityCondition(
			v1alpha1.TensegrityConsumed, corev1.ConditionFalse,
			v1alpha1.KeysNotConsumedReason, message)
		resource.Status.Consumed = ptr.To(v1alpha1.ConsumedFailure)
	}

	v1alpha1.SetTensegrityCondition(&resource.Status, *condition)
}
