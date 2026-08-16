package class_member_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/data_access"
)

// MethodExecutor executes methods with parameter assignment and invocation support.
type MethodExecutor struct {
	*MethodExecutorBase
}

func NewMethodExecutor(method reflect.Value, obj any) *MethodExecutor {
	return &MethodExecutor{MethodExecutorBase: NewMethodExecutorBase(method, obj)}
}

var _ data_access.IFunctionExecutor[string, any] = (*MethodExecutor)(nil)
