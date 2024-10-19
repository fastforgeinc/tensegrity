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
	"context"
	k8srest "k8s.io/client-go/rest"
	k8sretry "k8s.io/client-go/util/retry"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type retryClient struct {
	k8sclient.Client
}

var _ k8sclient.Client = &retryClient{}

func NewRetryClient(config *k8srest.Config, options k8sclient.Options) (k8sclient.Client, error) {
	client, err := k8sclient.New(config, options)
	if err != nil {
		return nil, err
	}
	return &retryClient{client}, nil
}

func (c *retryClient) Create(ctx context.Context, obj k8sclient.Object, opts ...k8sclient.CreateOption) error {
	return c.retryOnConflict(c.Client.Create(ctx, obj, opts...))
}

func (c *retryClient) Delete(ctx context.Context, obj k8sclient.Object, opts ...k8sclient.DeleteOption) error {
	return c.retryOnConflict(c.Client.Delete(ctx, obj, opts...))
}

func (c *retryClient) Update(ctx context.Context, obj k8sclient.Object, opts ...k8sclient.UpdateOption) error {
	return c.retryOnConflict(c.Client.Update(ctx, obj, opts...))
}

func (c *retryClient) Patch(
	ctx context.Context, obj k8sclient.Object, patch k8sclient.Patch, opts ...k8sclient.PatchOption) error {

	return c.retryOnConflict(c.Client.Patch(ctx, obj, patch, opts...))
}

func (c *retryClient) DeleteAllOf(
	ctx context.Context, obj k8sclient.Object, opts ...k8sclient.DeleteAllOfOption) error {

	return c.retryOnConflict(c.Client.DeleteAllOf(ctx, obj, opts...))
}

func (c *retryClient) retryOnConflict(err error) error {
	return k8sretry.RetryOnConflict(k8sretry.DefaultRetry, func() error {
		return err
	})
}
