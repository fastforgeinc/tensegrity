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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StaticSpec defines the desired state of Static
type StaticSpec struct {
	// TensegritySpec defines which keys a workload consumes and/or produces, and its delegates.
	TensegritySpec `json:",inline"`
}

// StaticStatus defines the observed state of Static
type StaticStatus struct {
	// Tensegrity status.
	TensegrityStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Produced",type=string,JSONPath=`.status.produced`
// +kubebuilder:printcolumn:name="Produced Config Map",type=string,JSONPath=`.status.producedConfigMapName`
// +kubebuilder:printcolumn:name="Produced Secret",type=string,JSONPath=`.status.producedSecretName`
// +kubebuilder:printcolumn:name="Consumed",type=string,JSONPath=`.status.consumed`
// +kubebuilder:printcolumn:name="Consumed Config Map",type=string,JSONPath=`.status.consumedConfigMapName`
// +kubebuilder:printcolumn:name="Consumed Secret",type=string,JSONPath=`.status.consumedSecretName`

// Static is the Schema for the statics API
type Static struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StaticSpec   `json:"spec,omitempty"`
	Status StaticStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StaticList contains a list of Static
type StaticList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Static `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Static{}, &StaticList{})
}
