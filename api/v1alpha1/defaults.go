package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
)

func (s *TensegritySpec) SetDefaultConsumesSecretName(name string) {
	if len(s.ConsumesSecretName) == 0 {
		s.ConsumesSecretName = name
	}
}

func (s *TensegritySpec) SetDefaultConsumesConfigMapName(name string) {
	if len(s.ConsumesConfigMapName) == 0 {
		s.ConsumesConfigMapName = name
	}
}

func (s *TensegritySpec) SetDefaultProducesSecretName(name string) {
	if len(s.ProducesSecretName) == 0 {
		s.ProducesSecretName = name
	}
}

func (s *TensegritySpec) SetDefaultProducesConfigMapName(name string) {
	if len(s.ProducesConfigMapName) == 0 {
		s.ProducesConfigMapName = name
	}
}

func (s *TensegritySpec) SetDefaultProducesName(name string) {
	for i, produce := range s.Produces {
		if len(produce.Name) == 0 {
			s.Produces[i].Name = name
		}
	}
}

func (s *TensegritySpec) SetDefaultNamespaceDelegate(namespace string) {
	if len(s.Delegates) == 0 {
		s.Delegates = append(s.Delegates, v1.ObjectReference{
			Kind: "Namespace",
			Name: namespace,
		})
	}
}
