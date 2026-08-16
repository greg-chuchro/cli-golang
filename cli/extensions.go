package cli

import (
	"os"
	"path/filepath"
)

// GetToolCommandName returns the command name for the current executable.
// Go has no .NET global tool registry; this returns the executable base name.
func GetToolCommandName() string {
	if exe := os.Args[0]; exe != "" {
		return filepath.Base(exe)
	}
	return "program"
}
