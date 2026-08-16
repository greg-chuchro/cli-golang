package cli

import (
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/greg-chuchro/cli-golang/cli/completion"
	ics "github.com/greg-chuchro/cli-golang/cli/data_access/interface_component_stores"
	"github.com/greg-chuchro/cli-golang/cli/exceptions"
	"github.com/greg-chuchro/cli-golang/cli/formatting"
	"github.com/greg-chuchro/cli-golang/cli/interface_components"
	"github.com/greg-chuchro/cli-golang/cli/parsing"
)

// CommandLineInterface interprets CLI commands and executes methods on any classlib without requiring programming.
type CommandLineInterface struct {
	HelpPrinter              formatting.IHelpPrinter
	CompletionHandler        completion.ICompletionHandler
	DotnetSuggestHandler     completion.IDotnetSuggestHandler
	CliComponentStoreFactory ics.ICliComponentStoreFactory
	SkipUnknownOption        bool
	interfaceOut             io.Writer
}

// NewCommandLineInterface creates a new CLI with default dependencies.
func NewCommandLineInterface() *CommandLineInterface {
	cli := &CommandLineInterface{}
	cli.CliComponentStoreFactory = ics.NewCliComponentStoreFactory()
	cli.HelpPrinter = formatting.NewHelpPrinter()
	cli.CompletionHandler = completion.NewCompletionHandler()
	cli.DotnetSuggestHandler = completion.NewDotnetSuggestHandler(cli.CompletionHandler)
	cli.SkipUnknownOption = false
	return cli
}

// SkipUnknown reports whether unknown options should be skipped instead of throwing.
func (this *CommandLineInterface) SkipUnknown() bool {
	return this.SkipUnknownOption
}

// InterfaceOut returns the writer for interface description output.
func (this *CommandLineInterface) InterfaceOut() io.Writer {
	if this.interfaceOut != nil {
		return this.interfaceOut
	}
	return os.Stdout
}

// SetInterfaceOut sets the writer for interface description output.
func (this *CommandLineInterface) SetInterfaceOut(w io.Writer) {
	this.interfaceOut = w
}

// Execute runs a CLI command using command-line arguments from the environment.
func (this *CommandLineInterface) Execute(target any) any {
	args := os.Args[1:]
	return this.ExecuteArgs(target, args)
}

// ExecuteArgs runs a CLI command using the provided arguments.
func (this *CommandLineInterface) ExecuteArgs(target any, args []string) any {
	if len(args) == 0 {
		return nil
	}
	firstArg := args[0]
	switch firstArg {
	case "__complete":
		this.CompletionHandler.HandleComplete(this, args, target)
	case "completion":
		this.CompletionHandler.HandleCompletion(this, args, target)
	case "[suggest", "[suggest:0]":
		if strings.HasPrefix(firstArg, "[suggest") {
			this.DotnetSuggestHandler.HandleDotnetSuggest(this, args, target)
		}
	default:
		return this.executeInvoke(target, args)
	}
	return nil
}

