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
