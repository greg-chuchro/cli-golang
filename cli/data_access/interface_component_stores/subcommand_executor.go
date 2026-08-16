package interface_component_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/cli/data_access"
	"github.com/greg-chuchro/cli-golang/cli/interface_components"
)

// SubcommandExecutor executes CLI subcommands with parameter handling and metadata.
type SubcommandExecutor struct {
	executor data_access.ICliFunctionExecutor[any]
}

func NewSubcommandExecutor(executor data_access.ICliFunctionExecutor[any]) *SubcommandExecutor {
	return &SubcommandExecutor{executor: executor}
}

func (this *SubcommandExecutor) Get(key string) any {
	return this.executor.Get(key)
}

func (this *SubcommandExecutor) Set(key string, value any) {
	this.executor.Set(key, value)
}

func (this *SubcommandExecutor) GetValueOrInitialize(key string) any {
	return this.executor.GetValueOrInitialize(key)
}

func (this *SubcommandExecutor) ContainsKey(key string) bool {
	return this.executor.ContainsKey(key)
}

func (this *SubcommandExecutor) GetValueType(key string) reflect.Type {
	return this.executor.GetValueType(key)
}

func (this *SubcommandExecutor) SetParameterValues() {
	this.executor.SetParameterValues()
}

func (this *SubcommandExecutor) AddArgument(value any) {
	this.executor.AddArgument(value)
}

func (this *SubcommandExecutor) Invoke() (any, error) {
	return this.executor.Invoke()
}

func (this *SubcommandExecutor) GetParameters() []interface_components.Parameter {
	result := make([]interface_components.Parameter, 0)
	for _, pair := range this.executor.GetAccessorKeysPairs() {
		key := pair.Keys[0]
		result = append(result, interface_components.Parameter{
			HasDefaultValue: false,
			IsMultiValue:    this.executor.IsMultiValue(key),
			Keys:            pair.Keys,
			ParameterInfo:   pair.Accessor,
			ValueType:       this.GetValueType(key),
			Value:           this.toString(this.Get(key)),
		})
	}
	return result
}

func (this *SubcommandExecutor) GetFirstUnassignedParameter() *interface_components.Parameter {
	this.executor.SetParameterValues()
	for _, p := range this.GetParameters() {
		if p.Value == "" {
			return &p
		}
	}
	return nil
}

func (this *SubcommandExecutor) toString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return reflect.ValueOf(v).String()
}
