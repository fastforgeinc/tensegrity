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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *Deployment) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-k8s-tensegrity-fastforge-io-v1alpha1-deployment,mutating=true,failurePolicy=fail,sideEffects=None,groups=k8s.tensegrity.fastforge.io,resources=deployments,verbs=create;update,versions=v1alpha1,name=mdeployment.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Deployment{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Deployment) Default() {
	r.Spec.TensegritySpec.SetDefaultProducesName(r.GetName())
	r.Spec.TensegritySpec.SetDefaultNamespaceDelegate(r.GetNamespace())
	r.Spec.TensegritySpec.SetDefaultConsumesConfigMapName(r.GetName() + v1alpha1.DefaultConsumesConfigMapNamePrefix)
	r.Spec.TensegritySpec.SetDefaultConsumesSecretName(r.GetName() + v1alpha1.DefaultConsumesSecretNamePrefix)
	r.Spec.TensegritySpec.SetDefaultProducesConfigMapName(r.GetName() + v1alpha1.DefaultProducesConfigMapNamePrefix)
	r.Spec.TensegritySpec.SetDefaultProducesSecretName(r.GetName() + v1alpha1.DefaultProducesSecretNamePrefix)
}

//+kubebuilder:webhook:path=/validate-k8s-tensegrity-fastforge-io-v1alpha1-deployment,mutating=false,failurePolicy=fail,sideEffects=None,groups=k8s.tensegrity.fastforge.io,resources=deployments,verbs=create;update,versions=v1alpha1,name=vdeployment.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Deployment{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Deployment) ValidateCreate() (admission.Warnings, error) {
	if errs := r.Spec.TensegritySpec.Validate(); len(errs) > 0 {
		return nil, apierrors.NewInvalid(r.GetObjectKind().GroupVersionKind().GroupKind(), r.GetName(), errs)
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Deployment) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	if errs := r.Spec.TensegritySpec.Validate(); len(errs) > 0 {
		return nil, apierrors.NewInvalid(r.GetObjectKind().GroupVersionKind().GroupKind(), r.GetName(), errs)
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Deployment) ValidateDelete() (admission.Warnings, error) {
	if errs := r.Spec.TensegritySpec.Validate(); len(errs) > 0 {
		return nil, apierrors.NewInvalid(r.GetObjectKind().GroupVersionKind().GroupKind(), r.GetName(), errs)
	}
	return nil, nil
}
