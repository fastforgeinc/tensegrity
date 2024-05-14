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
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatefulSetSpec defines the desired state of StatefulSet
type StatefulSetSpec struct {
	// Kubernetes stateful set spec.
	v1.StatefulSetSpec `json:",inline"`
	// WorkloadSpec defines which keys a stateful set consumes and produces, and a deployment delegates.
	WorkloadSpec `json:",inline"`
}

// StatefulSetStatus defines the observed state of StatefulSet
type StatefulSetStatus struct {
	// Kubernetes stateful set status.
	v1.StatefulSetStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// StatefulSet is the Schema for the statefulsets API
type StatefulSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StatefulSetSpec   `json:"spec,omitempty"`
	Status StatefulSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StatefulSetList contains a list of StatefulSet
type StatefulSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StatefulSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StatefulSet{}, &StatefulSetList{})
}
