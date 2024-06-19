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
