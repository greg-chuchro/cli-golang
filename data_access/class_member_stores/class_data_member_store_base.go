package class_member_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/data_access"
)

// BindingFlags controls the lookup scope for class members, mirroring System.Reflection.BindingFlags.
type BindingFlags struct {
	Instance   bool
	Static     bool
	Public     bool
	Private    bool
	IgnoreCase bool
}

// DefaultLookup is the default binding flags used by the factory.
var DefaultLookup = BindingFlags{Instance: true, Static: true, Public: true}

// ClassDataMemberStoreBase is the base for field/property stores that access members of a target object.
type ClassDataMemberStoreBase[TKey comparable, TValue any, TAccessor comparable, TInternalValue any] struct {
	TargetObject    any
	BindingFlags    BindingFlags
	AccessValidator data_access.IAccessValidator
	TargetType      reflect.Type
	TargetValue     reflect.Value
}

// isNil reports whether value is nil or a nil reference/collection.
func isNil(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}
