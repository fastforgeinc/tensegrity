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
	"encoding/json"
	rolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	"github.com/fastforgeinc/tensegrity/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// RolloutSpec defines the desired state of Rollout
type RolloutSpec struct {
	// DeploymentSpec is github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1.RolloutSpec type.
	rolloutsv1alpha1.RolloutSpec `json:",inline"`
	// TensegritySpec defines which keys a workload consumes and/or produces, and its delegates.
	v1alpha1.TensegritySpec `json:",inline"`
}

func (s *RolloutSpec) MarshalJSON() ([]byte, error) {
	type RolloutAlias rolloutsv1alpha1.RolloutSpec
	type TensegrityAlias v1alpha1.TensegritySpec

	if s.TemplateResolvedFromRef || s.SelectorResolvedFromRef {
		obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&struct {
			RolloutAlias    `json:",inline"`
			TensegrityAlias `json:",inline"`
		}{
			RolloutAlias:    (RolloutAlias)(s.RolloutSpec),
			TensegrityAlias: (TensegrityAlias)(s.TensegritySpec),
		})
		if err != nil {
			return nil, err
		}
		if s.TemplateResolvedFromRef {
			unstructured.RemoveNestedField(obj, "template")
		}
		if s.SelectorResolvedFromRef {
			unstructured.RemoveNestedField(obj, "selector")
		}

		return json.Marshal(obj)
	}
	return json.Marshal(&struct {
		RolloutAlias    `json:",inline"`
		TensegrityAlias `json:",inline"`
	}{
		RolloutAlias:    (RolloutAlias)(s.RolloutSpec),
		TensegrityAlias: (TensegrityAlias)(s.TensegritySpec),
	})
}

// RolloutStatus defines the observed state of Rollout
type RolloutStatus struct {
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

// Rollout is a wrapper type of github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1.Rollout type.
type Rollout struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RolloutSpec   `json:"spec"`
	Status RolloutStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RolloutList contains a list of Rollout
type RolloutList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rollout `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Rollout{}, &RolloutList{})
}
