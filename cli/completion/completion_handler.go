package completion

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/greg-chuchro/cli-golang/cli/common"
	ics "github.com/greg-chuchro/cli-golang/cli/data_access/interface_component_stores"
	"github.com/greg-chuchro/cli-golang/cli/exceptions"
)

// CompletionHandler handles completion and completion-script generation.
type CompletionHandler struct {
	CliComponentStoreFactory  ics.ICliComponentStoreFactory
	CompletionPrinter         ICompletionPrinter
	CompletionScriptGenerator ICompletionScriptGenerator
}

func NewCompletionHandler() *CompletionHandler {
	factory := ics.NewCliComponentStoreFactory()
	return &CompletionHandler{
		CliComponentStoreFactory:  factory,
		CompletionPrinter:         NewCompletionPrinter(),
		CompletionScriptGenerator: NewCompletionScriptGenerator(),
	}
}

func (this *CompletionHandler) HandleComplete(context common.ICliContext, args []string, target any) {
	argsList := args
	if len(argsList) == 0 {
		this.handleComplete(context, target, []string{}, "")
		return
	}
	toComplete := argsList[len(argsList)-1]
	argsBefore := argsList[:len(argsList)-1]
	this.handleComplete(context, target, argsBefore, toComplete)
}

func (this *CompletionHandler) HandleCompletion(context common.ICliContext, args []string, target any) {
	argsList := args
	if len(argsList) == 0 {
		panic(exceptions.CliErrors.CompletionRequiresShell())
	}
	shell := argsList[0]
	isAllShells := strings.EqualFold(shell, "all")
	if !isAllShells && !contains(this.CompletionScriptGenerator.SupportedShells(), strings.ToLower(shell)) {
		panic(exceptions.CliErrors.UnsupportedShell(shell))
	}
	remaining := argsList[1:]
	if len(remaining) > 0 {
		action := strings.ToLower(remaining[0])
		switch action {
		case "install":
			if isAllShells {
				this.handleCompletionInstallAll(context)
			} else {
				this.handleCompletionInstall(context, shell)
			}
			return
		case "uninstall":
			if isAllShells {
				this.handleCompletionUninstallAll(context)
			} else {
				this.handleCompletionUninstall(context, shell)
			}
			return
		default:
			panic(exceptions.CliErrors.UnknownCompletionAction(action))
		}
	}
	if isAllShells {
		this.handleCompletionScriptAll(context)
	} else {
		this.handleCompletionScript(context, shell)
	}
}

func (this *CompletionHandler) handleComplete(context common.ICliContext, target any, args []string, toComplete string) {
	en := newArgEnumerator(args)
	factory := this.CliComponentStoreFactory
	submodule := target
	submoduleStore := factory.CreateSubmoduleStore(submodule)
	subcommandName := ""
	for en.MoveNext() {
		arg := en.Current()
		if submoduleStore.ContainsKey(arg) {
			submodule = submoduleStore.Get(arg)
			submoduleStore = factory.CreateSubmoduleStore(submodule)
		} else {
			subcommandName = arg
			break
		}
	}

	subcommandStore := factory.CreateSubcommandStore(submodule)
	if subcommandName == "" {
		submodules := submoduleStore.GetSubmodules()
		subcommands := subcommandStore.GetSubcommands(factory.CreateSubcommandExecutor)
		this.CompletionPrinter.PrintSubmodules(context, submodules, toComplete)
		this.CompletionPrinter.PrintSubcommands(context, subcommands, toComplete)
		return
	}

	if !subcommandStore.ContainsKey(subcommandName) {
		this.CompletionPrinter.PrintSubcommands(context, subcommandStore.GetSubcommands(factory.CreateSubcommandExecutor), subcommandName)
		return
	}

	method := subcommandStore.Get(subcommandName)
	executorWithOptions := factory.CreateSubcommandExecutorWithOptions(method, submodule)

	completingOptionValue := false
	var previousArg string
	if len(args) > 0 {
		previousArg = args[len(args)-1]
		completingOptionValue = strings.HasPrefix(previousArg, "-") && !strings.HasPrefix(toComplete, "-")
	}

	if strings.HasPrefix(toComplete, "-") {
		this.CompletionPrinter.PrintParametersAndOptions(context, executorWithOptions.GetParameters(), executorWithOptions.GetOptions(), toComplete)
	} else if completingOptionValue && previousArg != "" {
		optionName := strings.TrimLeft(previousArg, "-+")
		if p := findParameter(executorWithOptions.GetParameters(), optionName); p != nil {
			this.CompletionPrinter.PrintParameterValue(context, *p, toComplete, submodule)
		} else if o := findOption(executorWithOptions.GetOptions(), optionName); o != nil {
			this.CompletionPrinter.PrintOptionValue(context, *o, toComplete, submodule)
		}
	} else {
		if p := executorWithOptions.GetFirstUnassignedParameter(); p != nil {
			this.CompletionPrinter.PrintParameterValue(context, *p, toComplete, submodule)
		}
	}
}

func (this *CompletionHandler) handleCompletionScript(context common.ICliContext, shell string) {
	programName := getProgramName()
	script := this.CompletionScriptGenerator.GenerateScript(shell, programName)
	io.WriteString(context.InterfaceOut(), script+"\n")
}

func (this *CompletionHandler) handleCompletionInstall(context common.ICliContext, shell string) {
	programName := getProgramName()
	this.CompletionScriptGenerator.InstallScript(shell, programName)
	io.WriteString(context.InterfaceOut(), "Completion script installed.\n")
}

func (this *CompletionHandler) handleCompletionUninstall(context common.ICliContext, shell string) {
	programName := getProgramName()
	if removed := this.CompletionScriptGenerator.UninstallScript(shell, programName); removed {
		io.WriteString(context.InterfaceOut(), "Completion script has been uninstalled.\n")
	} else {
		io.WriteString(context.InterfaceOut(), "No completion script found.\n")
	}
}

func (this *CompletionHandler) handleCompletionScriptAll(context common.ICliContext) {
	programName := getProgramName()
	for _, shell := range this.CompletionScriptGenerator.SupportedShells() {
		io.WriteString(context.InterfaceOut(), "# Completion script for "+shell+"\n\n")
		io.WriteString(context.InterfaceOut(), this.CompletionScriptGenerator.GenerateScript(shell, programName)+"\n\n")
	}
}

func (this *CompletionHandler) handleCompletionInstallAll(context common.ICliContext) {
	programName := getProgramName()
	for _, shell := range this.CompletionScriptGenerator.SupportedShells() {
		this.CompletionScriptGenerator.InstallScript(shell, programName)
		io.WriteString(context.InterfaceOut(), "["+shell+"] installed.\n")
	}
}

func (this *CompletionHandler) handleCompletionUninstallAll(context common.ICliContext) {
	programName := getProgramName()
	for _, shell := range this.CompletionScriptGenerator.SupportedShells() {
		this.CompletionScriptGenerator.UninstallScript(shell, programName)
		io.WriteString(context.InterfaceOut(), "["+shell+"] uninstalled.\n")
	}
}

func getProgramName() string {
	if exe := os.Args[0]; exe != "" {
		return filepath.Base(exe)
	}
	return "program"
}
