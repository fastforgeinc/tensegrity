/*
This file is part of the Tensegrity distribution (https://github.com/fastforgeinc/tensegrity)
Copyright (C) 2024 FastForge Inc.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package v1alpha1

import (
	"github.com/fastforgeinc/tensegrity/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentSpec defines the desired state of Deployment.
type DeploymentSpec struct {
	// DeploymentSpec is k8s.io/api/apps/v1.DeploymentSpec type.
	appsv1.DeploymentSpec `json:",inline"`
	// TensegritySpec defines which keys a workload consumes and/or produces, and its delegates.
	v1alpha1.TensegritySpec `json:",inline"`
}

// DeploymentStatus defines the observed state of Deployment.
type DeploymentStatus struct {
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

// Deployment is a wrapper type of the k8s.io/api/apps/v1.Deployment type.
type Deployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeploymentSpec   `json:"spec"`
	Status DeploymentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DeploymentList contains a list of Deployment.
type DeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Deployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Deployment{}, &DeploymentList{})
}
