package extensions

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
)

// GetUnderlyingType returns the type of a struct field, mirroring MemberInfo.GetUnderlyingType for FieldInfo.
func GetUnderlyingType(field reflect.StructField) reflect.Type {
	return field.Type
}

// GetUnderlyingTypeOfMethod returns the return type of a method, mirroring MemberInfo.GetUnderlyingType for MethodInfo.
func GetUnderlyingTypeOfMethod(method reflect.Method) reflect.Type {
	if method.Type.NumOut() > 0 {
		return method.Type.Out(0)
	}
	return reflect.TypeOf((*struct{})(nil)).Elem()
}

// GetUnderlyingTypeOfValue returns the type of a reflect.Value, covering field/property/method/event-like values.
func GetUnderlyingTypeOfValue(value reflect.Value) reflect.Type {
	if !value.IsValid() {
		panic(exceptions.DataAccessErrors.UnrecognizedMemberType())
	}
	return value.Type()
}
