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

// WorkloadName is a name of dependency Deployment or StatefulSet to consume keys from
type WorkloadName string

// Delegate is identifier of namespace used to resolve consumed keys from.
type Delegate string

// Key is workload configuration entity that can be produced and/or consumed by other workloads.
type Key string

// KeyMap is mapping between a Deployment environment variable and consumed key.
type KeyMap map[string]string

// ProduceFromType is a string enumeration of values used by ProduceSpec to produces values of keys.
type ProduceFromType string

const (
	// ProduceFromKind is a type of ProduceSpec to use Kubernetes resources to produce value of a key.
	ProduceFromKind ProduceFromType = "kind"
	// ProduceFromValue is a type of ProduceSpec to use static value of a key declared in this spec file.
	ProduceFromValue ProduceFromType = "value"
)

type ConsumeSpec struct {
	// Maps defines mappings between consumed keys and environment variables.
	Maps KeyMap `json:"maps,omitempty"`
}

type ProduceFromKindSpec struct {
	// Kind of a resource to produce a key from.
	Kind *string `json:"kind,omitempty"`
	// Expression is yq filter to get value from a Kind resource
	// (docs: https://mikefarah.gitbook.io/yq)
	Expression *string `json:"expression,omitempty"`
}

type ProduceFromValueSpec struct {
	// Value is static value defined is spec file.
	Value string `json:"value,omitempty"`
}

type ProduceSpec struct {
	// From defines source of values of produced keys.
	From ProduceFromType `json:"from,omitempty"`
	// ProduceFromKindSpec defines a Kubernetes resource as a source of value for produced key.
	ProduceFromKindSpec `json:",inline"`
	// ProduceFromValueSpec defines a static value for produced key.
	ProduceFromValueSpec `json:",inline"`
}

type WorkloadSpec struct {
	// Delegates is a slice of Delegate type used to resolve consumed keys by a Deployment in specified order.
	Delegates []Delegate `json:"delegates,omitempty"`
	// Consumes is a map of other workloads and ConsumeSpec.
	Consumes map[WorkloadName]ConsumeSpec `json:"consumes,omitempty"`
	// Produces is a map of keys and value sources to get from.
	Produces map[Key]ProduceSpec `json:"produces,omitempty"`
}
