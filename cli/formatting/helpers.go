package formatting

import (
	"os"
	"reflect"
	"strings"

	"github.com/greg-chuchro/cli-golang/cli/interface_components"
)

// detectColorSupport checks whether the output terminal supports ANSI colors.
func detectColorSupport() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("COLORTERM") != "" {
		return true
	}
	term := os.Getenv("TERM")
	if term != "" && (strings.Contains(term, "color") || strings.HasPrefix(term, "xterm") || strings.HasPrefix(term, "screen") || term == "linux" || term == "cygwin") {
		return true
	}
	return false
}

// Go has no runtime-accessible XML documentation, so summaries are empty.
func getSummary(v any) string {
	return ""
}

func getDescription(v any) string {
	return ""
}

func getParameterDescription(p interface_components.Parameter) string {
	return getDescription(p.ParameterInfo)
}

func getOptionDescription(o interface_components.Option) string {
	return getDescription(o.MemberInfo)
}

var _ reflect.Type
