package class_member_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/data_access"
	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
	"github.com/greg-chuchro/cli-golang/data_access/extensions"
)

// PropertyStoreBase provides key-value access to object properties by name.
type PropertyStoreBase struct {
	data_access.KeyValueStoreBase[string, any, string, any]
	ClassDataMemberStoreBase[string, any, string, any]
}

func NewPropertyStoreBase(targetObject any, bindingFlags BindingFlags, accessValidator data_access.IAccessValidator) *PropertyStoreBase {
	b := &PropertyStoreBase{}
	b.TargetObject = targetObject
	b.BindingFlags = bindingFlags
	b.AccessValidator = accessValidator
	b.TargetType = reflect.TypeOf(targetObject).Elem()
	b.TargetValue = reflect.ValueOf(targetObject).Elem()
	b.KeyValueStoreBase = *data_access.NewKeyValueStoreBase[string, any, string, any](b)
	return b
}

func (this *PropertyStoreBase) Accessors() []string {
	var result []string
	methodNames := this.propertyMethodNames()
	for i := 0; i < this.TargetType.NumField(); i++ {
		field := this.TargetType.Field(i)
		// A CLI "property" mirrors a C# property: a field X that also exposes
		// GetX()/SetX() accessor methods (C# GetProperties returns only such members).
		if methodNames["Get"+field.Name] && methodNames["Set"+field.Name] && this.ContainsAccessor(field.Name) {
			result = append(result, field.Name)
		}
	}
	return result
}

// propertyMethodNames returns the set of method names available on the target type,
// including pointer-receiver methods (where Go property accessors are typically declared).
func (this *PropertyStoreBase) propertyMethodNames() map[string]bool {
	names := map[string]bool{}
	t := this.TargetType
	if t.Kind() != reflect.Ptr {
		for i := 0; i < t.NumMethod(); i++ {
			names[t.Method(i).Name] = true
		}
		ptr := reflect.PointerTo(t)
		for i := 0; i < ptr.NumMethod(); i++ {
			names[ptr.Method(i).Name] = true
		}
	} else {
		for i := 0; i < t.NumMethod(); i++ {
			names[t.Method(i).Name] = true
		}
	}
	return names
}

func (this *PropertyStoreBase) GetByAccessor(accessor string) any {
	return this.TargetValue.FieldByName(accessor).Interface()
}

func (this *PropertyStoreBase) SetByAccessor(accessor string, value any) {
	this.TargetValue.FieldByName(accessor).Set(reflect.ValueOf(value))
}

func (this *PropertyStoreBase) ContainsAccessor(accessor string) bool {
	field, ok := this.TargetType.FieldByName(accessor)
	if !ok {
		return false
	}
	// Mirrors C# ClassDataMemberStoreBase.ContainsAccessor, which validates the member.
	return this.AccessValidator == nil || this.AccessValidator.IsValid(field)
}

func (this *PropertyStoreBase) GetAccessor(key string) string {
	accessor, ok := this.TryGetAccessor(key)
	if !ok {
		panic(exceptions.KeyNotFound(key))
	}
	return accessor
}

func (this *PropertyStoreBase) GetValueTypeByAccessor(accessor string) reflect.Type {
	field, _ := this.TargetType.FieldByName(accessor)
	return field.Type
}

func (this *PropertyStoreBase) GetValueOrInitializeByAccessor(accessor string) any {
	field := this.TargetValue.FieldByName(accessor)
	value := field.Interface()
	if isNil(value) {
		value = extensions.CreateInstance(field.Type())
		field.Set(reflect.ValueOf(value))
	}
	return value
}

func (this *PropertyStoreBase) TryGetAccessor(key string) (string, bool) {
	if _, ok := this.TargetType.FieldByName(key); ok {
		return key, true
	}
	return "", false
}

func (this *PropertyStoreBase) ConvertFromInternalValue(value any, accessor string) any {
	return value
}

func (this *PropertyStoreBase) ConvertToInternalValue(value any, accessor string) any {
	return value
}
