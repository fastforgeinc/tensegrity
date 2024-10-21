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
