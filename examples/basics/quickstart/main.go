package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/greg-chuchro/cli-golang/cli"
)

// QuickStartManager is the root CLI object. Its exported fields become submodules
// and its exported methods become subcommands. (C#: internal class QuickStartManager)
type QuickStartManager struct {
	// Submodule becomes the 'sub' submodule. (C#: public SubStart Submodule { get; } = new())
	Submodule *SubStart `cli:"sub"`
}

// SubStart is a submodule exposing the 'run' subcommand. (C#: internal class SubStart)
type SubStart struct {
	DefaultValue string `cli:"default-value"`
}

// NewSubStart constructs a SubStart with the C# default value "default".
func NewSubStart() *SubStart {
	return &SubStart{DefaultValue: "default"}
}

// Run is the 'run' subcommand. (C#: [CliName("run")] [CliName("r")] public QuickResult SubRun(...))
//
// In Go a method cannot carry struct tags, so the CliName rename is expressed by the
// method name itself (kebab-cased to "run"); the ["r"] alias has no direct equivalent.
func (this *SubStart) Run(requiredParameter int, optionalParameter int) QuickResult {
	if optionalParameter == 0 {
		optionalParameter = 1
	}
	return NewQuickResult(this.DefaultValue, requiredParameter, optionalParameter)
}

// QuickResult is the subcommand return value. (C#: public record QuickResult(string s, int a, int b))
type QuickResult struct {
	S string `cli:"s"`
	A int    `cli:"a"`
	B int    `cli:"b"`
}

// NewQuickResult constructs a QuickResult. (Go has no primary-constructor records.)
func NewQuickResult(s string, a int, b int) QuickResult {
	return QuickResult{S: s, A: a, B: b}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, r)
			os.Exit(1)
		}
	}()

	result := cli.NewCommandLineInterface().Execute(&QuickStartManager{Submodule: NewSubStart()})
	if _, ok := result.(struct{}); ok {
		return
	}
	if result != nil {
		out, err := json.Marshal(result)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(string(out))
	}
}
