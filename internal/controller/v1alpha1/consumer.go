package v1alpha1

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fastforgeinc/tensegrity/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/jsonpath"
	"reconciler.io/runtime/reconcilers"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

// ConsumerReconciler consumes keys from delegate workloads and puts into stash,
// for further processing by ConfigMapReconciler and SecretReconciler.
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=namespaces/status,verbs=get
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete
type ConsumerReconciler struct {
	workloadReconciler
}

func (r *ConsumerReconciler) Setup(ctx context.Context, _ ctrl.Manager, builder *builder.Builder) error {
	builder.Watches(new(corev1.Namespace), reconcilers.EnqueueTracked(ctx))
	return nil
}

func (r *ConsumerReconciler) Sync(ctx context.Context, resource *v1alpha1.Tensegrity) (err error) {
	if len(resource.Spec.Consumes) == 0 {
		return nil
	}

	if len(resource.Spec.Delegates) == 0 {
		return fmt.Errorf("no delegates specified")
	}

	keys, sensitiveKeys, err := r.getConsumedKeys(ctx, resource)
	if err != nil {
		return err
	}
	reconcilers.StashValue(ctx, secretKeysStashKey, sensitiveKeys)
	reconcilers.StashValue(ctx, secretNameStashKey, resource.Spec.SecretName)
	reconcilers.StashValue(ctx, configMapKeysStashKey, keys)
	reconcilers.StashValue(ctx, configMapNameStashKey, resource.Spec.ConfigMapName)
	return nil
}

func (r *ConsumerReconciler) getConsumedKeys(
	ctx context.Context, resource *v1alpha1.Tensegrity) (keys, sensitiveKeys map[string]string, err error) {

	keys = make(map[string]string)
	sensitiveKeys = make(map[string]string)
	consumesByRef := make(map[corev1.ObjectReference]v1alpha1.ConsumeSpec)
	for _, consume := range resource.Spec.Consumes {
		consumesByRef[consume.ObjectReference] = consume
	}
	for _, delegate := range resource.Spec.Delegates {
		if len(consumesByRef) == 0 {
			break
		}
		switch delegate.Kind {
		case "Namespace":
			if err = r.resolveKeysFromNamespace(ctx, delegate, consumesByRef, keys, sensitiveKeys); err != nil {
				return nil, nil, err
			}
		default:
			return nil, nil, fmt.Errorf("unsupported delegate kind: %s", delegate.Kind)
		}
	}

	return keys, sensitiveKeys, nil
}

func (r *ConsumerReconciler) resolveKeysFromNamespace(
	ctx context.Context, ref corev1.ObjectReference, consumesByRef map[corev1.ObjectReference]v1alpha1.ConsumeSpec,
	keys, sensitiveKeys map[string]string) error {

	config := reconcilers.RetrieveConfigOrDie(ctx)
	namespace := new(corev1.Namespace)
	namespace.SetName(ref.Name)
	err := config.TrackAndGet(ctx, client.ObjectKeyFromObject(namespace), namespace)
	if k8serrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	for consumeRef, consume := range consumesByRef {
		workload := v1alpha1.TensegrityFromRef(consumeRef)
		workload.SetNamespace(namespace.Name)
		err = config.TrackAndGet(ctx, client.ObjectKeyFromObject(workload), workload)
		if k8serrors.IsNotFound(err) {
			continue
		} else if err != nil {
			return err
		}
		reverseMaps := make(map[string]string, len(consume.Maps))
		for env, key := range consume.Maps {
			reverseMaps[key] = env
		}
		for _, produces := range workload.Status.Produces {
			env, ok := reverseMaps[produces.Key]
			if !ok {
				continue
			}

			obj := new(unstructured.Unstructured)
			obj.SetKind(produces.Kind)
			obj.SetName(produces.Name)
			obj.SetNamespace(namespace.Name)
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

			if produces.Sensitive {
				sensitiveKeys[env] = buf.String()
			} else {
				keys[env] = buf.String()
			}
		}
	}

	return nil
}
