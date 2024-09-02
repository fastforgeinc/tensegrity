package v1alpha1

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fastforgeinc/tensegrity/api/v1alpha1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/jsonpath"
	"k8s.io/utils/ptr"
	"reconciler.io/runtime/reconcilers"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

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
// for further processing by ConfigMapReconciler and SecretReconciler.
type ConsumerReconciler struct {
	workloadReconciler
}

func (r *ConsumerReconciler) Setup(ctx context.Context, _ ctrl.Manager, builder *builder.Builder) error {
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
		reconcilers.StashValue(ctx, configMapKeysStashKey, keys)
		reconcilers.StashValue(ctx, configMapNameStashKey, resource.Spec.ConfigMapName)
		resource.Status.ConfigMapName = resource.Spec.ConfigMapName
	}

	if len(sensitiveKeys) > 0 {
		reconcilers.StashValue(ctx, secretKeysStashKey, sensitiveKeys)
		reconcilers.StashValue(ctx, secretNameStashKey, resource.Spec.SecretName)
		resource.Status.SecretName = resource.Spec.SecretName
	}
	return nil
}

func (r *ConsumerReconciler) getKeys(
	ctx context.Context, resource *v1alpha1.Tensegrity) (keys, sensitiveKeys map[string]string, err error) {

	keys = make(map[string]string)
	sensitiveKeys = make(map[string]string)
	consumedByRef := make(map[corev1.ObjectReference]struct {
		v1alpha1.ConsumesSpec
		Delegate corev1.ObjectReference
	})
	consumesByRef := make(map[corev1.ObjectReference]v1alpha1.ConsumesSpec)
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
	consumedByRef map[corev1.ObjectReference]struct {
		v1alpha1.ConsumesSpec
		Delegate corev1.ObjectReference
	},
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

	for consumesRef, consumes := range consumesByRef {
		workload := v1alpha1.TensegrityFromRef(consumesRef)
		workload.SetNamespace(namespace.Name)
		err = config.TrackAndGet(ctx, client.ObjectKeyFromObject(workload), workload)
		if k8serrors.IsNotFound(err) {
			continue
		} else if err != nil {
			return err
		}
		reverseMaps := make(map[string]string, len(consumes.Maps))
		for env, key := range consumes.Maps {
			reverseMaps[key] = env
		}
		for _, produced := range workload.Status.ProducedKeys {
			env, ok := reverseMaps[produced.Key]
			if !ok {
				continue
			}

			object, err := r.getObject(ctx, namespace.Name, &produced)
			if k8serrors.IsNotFound(err) {
				continue
			} else if err != nil {
				return err
			}

			value, err := r.parseValue(object, &produced)
			if err != nil {
				return err
			}

			if produced.Sensitive {
				sensitiveKeys[env] = value
			} else {
				keys[env] = value
			}
			delete(reverseMaps, produced.Key)
		}
		if len(reverseMaps) == 0 {
			consumedByRef[consumesRef] = struct {
				v1alpha1.ConsumesSpec
				Delegate corev1.ObjectReference
			}{ConsumesSpec: consumes, Delegate: delegate}
			delete(consumesByRef, consumesRef)
		}
	}

	return nil
}

func (r *ConsumerReconciler) getObject(
	ctx context.Context, namespace string, produced *v1alpha1.ProducedKeyStatus) (*unstructured.Unstructured, error) {

	if len(produced.Kind) == 0 && len(produced.APIVersion) == 0 {
		return new(unstructured.Unstructured), nil
	}

	config := reconcilers.RetrieveConfigOrDie(ctx)
	obj := new(unstructured.Unstructured)
	obj.SetKind(produced.Kind)
	obj.SetName(produced.Name)
	obj.SetNamespace(namespace)
	obj.SetAPIVersion(produced.APIVersion)

	err := config.Get(ctx, client.ObjectKeyFromObject(obj), obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *ConsumerReconciler) parseValue(
	obj *unstructured.Unstructured, produced *v1alpha1.ProducedKeyStatus) (string, error) {

	jp := jsonpath.New(produced.Key)
	jp.AllowMissingKeys(false)
	if err := jp.Parse(produced.FieldPath); err != nil {
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
