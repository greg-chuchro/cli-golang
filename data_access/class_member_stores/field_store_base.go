package class_member_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/data_access"
	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
	"github.com/greg-chuchro/cli-golang/data_access/extensions"
)

// FieldStoreBase provides key-value access to object fields by name.
type FieldStoreBase struct {
	data_access.KeyValueStoreBase[string, any, string, any]
	ClassDataMemberStoreBase[string, any, string, any]
}

func NewFieldStoreBase(targetObject any, bindingFlags BindingFlags, accessValidator data_access.IAccessValidator) *FieldStoreBase {
	b := &FieldStoreBase{}
	b.TargetObject = targetObject
	b.BindingFlags = bindingFlags
	b.AccessValidator = accessValidator
	b.TargetType = reflect.TypeOf(targetObject).Elem()
	b.TargetValue = reflect.ValueOf(targetObject).Elem()
	b.KeyValueStoreBase = *data_access.NewKeyValueStoreBase[string, any, string, any](b)
	return b
}

func (this *FieldStoreBase) Accessors() []string {
	var result []string
	for i := 0; i < this.TargetType.NumField(); i++ {
		field := this.TargetType.Field(i)
		if this.ContainsAccessor(field.Name) {
			result = append(result, field.Name)
		}
	}
	return result
}

func (this *FieldStoreBase) GetByAccessor(accessor string) any {
	return this.TargetValue.FieldByName(accessor).Interface()
}

func (this *FieldStoreBase) SetByAccessor(accessor string, value any) {
	this.TargetValue.FieldByName(accessor).Set(reflect.ValueOf(value))
}

func (this *FieldStoreBase) ContainsAccessor(accessor string) bool {
	field, ok := this.TargetType.FieldByName(accessor)
	if !ok {
		return false
	}
	// Mirrors C# ClassDataMemberStoreBase.ContainsAccessor, which validates the member.
	return this.AccessValidator == nil || this.AccessValidator.IsValid(field)
}

func (this *FieldStoreBase) GetAccessor(key string) string {
	accessor, ok := this.TryGetAccessor(key)
	if !ok {
		panic(exceptions.KeyNotFound(key))
	}
	return accessor
}

func (this *FieldStoreBase) GetValueTypeByAccessor(accessor string) reflect.Type {
	field, _ := this.TargetType.FieldByName(accessor)
	return field.Type
}

func (this *FieldStoreBase) GetValueOrInitializeByAccessor(accessor string) any {
	field := this.TargetValue.FieldByName(accessor)
	value := field.Interface()
	if isNil(value) {
		value = extensions.CreateInstance(field.Type())
		field.Set(reflect.ValueOf(value))
	}
	return value
}

func (this *FieldStoreBase) TryGetAccessor(key string) (string, bool) {
	if _, ok := this.TargetType.FieldByName(key); ok {
		return key, true
	}
	return "", false
}

func (this *FieldStoreBase) ConvertFromInternalValue(value any, accessor string) any {
	return value
}

func (this *FieldStoreBase) ConvertToInternalValue(value any, accessor string) any {
	return value
}
