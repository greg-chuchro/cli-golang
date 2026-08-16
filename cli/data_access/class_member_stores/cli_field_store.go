package class_member_stores

import (
	"reflect"
	"strings"

	"github.com/greg-chuchro/cli-golang/cli/data_access"
	"github.com/greg-chuchro/cli-golang/cli/exceptions"
	da "github.com/greg-chuchro/cli-golang/data_access"
	cms "github.com/greg-chuchro/cli-golang/data_access/class_member_stores"
	dae "github.com/greg-chuchro/cli-golang/data_access/exceptions"
)

// CliFieldStore provides CLI key-value access to object fields with name resolution and conversion.
type CliFieldStore struct {
	*cms.FieldStoreBase
	accessorsByName    map[string]reflect.StructField
	accessValidator    data_access.IAccessValidator
	stringifier        data_access.IClassMemberStringifier
	compositeConverter data_access.ICompositeValueConverter[any]
}

func NewCliFieldStore(targetObject any, bindingFlags cms.BindingFlags, stringifier data_access.IClassMemberStringifier, accessValidator data_access.IAccessValidator, compositeConverter data_access.ICompositeValueConverter[any]) *CliFieldStore {
	base := cms.NewFieldStoreBase(targetObject, bindingFlags, accessValidator)
	s := &CliFieldStore{
		FieldStoreBase:     base,
		accessValidator:    accessValidator,
		stringifier:        stringifier,
		compositeConverter: compositeConverter,
	}
	s.FieldStoreBase.KeyValueStoreBase = *da.NewKeyValueStoreBase[string, any, string, any](s)
	s.accessorsByName = s.buildAccessorsByName()
	return s
}

func (this *CliFieldStore) buildAccessorsByName() map[string]reflect.StructField {
	cmp := func(a, b string) bool {
		if this.FieldStoreBase.BindingFlags.IgnoreCase {
			return strings.EqualFold(a, b)
		}
		return a == b
	}

	byRequired := map[string]reflect.StructField{}
	for _, name := range this.FieldStoreBase.Accessors() {
		stField, _ := this.FieldStoreBase.TargetType.FieldByName(name)
		for _, key := range this.stringifier.GetRequiredNames(stField) {
			found := false
			for k := range byRequired {
				if cmp(k, key) {
					found = true
					break
				}
			}
			if !found {
				byRequired[key] = stField
			} else {
				panic(exceptions.CliErrors.NameCollision(name, byRequired[key].Name))
			}
		}
	}

	byAlt := map[string]reflect.StructField{}
	collisions := map[string]bool{}
	for _, name := range this.FieldStoreBase.Accessors() {
		stField, _ := this.FieldStoreBase.TargetType.FieldByName(name)
		for _, key := range this.stringifier.GetAlternativeNames(stField) {
			found := false
			for k := range byAlt {
				if cmp(k, key) {
					found = true
					break
				}
			}
			if found {
				collisions[key] = true
			} else {
				byAlt[key] = stField
			}
		}
	}
	for key := range collisions {
		delete(byAlt, key)
	}

	for key, v := range byAlt {
		exists := false
		for k := range byRequired {
			if cmp(k, key) {
				exists = true
				break
			}
		}
		if !exists {
			byRequired[key] = v
		}
	}
	return byRequired
}

func (this *CliFieldStore) GetAccessorKeysPairs() []data_access.AccessorKeysPair[reflect.StructField] {
	grouped := map[string]reflect.StructField{}
	for _, v := range this.accessorsByName {
		grouped[v.Name] = v
	}
	result := make([]data_access.AccessorKeysPair[reflect.StructField], 0, len(grouped))
	for _, field := range grouped {
		keys := []string{}
		keys = append(keys, this.stringifier.GetRequiredNames(field)...)
		keys = append(keys, this.stringifier.GetAlternativeNames(field)...)
		result = append(result, *data_access.NewAccessorKeysPair(field, keys))
	}
	return result
}

func (this *CliFieldStore) IsMultiValue(key string) bool {
	accessor, ok := this.accessorsByName[key]
	if !ok {
		return false
	}
	return this.compositeConverter.IsMultiValue(accessor.Type)
}

func (this *CliFieldStore) ContainsAccessor(accessor string) bool {
	stField, ok := this.FieldStoreBase.TargetType.FieldByName(accessor)
	if !ok {
		return false
	}
	return this.accessValidator.IsValid(stField)
}

func (this *CliFieldStore) GetValueType(accessor string) reflect.Type {
	stField, _ := this.FieldStoreBase.TargetType.FieldByName(accessor)
	return this.compositeConverter.GetValueType(stField.Type)
}

func (this *CliFieldStore) TryGetAccessor(key string) (string, bool) {
	field, ok := this.accessorsByName[key]
	if !ok {
		return "", false
	}
	return field.Name, true
}

func (this *CliFieldStore) GetAccessor(key string) string {
	field, ok := this.accessorsByName[key]
	if !ok {
		panic(dae.KeyNotFound(key))
	}
	return field.Name
}

func (this *CliFieldStore) ConvertFromInternalValue(value any, accessor string) any {
	stField, _ := this.FieldStoreBase.TargetType.FieldByName(accessor)
	return this.compositeConverter.ConvertFrom(value, stField.Type)
}

func (this *CliFieldStore) ConvertToInternalValue(value any, accessor string) any {
	stField, _ := this.FieldStoreBase.TargetType.FieldByName(accessor)
	currentValue := this.FieldStoreBase.GetValueOrInitializeByAccessor(accessor)
	return this.compositeConverter.ConvertToOrUpdate(value, stField.Type, currentValue)
}
