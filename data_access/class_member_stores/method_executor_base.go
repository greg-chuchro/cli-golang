package class_member_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/data_access"
	"github.com/greg-chuchro/cli-golang/data_access/extensions"
)

// MethodExecutorBase executes methods with parameter assignment and invocation support.
type MethodExecutorBase struct {
	*ParameterStoreBase
	TargetObject any
	Arguments    []any
}

func NewMethodExecutorBase(method reflect.Value, targetObject any) *MethodExecutorBase {
	b := &MethodExecutorBase{}
	b.ParameterStoreBase = NewParameterStoreBase(method)
	b.TargetObject = targetObject
	b.Arguments = []any{}
	return b
}

func (this *MethodExecutorBase) GetByAccessor(accessor int) any {
	result := this.ParameterStoreBase.GetByAccessor(accessor)
	if result == paramUnassigned {
		return nil
	}
	return result
}

func (this *MethodExecutorBase) AddArgument(value any) {
	this.Arguments = append(this.Arguments, value)
}

func (this *MethodExecutorBase) Invoke() (any, error) {
	in := make([]reflect.Value, len(this.ParameterValues))
	for i, v := range this.ParameterValues {
		if IsParamUnassigned(v) {
			in[i] = reflect.Zero(this.ParameterInfos[i])
		} else {
			in[i] = reflect.ValueOf(v)
		}
	}

	result := this.ParentMethod.Call(in)
	if len(result) > 0 {
		return result[0].Interface(), nil
	}
	return nil, nil
}

var _ data_access.IFunctionExecutor[string, any] = (*MethodExecutorBase)(nil)
var _ = extensions.CreateInstance
