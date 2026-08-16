package data_access

import (
	"reflect"

	da "github.com/greg-chuchro/cli-golang/data_access"
	"github.com/greg-chuchro/cli-golang/data_access/extensions"
)

// IAccessValidator is re-exported from the base data_access package.
type IAccessValidator = da.IAccessValidator

// OptionAccessValidator validates whether a member is a valid CLI option (must be convertible type).
type OptionAccessValidator struct {
	valueConverter IValueConverter[any]
}

func NewOptionAccessValidator(valueConverter IValueConverter[any]) *OptionAccessValidator {
	return &OptionAccessValidator{valueConverter: valueConverter}
}

func (this *OptionAccessValidator) IsValid(accessor reflect.StructField) bool {
	t := extensions.GetUnderlyingType(accessor)
	return this.valueConverter.CanConvert(t)
}

// SubcommandAccessValidator validates whether a method is a valid CLI subcommand (excludes Go-specific methods).
type SubcommandAccessValidator struct{}

func NewSubcommandAccessValidator() *SubcommandAccessValidator {
	return &SubcommandAccessValidator{}
}

func (this *SubcommandAccessValidator) IsValid(accessor reflect.StructField) bool {
	return true
}

// SubmoduleAccessValidator validates whether a member is a valid CLI submodule (must be non-convertible type).
type SubmoduleAccessValidator struct {
	valueConverter IValueConverter[any]
}

func NewSubmoduleAccessValidator(valueConverter IValueConverter[any]) *SubmoduleAccessValidator {
	return &SubmoduleAccessValidator{valueConverter: valueConverter}
}

func (this *SubmoduleAccessValidator) IsValid(accessor reflect.StructField) bool {
	t := extensions.GetUnderlyingType(accessor)
	return !this.valueConverter.CanConvert(t)
}
