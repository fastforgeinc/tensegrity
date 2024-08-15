package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
)

func (s *TensegritySpec) SetDefaultSecretName(name string) {
	if len(s.SecretName) == 0 {
		s.SecretName = name
	}
}

func (s *TensegritySpec) SetDefaultConfigMapName(name string) {
	if len(s.ConfigMapName) == 0 {
		s.ConfigMapName = name
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
