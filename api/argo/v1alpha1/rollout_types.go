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
	"encoding/json"
	rolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	"github.com/fastforgeinc/tensegrity/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"reconciler.io/runtime/apis"
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
	// Kubernetes status.
	apis.Status `json:",inline"`
	// Tensegrity status.
	v1alpha1.TensegrityStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

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
