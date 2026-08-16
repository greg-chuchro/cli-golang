package completion

import (
	"github.com/greg-chuchro/cli-golang/cli/common"
	"github.com/greg-chuchro/cli-golang/cli/interface_components"
)

// ICompletionHandler handles completion and completion-script generation.
type ICompletionHandler interface {
	HandleComplete(context common.ICliContext, args []string, target any)
	HandleCompletion(context common.ICliContext, args []string, target any)
}

// ICompletionPrinter prints completion suggestions.
type ICompletionPrinter interface {
	PrintSubmodules(context common.ICliContext, submodules []interface_components.Submodule, partialInput string)
	PrintSubcommands(context common.ICliContext, subcommands []interface_components.Subcommand, partialInput string)
	PrintParametersAndOptions(context common.ICliContext, parameters []interface_components.Parameter, options []interface_components.Option, partialInput string)
	PrintSubmoduleValue(context common.ICliContext, submodule interface_components.Submodule, partialInput string)
	PrintSubcommandValue(context common.ICliContext, subcommand interface_components.Subcommand, partialInput string)
	PrintOptionValue(context common.ICliContext, option interface_components.Option, partialInput string, submodule any)
	PrintParameterValue(context common.ICliContext, parameter interface_components.Parameter, partialInput string, submodule any)
}

// ICompletionScriptGenerator generates and installs shell completion scripts.
type ICompletionScriptGenerator interface {
	SupportedShells() []string
	GenerateScript(shell string, programName string) string
	GetInstallPath(shell string, programName string) string
	InstallScript(shell string, programName string)
	UninstallScript(shell string, programName string) bool
}

// IDotnetSuggestHandler handles the dotnet-suggest protocol.
type IDotnetSuggestHandler interface {
	HandleDotnetSuggest(context common.ICliContext, args []string, target any)
	Register(commandPath string)
	Unregister(commandPath string)
}
