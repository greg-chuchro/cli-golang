package interface_component_stores

import (
	"github.com/greg-chuchro/cli-golang/cli/interface_components"
	da "github.com/greg-chuchro/cli-golang/data_access"
)

// SubcommandExecutorWithOptions executes CLI subcommands with both parameter and option support (enables shadowing).
type SubcommandExecutorWithOptions struct {
	*da.DistinctDualKeyValueStoreBase[string, any]
	subcommandExecutor ISubcommandExecutor
	optionStore        IOptionStore
}

func NewSubcommandExecutorWithOptions(subcommandExecutor ISubcommandExecutor, optionStore IOptionStore) *SubcommandExecutorWithOptions {
	s := &SubcommandExecutorWithOptions{
		subcommandExecutor: subcommandExecutor,
		optionStore:        optionStore,
	}
	s.DistinctDualKeyValueStoreBase = da.NewDistinctDualKeyValueStoreBase[string, any](subcommandExecutor, optionStore)
	return s
}

func (this *SubcommandExecutorWithOptions) AddArgument(value any) {
	this.subcommandExecutor.AddArgument(value)
}

func (this *SubcommandExecutorWithOptions) SetParameterValues() {
	this.subcommandExecutor.SetParameterValues()
}

func (this *SubcommandExecutorWithOptions) Invoke() (any, error) {
	this.subcommandExecutor.SetParameterValues()
	return this.subcommandExecutor.Invoke()
}

func (this *SubcommandExecutorWithOptions) GetParameters() []interface_components.Parameter {
	return this.subcommandExecutor.GetParameters()
}

func (this *SubcommandExecutorWithOptions) GetFirstUnassignedParameter() *interface_components.Parameter {
	return this.subcommandExecutor.GetFirstUnassignedParameter()
}

func (this *SubcommandExecutorWithOptions) GetOptions() []interface_components.Option {
	return this.optionStore.GetOptions()
}
