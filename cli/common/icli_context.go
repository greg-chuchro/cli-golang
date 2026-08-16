package common

import "io"

// ICliContext provides context and configuration for CLI operations.
type ICliContext interface {
	// SkipUnknown reports whether unknown options should be skipped instead of throwing.
	SkipUnknown() bool

	// InterfaceOut is the writer for interface description output (help, completions, framework messages).
	InterfaceOut() io.Writer
}
