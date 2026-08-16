package formatting

import (
	"io"
	"reflect"
	"runtime/debug"

	"github.com/greg-chuchro/cli-golang/cli/common"
)

// VersionPrinter prints version information for the CLI application.
type VersionPrinter struct {
	UseRevisionVersion bool
}

func NewVersionPrinter() *VersionPrinter {
	return &VersionPrinter{UseRevisionVersion: false}
}

func (this *VersionPrinter) PrintVersion(context common.ICliContext, rootSubmoduleType reflect.Type) {
	version := readVersion()
	if version != "" {
		io.WriteString(context.InterfaceOut(), version+"\n")
	}
}

func readVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return ""
}