func (this *CommandLineInterface) executeInvoke(target any, args []string) any {
	en := newSliceEnumerator(args)

	parentSubmodule := any(nil)
	submodule := target
	submoduleStore := this.CliComponentStoreFactory.CreateSubmoduleStore(submodule)
	subcommandName := ""
	for en.MoveNext() {
		arg := en.Current()
		if submoduleStore.ContainsKey(arg) {
			parentSubmodule = submodule
			submodule = submoduleStore.Get(arg)
			submoduleStore = this.CliComponentStoreFactory.CreateSubmoduleStore(submodule)
		} else {
			subcommandName = arg
			break
		}
	}

	if subcommandName == "--version" || subcommandName == "-v" {
		this.VersionPrinter().PrintVersion(this, reflect.TypeOf(target))
		return nil
	}

	subcommandStore := this.CliComponentStoreFactory.CreateSubcommandStore(submodule)

	if subcommandName == "" || subcommandName == "--help" || subcommandName == "-h" {
		optionStore := this.CliComponentStoreFactory.CreateOptionStore(submodule)
		isRoot := submodule == target
		if isRoot {
			this.HelpPrinter.PrintRootHelp(this, reflect.TypeOf(target), submoduleStore.GetSubmodules(), subcommandStore.GetSubcommands(this.CliComponentStoreFactory.CreateSubcommandExecutor), optionStore.GetOptions())
			return nil
		}
		parentSubmoduleStore := this.CliComponentStoreFactory.CreateSubmoduleStore(parentSubmodule)
		var parentModuleInfo interface_components.Submodule
		for _, s := range parentSubmoduleStore.GetSubmodules() {
			if parentSubmoduleStore.ContainsKey(s.Keys[0]) && parentSubmoduleStore.Get(s.Keys[0]) == submodule {
				parentModuleInfo = s
				break
			}
		}
		this.HelpPrinter.PrintSubmoduleHelp(this, reflect.TypeOf(target), parentModuleInfo, submoduleStore.GetSubmodules(), subcommandStore.GetSubcommands(this.CliComponentStoreFactory.CreateSubcommandExecutor), optionStore.GetOptions())
		return nil
	}

	method := subcommandStore.Get(subcommandName)
	executorWithOptions := this.CliComponentStoreFactory.CreateSubcommandExecutorWithOptions(method, submodule)

	if en.MoveNext() {
		firstArg := en.Current()
		if firstArg == "--help" || firstArg == "-h" {
			var matched interface_components.Subcommand
			for _, sc := range subcommandStore.GetSubcommands(this.CliComponentStoreFactory.CreateSubcommandExecutor) {
				if sc.MethodInfo.Name == method.Name {
					matched = sc
					break
				}
			}
			this.HelpPrinter.PrintSubcommandHelp(this, reflect.TypeOf(target), matched, executorWithOptions.GetOptions())
			return nil
		}
		skipped := newSliceEnumerator(sliceFrom(en))
		this.readParametersAndOptions(skipped, executorWithOptions)
	}

	result, err := executorWithOptions.Invoke()
	if err != nil {
		panic(exceptions.CliErrors.FailedToAccessData(err.Error(), err))
	}
	return result
}

func (this *CommandLineInterface) readParametersAndOptions(args *sliceEnumerator, executorWithOptions ics.ISubcommandExecutorWithOptions) {
	optionReader := parsing.NewOptionReader(args.tokens, executorWithOptions)
	for _, tuple := range optionReader.ReadOptions() {
		if tuple.Attr.HasFlag(parsing.OptionFlagAmbigousValue) {
			panic(exceptions.CliErrors.AmbiguousValue(args.Current(), tuple.Option))
		}
		if tuple.Attr.HasFlag(parsing.OptionFlagUnknown) {
			if this.SkipUnknownOption {
				continue
			}
			panic(exceptions.CliErrors.UnknownOption(tuple.Option))
		}
		if tuple.Attr.HasFlag(parsing.OptionFlagValueUnassigned) && !tuple.Attr.HasFlag(parsing.OptionFlagNotAnOption) {
			panic(exceptions.CliErrors.OptionRequiresValue(tuple.Option))
		}
		if tuple.Attr.HasFlag(parsing.OptionFlagNotAnOption) {
			executorWithOptions.AddArgument(tuple.Option)
		} else {
			executorWithOptions.Set(tuple.Option, tuple.Value)
		}
	}
}

func (this *CommandLineInterface) VersionPrinter() formatting.IVersionPrinter {
	return formatting.NewVersionPrinter()
}

func sliceFrom(en *sliceEnumerator) []string {
	if en.index < 0 || en.index >= len(en.tokens) {
		return []string{}
	}
	out := make([]string, 0, len(en.tokens)-en.index)
	out = append(out, en.tokens[en.index:]...)
	return out
}
