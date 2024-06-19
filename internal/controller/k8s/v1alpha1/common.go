package v1alpha1

import (
	k8sv1alpha1 "github.com/fastforgeinc/tensegrity/api/k8s/v1alpha1"
	"github.com/fastforgeinc/tensegrity/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TensegrityFromRef(ref corev1.ObjectReference) *v1alpha1.Tensegrity {
	w := new(v1alpha1.Tensegrity)
	w.SetName(ref.Name)
	w.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    ref.Kind,
		Group:   k8sv1alpha1.GroupVersion.Group,
		Version: k8sv1alpha1.GroupVersion.Version,
	})
	return w
}
