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

package v1alpha1

import (
	"github.com/fastforgeinc/tensegrity/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DaemonSetSpec defines the desired state of DaemonSet.
type DaemonSetSpec struct {
	// DaemonSetSpec is k8s.io/api/apps/v1.DaemonSetSpec type.
	appsv1.DaemonSetSpec `json:",inline"`
	// TensegritySpec defines which keys a workload consumes and/or produces, and its delegates.
	v1alpha1.TensegritySpec `json:",inline"`
}

// DaemonSetStatus defines the observed state of DaemonSet
type DaemonSetStatus struct {
	// Tensegrity status.
	v1alpha1.TensegrityStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Produced",type=string,JSONPath=`.status.produced`
// +kubebuilder:printcolumn:name="Produced Config Map",type=string,JSONPath=`.status.producedConfigMapName`
// +kubebuilder:printcolumn:name="Produced Secret",type=string,JSONPath=`.status.producedSecretName`
// +kubebuilder:printcolumn:name="Consumed",type=string,JSONPath=`.status.consumed`
// +kubebuilder:printcolumn:name="Consumed Config Map",type=string,JSONPath=`.status.consumedConfigMapName`
// +kubebuilder:printcolumn:name="Consumed Secret",type=string,JSONPath=`.status.consumedSecretName`

// DaemonSet is a wrapper type of the k8s.io/api/apps/v1.DaemonSet type.
type DaemonSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DaemonSetSpec   `json:"spec"`
	Status DaemonSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DaemonSetList contains a list of DaemonSet.
type DaemonSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DaemonSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DaemonSet{}, &DaemonSetList{})
}
