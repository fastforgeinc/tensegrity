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
