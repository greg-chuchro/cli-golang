package formatting

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/greg-chuchro/cli-golang/cli/common"
	"github.com/greg-chuchro/cli-golang/cli/interface_components"
)

// IHelpPrinter prints formatted help information for CLI components.
type IHelpPrinter interface {
	PrintRootHelp(context common.ICliContext, rootType reflect.Type, submodules []interface_components.Submodule, subcommands []interface_components.Subcommand, options []interface_components.Option)
	PrintSubmoduleHelp(context common.ICliContext, rootType reflect.Type, submodule interface_components.Submodule, submodules []interface_components.Submodule, subcommands []interface_components.Subcommand, options []interface_components.Option)
	PrintSubcommandHelp(context common.ICliContext, rootType reflect.Type, subcommand interface_components.Subcommand, options []interface_components.Option)
}

// IVersionPrinter prints version information for the CLI application.
type IVersionPrinter interface {
	PrintVersion(context common.ICliContext, rootSubmoduleType reflect.Type)
}

type itemInfo struct {
	Description string
	Keys        []string
}

type sectionInfo struct {
	Title     string
	ItemInfos []itemInfo
}

// HelpPrinter prints formatted help information using ANSI color output.
type HelpPrinter struct {
	supportsColor bool
}

func NewHelpPrinter() *HelpPrinter {
	return &HelpPrinter{supportsColor: detectColorSupport()}
}

func (this *HelpPrinter) PrintSubmoduleHelp(context common.ICliContext, rootType reflect.Type, submodule interface_components.Submodule, submodules []interface_components.Submodule, subcommands []interface_components.Subcommand, options []interface_components.Option) {
	description := getSummary(submodule.MemberInfo)
	if description != "" {
		this.writeln(context, description)
		this.writeln(context, "")
	}
	this.printHelp(context, submodules, subcommands, options)
}

func (this *HelpPrinter) PrintRootHelp(context common.ICliContext, rootType reflect.Type, submodules []interface_components.Submodule, subcommands []interface_components.Subcommand, options []interface_components.Option) {
	rootDescription := getSummary(rootType)
	if rootDescription != "" {
		this.writeln(context, rootDescription)
		this.writeln(context, "")
	}
	this.printHelp(context, submodules, subcommands, options)
}

func (this *HelpPrinter) PrintSubcommandHelp(context common.ICliContext, rootType reflect.Type, subcommand interface_components.Subcommand, options []interface_components.Option) {
	this.printSubcommandDescription(context, subcommand)
	sections := []sectionInfo{
		{
			Title: "Parameters",
			ItemInfos: func() []itemInfo {
				out := make([]itemInfo, 0, len(subcommand.Parameters))
				for _, p := range subcommand.Parameters {
					out = append(out, itemInfo{Description: getParameterDescription(p), Keys: mapKeys(p.Keys)})
				}
				return out
			}(),
		},
		{
			Title: "Options",
			ItemInfos: func() []itemInfo {
				out := make([]itemInfo, 0, len(options))
				for _, o := range options {
					out = append(out, itemInfo{Description: getOptionDescription(o), Keys: mapKeys(o.Keys)})
				}
				return out
			}(),
		},
	}
	this.printSections(context, sections)
}

func (this *HelpPrinter) printHelp(context common.ICliContext, submodules []interface_components.Submodule, subcommands []interface_components.Subcommand, options []interface_components.Option) {
	sections := []sectionInfo{
		{Title: "Submodules", ItemInfos: mapItems(submodules, func(s interface_components.Submodule) (string, []string) { return getDescription(s.MemberInfo), s.Keys })},
		{Title: "Subcommands", ItemInfos: mapItems(subcommands, func(s interface_components.Subcommand) (string, []string) {
			return getDescription(s.MethodInfo), s.Keys
		})},
		{Title: "Options", ItemInfos: mapItems(options, func(o interface_components.Option) (string, []string) {
			return getOptionDescription(o), mapKeys(o.Keys)
		})},
	}
	this.printSections(context, sections)
}

