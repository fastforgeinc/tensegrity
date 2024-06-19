/*
Copyright 2024 FastForge Inc. support@fastforge.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package v1alpha1 is common types for other versioned resources.
// +kubebuilder:object:generate=true
package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConsumeSpec struct {
	// Reference to a consumed object.
	corev1.ObjectReference `json:",inline"`
	// Maps defines mappings between consumed object keys and environment variables.
	Maps map[string]string `json:"maps,omitempty"`
}

type ProduceSpec struct {
	// Key is a name of a key is being produced.
	Key string `json:"key"`
	// ObjectReference is a reference to a Kubernetes resource as a source of value of the key is being produced.
	corev1.ObjectReference `json:",inline"`
	// Sensitive indicates that the produced key value must be hidden and consumed as a Secret.
	// +optional
	Sensitive bool `json:"sensitive,omitempty"`
}

type ProducesStatus struct {
	// Key is a name of a produced key.
	Key string `json:"key"`
	// ProduceFromKindSpec is a reference to a Kubernetes resource as a source of value of the produced key.
	corev1.ObjectReference `json:",inline"`
	// Sensitive indicates that the produced key value must be hidden and consumed as a Secret.
	// +optional
	Sensitive bool `json:"sensitive,omitempty"`
	// Value of the key resolved from Kubernetes resource.
	Value string `json:"value"`
}

// Tensegrity is a shared duck type for other reconcilers.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Tensegrity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TensegritySpec   `json:"spec"`
	Status TensegrityStatus `json:"status,omitempty"`
}

// TensegritySpec is Tensegrity controller specs.
type TensegritySpec struct {
	// Delegates is a slice of references to a Kubernetes resource used to resolve consumed keys.
	Delegates []corev1.ObjectReference `json:"delegates,omitempty"`
	// Consumes is a map of other workloads and ConsumeSpec.
	Consumes []ConsumeSpec `json:"consumes,omitempty"`
	// Produces is a map of keys and value sources to get from.
	Produces []ProduceSpec `json:"produces,omitempty"`
	// ConfigMapName is name of a ConfigMap is being generates by tensegrity controller,
	// defaults to <workload-name>-tensegrity.
	// +optional
	ConfigMapName string `json:"configMapName,omitempty"`
	// SecretName is name of a Secret is being generates by Tensegrity controller,
	// defaults to <workload-name>-tensegrity.
	// +optional
	SecretName string `json:"secretName,omitempty"`
}

// TensegrityStatus is Tensegrity controller status.
type TensegrityStatus struct {
	// Produces is a list of verified keys and their respective references with values.
	Produces []ProducesStatus `json:"produces,omitempty"`
	// SecretName is a name of a Secret with resolved environment variable and sensitive values
	// programmatically generated for a workload by Tensegrity controller.
	SecretName string `json:"secretName,omitempty"`
	// ConfigMapName is a name of a Secret with resolved environment variable and values
	// programmatically generated for a workload by Tensegrity controller.
	ConfigMapName string `json:"configMapName,omitempty"`
}

func TensegrityFromRef(ref corev1.ObjectReference) *Tensegrity {
	w := new(Tensegrity)
	w.SetName(ref.Name)
	w.SetGroupVersionKind(ref.GroupVersionKind())
	return w
}
