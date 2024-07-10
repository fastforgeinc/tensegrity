package v1alpha1

import (
	"bytes"
	"context"
	"github.com/fastforgeinc/tensegrity/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/jsonpath"
	"reconciler.io/runtime/reconcilers"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func (r *ProducerReconciler) Sync(ctx context.Context, resource *v1alpha1.Tensegrity) (err error) {
	if len(resource.Spec.Produces) == 0 {
		return nil
	}
	resource.Status.Produces = make([]v1alpha1.ProducesStatus, 0, len(resource.Spec.Produces))
	for _, produces := range resource.Spec.Produces {
		if err = r.updateStatus(ctx, resource, &produces); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProducerReconciler) updateStatus(
	ctx context.Context, resource *v1alpha1.Tensegrity, produces *v1alpha1.ProduceSpec) (err error) {

	config := reconcilers.RetrieveConfigOrDie(ctx)
	obj := new(unstructured.Unstructured)
	obj.SetKind(produces.Kind)
	obj.SetName(produces.Name)
	obj.SetNamespace(resource.Namespace)
	obj.SetAPIVersion(produces.APIVersion)

	err = config.Get(ctx, client.ObjectKeyFromObject(obj), obj)
	if err != nil {
		return err
	}

	jp := jsonpath.New(produces.Key)
	jp.AllowMissingKeys(false)
	if err = jp.Parse(produces.FieldPath); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = jp.Execute(buf, obj.Object); err != nil {
		return err
	}

	producesStatus := v1alpha1.ProducesStatus{
		Key:       produces.Key,
		Value:     "<redacted>",
		Sensitive: produces.Sensitive,
		ObjectReference: corev1.ObjectReference{
			Kind:            obj.GetKind(),
			Namespace:       obj.GetNamespace(),
			Name:            obj.GetName(),
			UID:             obj.GetUID(),
			APIVersion:      obj.GetAPIVersion(),
			FieldPath:       produces.FieldPath,
			ResourceVersion: obj.GetResourceVersion(),
		},
	}
	if !produces.Sensitive {
		producesStatus.Value = buf.String()
	}
	resource.Status.Produces = append(resource.Status.Produces, producesStatus)
	return nil
}
