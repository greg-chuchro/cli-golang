package completion

import (
	"context"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/greg-chuchro/cli-golang/cli/common"
	"github.com/greg-chuchro/cli-golang/cli/interface_components"
	attr "github.com/greg-chuchro/cli-golang/cli_attributes"
	"github.com/greg-chuchro/cli-golang/cli_attributes/completion_providers"
)

// CompletionPrinter prints completion suggestions based on CLI state and metadata.
type CompletionPrinter struct{}

func NewCompletionPrinter() *CompletionPrinter {
	return &CompletionPrinter{}
}

func (this *CompletionPrinter) PrintSubmodules(context common.ICliContext, submodules []interface_components.Submodule, partialInput string) {
	completions := uniqueSorted(selectKeys(submodules, partialInput))
	for _, c := range completions {
		writeln(context, c)
	}
}

func (this *CompletionPrinter) PrintSubcommands(context common.ICliContext, subcommands []interface_components.Subcommand, partialInput string) {
	completions := uniqueSorted(selectKeys(subcommands, partialInput))
	for _, c := range completions {
		writeln(context, c)
	}
}

func (this *CompletionPrinter) PrintParametersAndOptions(context common.ICliContext, parameters []interface_components.Parameter, options []interface_components.Option, partialInput string) {
	partial := strings.TrimLeft(partialInput, "-+")
	combined := []string{}
	for _, p := range parameters {
		combined = append(combined, p.Keys...)
	}
	for _, o := range options {
		combined = append(combined, o.Keys...)
	}
	completions := uniqueSortedFiltered(combined, partial)
	for _, k := range completions {
		if len(k) == 1 {
			writeln(context, "-"+k)
		} else {
			writeln(context, "--"+k)
		}
	}
}

func (this *CompletionPrinter) PrintSubmoduleValue(context common.ICliContext, submodule interface_components.Submodule, partialInput string) {
}

func (this *CompletionPrinter) PrintSubcommandValue(context common.ICliContext, subcommand interface_components.Subcommand, partialInput string) {
}

func (this *CompletionPrinter) PrintOptionValue(context common.ICliContext, option interface_components.Option, partialInput string, submodule any) {
	completions := this.getValueCompletions(option.ValueType, option.MemberInfo, partialInput, submodule)
	for _, c := range completions {
		writeln(context, c)
	}
}

func (this *CompletionPrinter) PrintParameterValue(context common.ICliContext, parameter interface_components.Parameter, partialInput string, submodule any) {
	completions := this.getValueCompletions(parameter.ValueType, parameter.ParameterInfo, partialInput, submodule)
	for _, c := range completions {
		writeln(context, c)
	}
}

func (this *CompletionPrinter) getValueCompletions(t reflect.Type, attributeProvider reflect.StructField, partialInput string, submodule any) []string {
	attrs, _ := attr.GetCliCompletionAttributes(attributeProvider)
	if len(attrs) > 0 {
		ctx := completion_providers.NewCompletionProviderContext(partialInput, submodule, attrs[0].Filter)
		provider := reflect.New(attrs[0].ProviderType.Elem()).Interface().(completion_providers.ICompletionProvider)
		return provider.GetCompletions(context.Background(), ctx)
	}

	if t != nil && t.Kind() == reflect.Bool {
		return filterPrefix([]string{"true", "false"}, partialInput)
	}

	return []string{}
}

func selectKeys[T any](items []T, partialInput string) []string {
	out := []string{}
	for _, it := range items {
		v := reflect.ValueOf(it)
		keys := v.FieldByName("Keys").Interface().([]string)
		for _, k := range keys {
			if strings.HasPrefix(strings.ToLower(k), strings.ToLower(partialInput)) {
				out = append(out, k)
			}
		}
	}
	return out
}

func uniqueSorted(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, v := range values {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}

func uniqueSortedFiltered(values []string, partial string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, v := range values {
		if seen[v] {
			continue
		}
		if strings.HasPrefix(strings.ToLower(v), strings.ToLower(partial)) {
			seen[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}

func filterPrefix(values []string, partial string) []string {
	out := []string{}
	for _, v := range values {
		if strings.HasPrefix(strings.ToLower(v), strings.ToLower(partial)) {
			out = append(out, v)
		}
	}
	return out
}

func writeln(context common.ICliContext, s string) {
	io.WriteString(context.InterfaceOut(), s+"\n")
}
