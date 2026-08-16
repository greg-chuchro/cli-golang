package class_member_stores

import (
	"reflect"
	"strconv"

	"github.com/greg-chuchro/cli-golang/data_access"
	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
	"github.com/greg-chuchro/cli-golang/data_access/extensions"
)

// paramUnassigned is the sentinel used in place of C#'s DBNull.Value to mark unassigned parameters.
var paramUnassigned = struct{}{}

// ParameterStoreBase provides key-value access to method parameters by index.
type ParameterStoreBase struct {
	data_access.KeyValueStoreBase[string, any, int, any]
	ParameterInfos            []reflect.Type
	ParameterValues           []any
	ParameterIndexByParameter map[int]int
	ParentMethod              reflect.Value
}

func NewParameterStoreBase(method reflect.Value) *ParameterStoreBase {
	b := &ParameterStoreBase{}
	b.ParentMethod = method
	methodType := method.Type()
	count := methodType.NumIn()
	b.ParameterInfos = make([]reflect.Type, count)
	b.ParameterValues = make([]any, count)
	b.ParameterIndexByParameter = make(map[int]int, count)
	for i := 0; i < count; i++ {
		b.ParameterInfos[i] = methodType.In(i)
		b.ParameterIndexByParameter[i] = i
		b.ParameterValues[i] = paramUnassigned
	}
	b.KeyValueStoreBase = *data_access.NewKeyValueStoreBase[string, any, int, any](b)
	return b
}

func (this *ParameterStoreBase) Accessors() []int {
	result := make([]int, len(this.ParameterInfos))
	for i := range this.ParameterInfos {
		result[i] = i
	}
	return result
}

func (this *ParameterStoreBase) GetByAccessor(accessor int) any {
	return this.ParameterValues[accessor]
}

func (this *ParameterStoreBase) SetByAccessor(accessor int, value any) {
	this.ParameterValues[accessor] = value
}

func (this *ParameterStoreBase) ContainsAccessor(accessor int) bool {
	return accessor >= 0 && accessor < len(this.ParameterInfos)
}

func (this *ParameterStoreBase) GetAccessor(key string) int {
	accessor, ok := this.TryGetAccessor(key)
	if !ok {
		panic(exceptions.KeyNotFound(key))
	}
	return accessor
}

func (this *ParameterStoreBase) GetValueTypeByAccessor(accessor int) reflect.Type {
	return this.ParameterInfos[accessor]
}

func (this *ParameterStoreBase) GetValueOrInitializeByAccessor(accessor int) any {
	value := this.ParameterValues[accessor]
	if value == paramUnassigned {
		value = nil
	}
	if value == nil {
		value = extensions.CreateInstance(this.ParameterInfos[accessor])
		this.ParameterValues[accessor] = value
	}
	return value
}

func (this *ParameterStoreBase) TryGetAccessor(key string) (int, bool) {
	for i := range this.ParameterInfos {
		if strconv.Itoa(i) == key {
			return i, true
		}
	}
	return 0, false
}

func (this *ParameterStoreBase) ConvertFromInternalValue(value any, accessor int) any {
	return value
}

func (this *ParameterStoreBase) ConvertToInternalValue(value any, accessor int) any {
	return value
}

// ParamUnassigned returns the sentinel value used in place of C#'s DBNull.Value to mark unassigned parameters.
func ParamUnassigned() any {
	return paramUnassigned
}

// IsParamUnassigned reports whether v is the unassigned-parameter sentinel.
func IsParamUnassigned(v any) bool {
	return v == paramUnassigned
}
