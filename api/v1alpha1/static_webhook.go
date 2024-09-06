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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *Static) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-tensegrity-fastforge-io-v1alpha1-static,mutating=true,failurePolicy=fail,sideEffects=None,groups=tensegrity.fastforge.io,resources=statics,verbs=create;update,versions=v1alpha1,name=mstatic.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Static{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Static) Default() {
	r.Spec.TensegritySpec.SetDefaultProducesName(r.GetName())
	r.Spec.TensegritySpec.SetDefaultNamespaceDelegate(r.GetNamespace())
	r.Spec.TensegritySpec.SetDefaultProducesConfigMapName(r.GetName() + DefaultProducesConfigMapNamePrefix)
	r.Spec.TensegritySpec.SetDefaultProducesSecretName(r.GetName() + DefaultProducesSecretNamePrefix)
}

//+kubebuilder:webhook:path=/validate-tensegrity-fastforge-io-v1alpha1-static,mutating=false,failurePolicy=fail,sideEffects=None,groups=tensegrity.fastforge.io,resources=statics,verbs=create;update,versions=v1alpha1,name=vstatic.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Static{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Static) ValidateCreate() (admission.Warnings, error) {
	if errs := r.Spec.TensegritySpec.Validate(); len(errs) > 0 {
		return nil, apierrors.NewInvalid(r.GetObjectKind().GroupVersionKind().GroupKind(), r.GetName(), errs)
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Static) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	if errs := r.Spec.TensegritySpec.Validate(); len(errs) > 0 {
		return nil, apierrors.NewInvalid(r.GetObjectKind().GroupVersionKind().GroupKind(), r.GetName(), errs)
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Static) ValidateDelete() (admission.Warnings, error) {
	if errs := r.Spec.TensegritySpec.Validate(); len(errs) > 0 {
		return nil, apierrors.NewInvalid(r.GetObjectKind().GroupVersionKind().GroupKind(), r.GetName(), errs)
	}
	return nil, nil
}
