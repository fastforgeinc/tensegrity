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
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/jsonpath"
	"k8s.io/utils/ptr"
	"reconciler.io/runtime/reconcilers"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/fastforgeinc/tensegrity/api/v1alpha1"
)

func NewProducerReconciler() *ProducerReconciler {
	r := new(ProducerReconciler)
	r.workloadReconciler = workloadReconciler{
		Name: "ProducerReconciler",
		Sync: r.Sync,
	}
	return r
}

type ProducerReconciler struct {
	workloadReconciler
}

func (r *ProducerReconciler) Sync(ctx context.Context, resource *v1alpha1.Tensegrity) error {
	resource.Status.Produced = nil
	resource.Status.ProducedKeys = nil

	if len(resource.Spec.Produces) == 0 {
		return nil
	}

	var seenError bool
	keys := make(map[string]string, len(resource.Spec.Produces))
	sensitiveKeys := make(map[string]string, len(resource.Spec.Produces))

	producedKeys := make([]v1alpha1.ProducedKeyStatus, 0, len(resource.Spec.Produces))
	for _, produces := range resource.Spec.Produces {
		var err error
		var value string
		var object *unstructured.Unstructured
		object, err = r.getObject(ctx, resource.Namespace, &produces)
		if err == nil {
			value, err = r.parseValue(object, &produces)
			switch {
			case produces.Sensitive && produces.Encoded:
				sensitiveKeys[produces.Key] = value
			case produces.Sensitive && !produces.Encoded:
				sensitiveKeys[produces.Key] = base64.StdEncoding.EncodeToString([]byte(value))
			default:
				keys[produces.Key] = value
			}
		}
		r.updateKeyStatus(object, resource, &produces, value, err)
		if err != nil {
			seenError = true
		}
	}
	sort.Slice(producedKeys, func(i, j int) bool {
		return producedKeys[i].Key < producedKeys[j].Key
	})
	resource.Status.ProducedKeys = producedKeys
	r.updateStatus(resource)

	if !seenError && len(keys) > 0 {
		reconcilers.StashValue(ctx, producerConfigMapKeysStashKey, keys)
		reconcilers.StashValue(ctx, producerConfigMapNameStashKey, resource.Spec.ProducesConfigMapName)
		resource.Status.ProducedConfigMapName = resource.Spec.ProducesConfigMapName
	} else {
		reconcilers.ClearValue(ctx, producerConfigMapKeysStashKey)
		reconcilers.ClearValue(ctx, producerConfigMapNameStashKey)
		resource.Status.ProducedConfigMapName = ""
	}

	if !seenError && len(sensitiveKeys) > 0 {
		reconcilers.StashValue(ctx, producerSecretKeysStashKey, sensitiveKeys)
		reconcilers.StashValue(ctx, producerSecretNameStashKey, resource.Spec.ProducesSecretName)
		resource.Status.ProducedSecretName = resource.Spec.ProducesSecretName
	} else {
		reconcilers.ClearValue(ctx, producerSecretKeysStashKey)
		reconcilers.ClearValue(ctx, producerSecretNameStashKey)
		resource.Status.ProducedSecretName = ""
	}

	return nil
}

func (r *ProducerReconciler) getObject(
	ctx context.Context, namespace string, produces *v1alpha1.ProducesSpec) (*unstructured.Unstructured, error) {

	if len(produces.Kind) == 0 && len(produces.APIVersion) == 0 {
		return new(unstructured.Unstructured), nil
	}

	config := reconcilers.RetrieveConfigOrDie(ctx)
	obj := new(unstructured.Unstructured)
	obj.SetKind(produces.Kind)
	obj.SetName(produces.Name)
	obj.SetNamespace(namespace)
	obj.SetAPIVersion(produces.APIVersion)

	err := config.TrackAndGet(ctx, client.ObjectKeyFromObject(obj), obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *ProducerReconciler) parseValue(
	obj *unstructured.Unstructured, produces *v1alpha1.ProducesSpec) (string, error) {

	jp := jsonpath.New(produces.Key)
	jp.AllowMissingKeys(false)
	if err := jp.Parse(produces.FieldPath); err != nil {
		return "", errors.Wrap(err, "fieldPath")
	}

	buf := new(bytes.Buffer)
	if err := jp.Execute(buf, obj.Object); err != nil {
		return "", errors.Wrap(err, "fieldPath")
	}

	value := buf.String()
	if len(value) == 0 {
		return "", errors.Wrap(errors.New("value is empty"), "fieldPath")
	}

	return buf.String(), nil
}

func (r *ProducerReconciler) updateKeyStatus(
	obj *unstructured.Unstructured, resource *v1alpha1.Tensegrity,
	produces *v1alpha1.ProducesSpec, value string, err error) {

	producedKeyStatus := v1alpha1.ProducedKeyStatus{
		Status:    v1alpha1.ProducedSuccess,
		Key:       produces.Key,
		Sensitive: produces.Sensitive,
	}

	if obj != nil {
		producedKeyStatus.ObjectReference = corev1.ObjectReference{
			Kind:            obj.GetKind(),
			Namespace:       obj.GetNamespace(),
			Name:            obj.GetName(),
			UID:             obj.GetUID(),
			APIVersion:      obj.GetAPIVersion(),
			FieldPath:       produces.FieldPath,
			ResourceVersion: obj.GetResourceVersion(),
		}
	}

	if len(value) > 0 && !produces.Sensitive {
		producedKeyStatus.Value = ptr.To(value)
	}

	if err != nil {
		producedKeyStatus.Status = v1alpha1.ProducedFailure
		producedKeyStatus.Reason = ptr.To(err.Error())
	}

	resource.Status.ProducedKeys = append(resource.Status.ProducedKeys, producedKeyStatus)
}

func (r *ProducerReconciler) updateStatus(resource *v1alpha1.Tensegrity) {
	resource.Status.Produced = ptr.To(v1alpha1.ProducedSuccess)
	condition := v1alpha1.NewTensegrityCondition(
		v1alpha1.TensegrityProduced, corev1.ConditionTrue, v1alpha1.KeysProducedReason, v1alpha1.KeysProducedMessage)

	var failedKeys []string
	for _, produced := range resource.Status.ProducedKeys {
		if produced.Status == v1alpha1.ProducedFailure {
			failedKeys = append(failedKeys, produced.Key)
		}
	}

	if len(failedKeys) > 0 {
		message := fmt.Sprintf(v1alpha1.KeysNotProducedMessage, strings.Join(failedKeys, ", "))
		condition = v1alpha1.NewTensegrityCondition(
			v1alpha1.TensegrityProduced, corev1.ConditionFalse,
			v1alpha1.KeysNotProducedReason, message)
		resource.Status.Produced = ptr.To(v1alpha1.ProducedFailure)
	}

	v1alpha1.SetTensegrityCondition(&resource.Status, *condition)
}
