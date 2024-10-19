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