func (this *HelpPrinter) printSections(context common.ICliContext, sections []sectionInfo) {
	maxLengths := normalizeKeyCounts(sections)

	first := true
	for _, section := range sections {
		if len(section.ItemInfos) == 0 {
			continue
		}
		if !first {
			this.writeln(context, "")
		}
		first = false
		this.setColor(context, 90, 147, 241)
		this.writeln(context, section.Title)
		for _, item := range section.ItemInfos {
			this.write(context, "  ")
			parts := make([]string, len(maxLengths))
			keys := item.Keys
			for i := 0; i < len(maxLengths); i++ {
				if i < len(keys) {
					if i == 0 {
						parts[i] = padRight(keys[i], maxLengths[i])
					} else {
						parts[i] = padLeft(keys[i], maxLengths[i])
					}
				} else {
					parts[i] = ""
				}
			}
			this.setColor(context, 149, 184, 204)
			this.write(context, strings.Join(parts, " "))
			this.resetColor(context)
			this.write(context, "  ")
			this.writeln(context, item.Description)
		}
	}
	this.resetColor(context)
}

func (this *HelpPrinter) printSubcommandDescription(context common.ICliContext, subcommand interface_components.Subcommand) {
	parts := []string{}
	summary := getSummary(subcommand.MethodInfo)
	if summary != "" {
		parts = append(parts, summary)
	}
	description := strings.Join(parts, "\n")
	if description != "" {
		this.writeln(context, description)
		this.writeln(context, "")
	}
}

func (this *HelpPrinter) write(context common.ICliContext, s string) {
	io.WriteString(context.InterfaceOut(), s)
}

func (this *HelpPrinter) writeln(context common.ICliContext, s string) {
	io.WriteString(context.InterfaceOut(), s+"\n")
}

func (this *HelpPrinter) setColor(context common.ICliContext, r, g, b int) {
	if this.supportsColor {
		this.write(context, fmt.Sprintf("\u001b[38;2;%d;%d;%dm", r, g, b))
	}
}

func (this *HelpPrinter) resetColor(context common.ICliContext) {
	if this.supportsColor {
		this.write(context, "\u001b[0m")
	}
}

func mapKeys(keys []string) []string {
	out := make([]string, len(keys))
	for i, k := range keys {
		if len(k) > 1 {
			out[i] = "--" + k
		} else {
			out[i] = "-" + k
		}
	}
	return out
}

func mapItems[T any](items []T, f func(T) (string, []string)) []itemInfo {
	out := make([]itemInfo, 0, len(items))
	for _, it := range items {
		desc, keys := f(it)
		out = append(out, itemInfo{Description: desc, Keys: keys})
	}
	return out
}

func getMaxKeyLengths(sections []sectionInfo) []int {
	maxLengths := []int{}
	for _, section := range sections {
		for _, item := range section.ItemInfos {
			for i, k := range item.Keys {
				if len(maxLengths) <= i {
					maxLengths = append(maxLengths, 0)
				}
				if len(k) > maxLengths[i] {
					maxLengths[i] = len(k)
				}
			}
		}
	}
	return maxLengths
}

func normalizeKeyCounts(sections []sectionInfo) []int {
	maxLengths := getMaxKeyLengths(sections)
	for _, section := range sections {
		for _, item := range section.ItemInfos {
			if len(item.Keys) < len(maxLengths) {
				for i := len(item.Keys); i < len(maxLengths); i++ {
					item.Keys = append(item.Keys, "")
				}
			}
		}
	}
	return getMaxKeyLengths(sections)
}

func padRight(s string, n int) string {
	for len(s) < n {
		s += " "
	}
	return s
}

func padLeft(s string, n int) string {
	for len(s) < n {
		s = " " + s
	}
	return s
}
