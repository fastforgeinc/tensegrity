package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
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
	seenDelegates := make(map[v1.ObjectReference]struct{}, len(s.Delegates))
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
