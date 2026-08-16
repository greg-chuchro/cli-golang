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

// CliMethodExecutor provides CLI key-value access to method parameters with name resolution and conversion.
type CliMethodExecutor struct {
	*cms.MethodExecutorBase
	accessorsByName    map[string]int
	bindingFlags       cms.BindingFlags
	stringifier        data_access.IClassMemberStringifier
	compositeConverter data_access.ICompositeValueConverter[any]
}

func NewCliMethodExecutor(method reflect.Value, obj any, bindingFlags cms.BindingFlags, stringifier data_access.IClassMemberStringifier, compositeConverter data_access.ICompositeValueConverter[any]) *CliMethodExecutor {
	base := cms.NewMethodExecutorBase(method, obj)
	s := &CliMethodExecutor{
		MethodExecutorBase: base,
		bindingFlags:       bindingFlags,
		stringifier:        stringifier,
		compositeConverter: compositeConverter,
	}
	s.MethodExecutorBase.KeyValueStoreBase = *da.NewKeyValueStoreBase[string, any, int, any](s)
	s.accessorsByName = s.buildAccessorsByName()
	return s
}

func (this *CliMethodExecutor) buildAccessorsByName() map[string]int {
	cmp := func(a, b string) bool {
		if this.bindingFlags.IgnoreCase {
			return strings.EqualFold(a, b)
		}
		return a == b
	}

	byRequired := map[string]int{}
	for i, t := range this.MethodExecutorBase.ParameterInfos {
		_ = t
		stField := this.synthParamField(i)
		for _, key := range this.stringifier.GetRequiredNames(stField) {
			found := false
			for k := range byRequired {
				if cmp(k, key) {
					found = true
					break
				}
			}
			if !found {
				byRequired[key] = i
			} else {
				panic(exceptions.CliErrors.NameCollision(this.paramName(i), this.paramName(byRequired[key])))
			}
		}
	}

	byAlt := map[string]int{}
	collisions := map[string]bool{}
	for i := range this.MethodExecutorBase.ParameterInfos {
		stField := this.synthParamField(i)
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
				byAlt[key] = i
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

func (this *CliMethodExecutor) paramName(i int) string {
	methodType := this.MethodExecutorBase.ParentMethod.Type()
	if methodType.NumIn() > i {
		return "arg" + itoaInt(i)
	}
	return "arg" + itoaInt(i)
}

func (this *CliMethodExecutor) synthParamField(i int) reflect.StructField {
	t := this.MethodExecutorBase.ParameterInfos[i]
	return reflect.StructField{Name: this.paramName(i), Type: t}
}

func (this *CliMethodExecutor) GetAccessorKeysPairs() []data_access.AccessorKeysPair[reflect.StructField] {
	grouped := map[int]reflect.StructField{}
	for _, i := range this.accessorsByName {
		grouped[i] = this.synthParamField(i)
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

func (this *CliMethodExecutor) IsMultiValue(key string) bool {
	i, ok := this.accessorsByName[key]
	if !ok {
		return false
	}
	return this.compositeConverter.IsMultiValue(this.MethodExecutorBase.ParameterInfos[i])
}

func (this *CliMethodExecutor) ContainsAccessor(accessor int) bool {
	return accessor >= 0 && accessor < len(this.MethodExecutorBase.ParameterInfos)
}

func (this *CliMethodExecutor) GetValueTypeByAccessor(accessor int) reflect.Type {
	if accessor < 0 || accessor >= len(this.MethodExecutorBase.ParameterInfos) {
		return nil
	}
	return this.compositeConverter.GetValueType(this.MethodExecutorBase.ParameterInfos[accessor])
}

func (this *CliMethodExecutor) GetAccessor(key string) int {
	accessor, ok := this.TryGetAccessor(key)
	if !ok {
		panic(dae.KeyNotFound(key))
	}
	return accessor
}

func (this *CliMethodExecutor) TryGetAccessor(key string) (int, bool) {
	i, ok := this.accessorsByName[key]
	return i, ok
}

func (this *CliMethodExecutor) ConvertFromInternalValue(value any, accessor int) any {
	return this.compositeConverter.ConvertFrom(value, this.MethodExecutorBase.ParameterInfos[accessor])
}

func (this *CliMethodExecutor) ConvertToInternalValue(value any, accessor int) any {
	currentValue := this.MethodExecutorBase.GetValueOrInitializeByAccessor(accessor)
	return this.compositeConverter.ConvertToOrUpdate(value, this.MethodExecutorBase.ParameterInfos[accessor], currentValue)
}

func (this *CliMethodExecutor) SetParameterValues() {
	argumentIndex := 0
	for parameterIndex := 0; parameterIndex < len(this.MethodExecutorBase.ParameterValues); parameterIndex++ {
		if cms.IsParamUnassigned(this.MethodExecutorBase.ParameterValues[parameterIndex]) {
			if argumentIndex < len(this.MethodExecutorBase.Arguments) {
				value := this.ConvertToInternalValue(this.MethodExecutorBase.Arguments[argumentIndex], parameterIndex)
				this.MethodExecutorBase.ParameterValues[parameterIndex] = value
				argumentIndex++
			}
		}
	}
	if argumentIndex < len(this.MethodExecutorBase.Arguments) {
		panic(dae.DataAccessErrors.UnexpectedArgument(this.MethodExecutorBase.Arguments[argumentIndex]))
	}
}

func itoaInt(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	buf := make([]byte, 0, 12)
	for i > 0 {
		buf = append([]byte{byte('0' + i%10)}, buf...)
		i /= 10
	}
	if neg {
		buf = append([]byte{'-'}, buf...)
	}
	return string(buf)
}
