package data_access

import "reflect"

// IAccessValidator validates whether a class member should be accessible in the CLI.
// Mirrors CalqFramework.DataAccess.IAccessValidator.
type IAccessValidator interface {
	IsValid(accessor reflect.StructField) bool
}
