package v1alpha1

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fastforgeinc/tensegrity/api/v1alpha1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/jsonpath"
	"k8s.io/utils/ptr"
	"reconciler.io/runtime/reconcilers"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
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
	v1alpha1.RemoveTensegrityCondition(&resource.Status, v1alpha1.TensegrityProduced)

	if len(resource.Spec.Produces) == 0 {
		return nil
	}

	resource.Status.ProducedKeys = make([]v1alpha1.ProducedKeyStatus, 0, len(resource.Spec.Produces))
	for _, produces := range resource.Spec.Produces {
		var err error
		var value string
		var object *unstructured.Unstructured
		object, err = r.getObject(ctx, resource.Namespace, &produces)
		if err == nil {
			value, err = r.parseValue(object, &produces)
		}
		r.updateKeyStatus(object, resource, &produces, value, err)
	}
	r.updateStatus(resource)
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
