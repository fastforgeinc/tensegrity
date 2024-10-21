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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/util/jsonpath"
)

func (s *TensegritySpec) Validate() (allErrs field.ErrorList) {
	if errs := s.validateProduces(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := s.validateConsumes(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := s.validateDelegates(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	return
}

func (s *TensegritySpec) validateConsumes() (errs field.ErrorList) {
	seenEnvs := make(map[string]struct{})
	seenRefs := make(map[corev1.ObjectReference]struct{})
	for i, c := range s.Consumes {
		if len(c.APIVersion) == 0 {
			errs = append(errs, field.Required(
				field.NewPath("spec").Child("consumes").Index(i).Key("apiVersion"), "valid resource api version"))
		}
		if len(c.Kind) == 0 {
			errs = append(errs, field.Required(
				field.NewPath("spec").Child("consumes").Index(i).Child("kind"), "valid resource kind"))
		}
		if len(c.Name) == 0 {
			errs = append(errs, field.Required(
				field.NewPath("spec").Child("consumes").Index(i).Child("name"), "valid resource name"))
		}
		if len(c.Maps) == 0 {
			errs = append(errs, field.Required(
				field.NewPath("spec").Child("consumes").Index(i).Child("maps"),
				"valid environment variables to keys mapping"))
		}
		if _, ok := seenRefs[c.ObjectReference]; ok {
			errs = append(errs, field.Duplicate(
				field.NewPath("spec").Child("consumes").Index(i), c.ObjectReference))
		}
		seenRefs[c.ObjectReference] = struct{}{}
		for env := range c.Maps {
			if _, ok := seenEnvs[env]; ok {
				errs = append(errs, field.Duplicate(
					field.NewPath("spec").Child("consumes").Index(i).Child("maps"), env))
			}
			seenEnvs[env] = struct{}{}
		}
	}
	return errs
}

func (s *TensegritySpec) validateProduces() (errs field.ErrorList) {
	seenKeys := make(map[string]struct{}, len(s.Produces))
	for i, p := range s.Produces {
		if len(p.Key) == 0 {
			errs = append(errs, field.Required(
				field.NewPath("spec").Child("produces").Index(i).Child("key"), "valid key name"))
		}
		if _, ok := seenKeys[p.Key]; ok {
			errs = append(errs, field.Duplicate(
				field.NewPath("spec").Child("produces").Index(i).Child("key"), p.Key))
		}
		seenKeys[p.Key] = struct{}{}
		if len(p.APIVersion) == 0 {
			errs = append(errs, field.Required(
				field.NewPath("spec").Child("produces").Index(i).Child("apiVersion"), "valid resource api version"))
		}
		if len(p.Kind) == 0 {
			errs = append(errs, field.Required(
				field.NewPath("spec").Child("produces").Index(i).Child("kind"), "valid resource kind"))
		}
		if len(p.Name) == 0 {
			errs = append(errs, field.Required(
				field.NewPath("spec").Child("produces").Index(i).Child("name"), "valid resource name"))
		}
		if len(p.FieldPath) == 0 {
			errs = append(errs, field.Required(
				field.NewPath("spec").Child("produces").Index(i).Child("fieldPath"), "valid resource JSONPath"))
		}
		if !p.Sensitive && p.Encoded {
			errs = append(errs, field.Invalid(
				field.NewPath("spec").Child("produces").Index(i).Child("encoded"),
				p.FieldPath, "encoded field is allowed only when key is sensitive"))
		}
		jp := jsonpath.New(p.Key)
		jp.AllowMissingKeys(false)
		if err := jp.Parse(p.FieldPath); err != nil {
			errs = append(errs, field.Invalid(
				field.NewPath("spec").Child("produces").Index(i).Child("fieldPath"),
				p.FieldPath, "valid resource JSONPath"))
		}
	}
	return errs
}

func (s *TensegritySpec) validateDelegates() (errs field.ErrorList) {
	seenDelegates := make(map[corev1.ObjectReference]struct{}, len(s.Delegates))
	for i, d := range s.Delegates {
		if len(d.Kind) == 0 {
			errs = append(errs, field.Required(
				field.NewPath("spec").Child("delegates").Index(i).Child("kind"), "valid resource kind"))
		}
		if len(d.Name) == 0 {
			errs = append(errs, field.Required(
				field.NewPath("spec").Child("delegates").Index(i).Child("name"), "valid resource name"))
		}
		switch d.Kind {
		case "Namespace":
		default:
			errs = append(errs, field.Invalid(
				field.NewPath("spec").Child("delegates").Index(i).Child("kind"),
				d.Kind, "kind must be on of these values: Namespace"))
		}
		if _, ok := seenDelegates[d]; ok {
			errs = append(errs, field.Duplicate(
				field.NewPath("spec").Child("delegates").Index(i), d))
		}
		seenDelegates[d] = struct{}{}
	}
	return
}
