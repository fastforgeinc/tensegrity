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
	"reconciler.io/runtime/apis"
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
	// Kubernetes status.
	apis.Status `json:",inline"`
	// Tensegrity status.
	v1alpha1.TensegrityStatus `json:",inline"`
}

// Deployment is a wrapper type of the k8s.io/api/apps/v1.Deployment type.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Deployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeploymentSpec   `json:"spec"`
	Status DeploymentStatus `json:"status,omitempty"`
}

// DeploymentList contains a list of Deployment.
// +kubebuilder:object:root=true
type DeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Deployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Deployment{}, &DeploymentList{})
}
