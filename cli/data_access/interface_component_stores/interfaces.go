package interface_component_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/cli/interface_components"
	da "github.com/greg-chuchro/cli-golang/data_access"
)

// IOptionStore provides access to CLI options (fields and properties) with string values.
type IOptionStore interface {
	da.IKeyValueStore[string, any]
	GetOptions() []interface_components.Option
}

// ISubcommandExecutor executes CLI subcommands (methods) with parameter handling.
type ISubcommandExecutor interface {
	da.IFunctionExecutor[string, any]
	GetParameters() []interface_components.Parameter
	GetFirstUnassignedParameter() *interface_components.Parameter
	SetParameterValues()
}

// ISubcommandExecutorWithOptions executes CLI subcommands with both parameter and option support.
type ISubcommandExecutorWithOptions interface {
	ISubcommandExecutor
	IOptionStore
}

// ISubcommandStore provides read-only access to CLI subcommands (methods) available on an object.
type ISubcommandStore interface {
	da.IReadOnlyKeyValueStore[string, reflect.Method]
	GetSubcommands(createSubcommandExecutor func(reflect.Method, any) ISubcommandExecutor) []interface_components.Subcommand
}

// ISubmoduleStore provides access to CLI submodules (nested objects) available on an object.
type ISubmoduleStore interface {
	da.IKeyValueStore[string, any]
	GetSubmodules() []interface_components.Submodule
}

// ICliComponentStoreFactory creates CLI component stores.
type ICliComponentStoreFactory interface {
	CreateOptionStore(obj any) IOptionStore
	CreateSubcommandExecutor(method reflect.Method, obj any) ISubcommandExecutor
	CreateSubcommandExecutorWithOptions(method reflect.Method, obj any) ISubcommandExecutorWithOptions
	CreateSubcommandStore(obj any) ISubcommandStore
	CreateSubmoduleStore(obj any) ISubmoduleStore
}
