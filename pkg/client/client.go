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

package client

import (
	k8srest "k8s.io/client-go/rest"
	"reconciler.io/runtime/duck"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func New(config *k8srest.Config, options k8sclient.Options) (k8sclient.Client, error) {
	client, err := NewRetryClient(config, options)
	if err != nil {
		return nil, err
	}
	return duck.NewDuckAwareClientWrapper(client), nil
}
