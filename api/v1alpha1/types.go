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

// +kubebuilder:object:generate=true

// Package v1alpha1 is common types for other versioned resources.
package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConsumedStatus string

const (
	ConsumedSuccess ConsumedStatus = "Success"
	ConsumedFailure ConsumedStatus = "Failure"
)

type ConsumesSpec struct {
	// ObjectReference to an object is being consumed.
	corev1.ObjectReference `json:",inline"`
	// Maps defines mappings between consumed object keys and ConfigMap/Secret keys.
	Maps map[string]string `json:"maps,omitempty"`
}

type ConsumedKeyStatus struct {
	// ObjectReference to a Tensegrity resource a key consumed from.
	corev1.ObjectReference `json:",inline"`
	// Delegate is a ObjectReference to a resource key is consumed from.
	Delegate *corev1.ObjectReference `json:"delegate,omitempty"`
	// Status of a key.
	Status ConsumedStatus `json:"status"`
	// Reason of a status.
	Reason *string `json:"reason,omitempty"`
	// Key is a name of a consumed key.
	Key string `json:"key"`
	// Env is a name of a consumed env.
	Env string `json:"env"`
}

type ProducedStatus string

const (
	ProducedSuccess ProducedStatus = "Success"
	ProducedFailure ProducedStatus = "Failure"
)

type ProducesSpec struct {
	// Key is a name of a key is being produced.
	Key string `json:"key"`
	// ObjectReference is a reference to a Kubernetes resource as a source of value of the key is being produced.
	corev1.ObjectReference `json:",inline"`
	// Sensitive indicates that the produced key value must be hidden and consumed as a Secret.
	// +optional
	Sensitive bool `json:"sensitive,omitempty"`
}

type ProducedKeyStatus struct {
	// ObjectReference to a Kubernetes resource a key produced from.
	corev1.ObjectReference `json:",inline"`
	// Status of a key.
	Status ProducedStatus `json:"status"`
	// Reason of a status.
	Reason *string `json:"reason,omitempty"`
	// Key is a name of a produced key.
	Key string `json:"key"`
	// Sensitive indicates that the produced key value must be hidden and represented as a Secret.
	// +optional
	Sensitive bool `json:"sensitive,omitempty"`
	// Value of the key resolved from Kubernetes resource.
	// +optional
	Value *string `json:"value,omitempty"`
}

// +kubebuilder:skipversion
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Tensegrity is a shared duck type for other reconcilers.
type Tensegrity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TensegritySpec   `json:"spec"`
	Status TensegrityStatus `json:"status,omitempty"`
}

// TensegritySpec is Tensegrity controller specs.
type TensegritySpec struct {
	// Delegates is a list of ObjectReference to a Kubernetes resource used to resolve consumed keys,
	// if empty defaults to a resource namespace.
	// +optional
	Delegates []corev1.ObjectReference `json:"delegates,omitempty"`
	// Consumes is a map of other workloads and ConsumeSpec.
	// +optional
	Consumes []ConsumesSpec `json:"consumes,omitempty"`
	// Produces is a map of keys and value sources to get from.
	// +optional
	Produces []ProducesSpec `json:"produces,omitempty"`
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
	// Consumed indicates whether all keys were consumed.
	Consumed *ConsumedStatus `json:"consumed,omitempty"`
	// ConsumedKeys indicates consumed keys and their statuses.
	ConsumedKeys []ConsumedKeyStatus `json:"consumedKeys,omitempty"`
	// Produced indicates whether all keys were produced.
	Produced *ProducedStatus `json:"produced,omitempty"`
	// ProducedKeys indicates produced keys and their statuses.
	ProducedKeys []ProducedKeyStatus `json:"producedKeys,omitempty"`
	// SecretName is a name of a Secret with resolved environment variable and sensitive values
	// programmatically generated for a workload by Tensegrity controller.
	SecretName string `json:"secretName,omitempty"`
	// ConfigMapName is a name of a Secret with resolved environment variable and values
	// programmatically generated for a workload by Tensegrity controller.
	ConfigMapName string `json:"configMapName,omitempty"`
	// Conditions a list of conditions a tensegrity resource can have.
	// +optional
	Conditions []TensegrityCondition `json:"conditions,omitempty"`
	// ObservedGeneration is the 'Generation' of the resource that
	// was last processed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

func TensegrityFromRef(ref corev1.ObjectReference) *Tensegrity {
	w := new(Tensegrity)
	w.SetName(ref.Name)
	w.SetGroupVersionKind(ref.GroupVersionKind())
	return w
}
