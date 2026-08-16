package completion_providers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/greg-chuchro/cli-golang/cli_attributes/exceptions"
)

type MethodCompletionProvider struct{}

func NewMethodCompletionProvider() *MethodCompletionProvider {
	return &MethodCompletionProvider{}
}

func (p *MethodCompletionProvider) GetCompletions(ctx context.Context, completionContext ICompletionProviderContext) []string {
	filter := completionContext.Filter()
	if filter == "" {
		return []string{}
	}

	submodule := completionContext.Submodule()
	if submodule == nil {
		return []string{}
	}

	submoduleValue := reflect.ValueOf(submodule)
	method := submoduleValue.MethodByName(filter)
	if !method.IsValid() {
		typeName := reflect.TypeOf(submodule).Name()
		panic(exceptions.NewCliAttributeException(
			fmt.Sprintf("Completion method '%s' not found on type '%s'", filter, typeName),
		))
	}

	methodType := method.Type()
	if methodType.NumIn() != 1 || methodType.In(0).Kind() != reflect.String {
		typeName := reflect.TypeOf(submodule).Name()
		panic(exceptions.NewCliAttributeException(
			fmt.Sprintf("Completion method '%s' on type '%s' must have signature: []string %s(string partialInput)", filter, typeName, filter),
		))
	}

	partialInput := completionContext.PartialInput()
	result := method.Call([]reflect.Value{reflect.ValueOf(partialInput)})

	if result[0].Len() == 0 {
		return []string{}
	}

	completions := make([]string, result[0].Len())
	for i := 0; i < result[0].Len(); i++ {
		completions[i] = result[0].Index(i).String()
	}

	return completions
}
