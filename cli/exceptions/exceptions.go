package exceptions

import (
	"errors"
	"fmt"
)

// CliException is thrown when CLI command parsing or execution fails.
type CliException struct {
	message string
	cause   error
}

func (e *CliException) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

func (e *CliException) Unwrap() error {
	return e.cause
}

func NewCliException(message string) *CliException {
	return &CliException{message: message}
}

func NewCliExceptionWithCause(message string, cause error) *CliException {
	return &CliException{message: message, cause: cause}
}

// CliErrors is a factory for common CliException instances with consistent messaging.
var CliErrors = struct {
	NotAnOption                  func(option string) error
	UnknownOption                func(option string) error
	OptionRequiresValue          func(option string) error
	UnexpectedValue              func(value string) error
	AmbiguousSyntax              func(option string) error
	AmbiguousValue               func(value string, option string) error
	InvalidSubcommand            func(subcommand string) error
	NameCollision                func(name1 string, name2 string) error
	FailedToParseArgument        func(message string, cause error) error
	FailedToAccessData           func(message string, cause error) error
	OptionError                  func(option string, message string, cause error) error
	OptionValueError             func(option string, value string, message string, cause error) error
	EmptyArgument                func() error
	UnsupportedMemberType        func(memberTypeName string) error
	UnsupportedShell             func(shell string) error
	UnknownCompletionAction      func(action string) error
	CompletionRequiresShell      func() error
	CompletionInstallFailed      func(shell string, message string, cause error) error
	CompletionUninstallFailed    func(shell string, message string, cause error) error
	UnableToDetermineProgramName func() error
	ValueOutOfRange              func(min any, max any, cause error) error
	InvalidValueFormat           func(expectedTypeName string, cause error) error
}{
	NotAnOption:         func(option string) error { return NewCliException(fmt.Sprintf("Not an option: %s", option)) },
	UnknownOption:       func(option string) error { return NewCliException(fmt.Sprintf("Unknown option: %s", option)) },
	OptionRequiresValue: func(option string) error { return NewCliException(fmt.Sprintf("Option '%s' requires a value", option)) },
	UnexpectedValue:     func(value string) error { return NewCliException(fmt.Sprintf("Unexpected value: %s", value)) },
	AmbiguousSyntax: func(option string) error {
		return NewCliException(fmt.Sprintf("Ambiguous syntax around '%s' (try using --)", option))
	},
	AmbiguousValue: func(value string, option string) error {
		return NewCliException(fmt.Sprintf("Ambiguous value '%s' for '%s', use option=value format for values starting with '-' or '+'", value, option))
	},
	InvalidSubcommand: func(subcommand string) error {
		return NewCliException(fmt.Sprintf("Invalid subcommand: %s", subcommand))
	},
	NameCollision: func(name1 string, name2 string) error {
		return NewCliException(fmt.Sprintf("CLI name of '%s' collides with '%s'", name1, name2))
	},
	FailedToParseArgument: func(message string, cause error) error {
		return NewCliExceptionWithCause(fmt.Sprintf("Failed to parse argument: %s", message), cause)
	},
	FailedToAccessData: func(message string, cause error) error {
		return NewCliExceptionWithCause(fmt.Sprintf("Failed to access data: %s", message), cause)
	},
	OptionError: func(option string, message string, cause error) error {
		return NewCliExceptionWithCause(fmt.Sprintf("option '%s': %s", option, message), cause)
	},
	OptionValueError: func(option string, value string, message string, cause error) error {
		return NewCliExceptionWithCause(fmt.Sprintf("option '%s=%s': %s", option, value, message), cause)
	},
	EmptyArgument: func() error { return NewCliException("Argument cannot be empty") },
	UnsupportedMemberType: func(memberTypeName string) error {
		return NewCliException(fmt.Sprintf("Unsupported member type: %s", memberTypeName))
	},
	UnsupportedShell: func(shell string) error {
		return NewCliException(fmt.Sprintf("Unsupported shell: %s. Supported shells: bash, zsh, powershell, fish", shell))
	},
	UnknownCompletionAction: func(action string) error {
		return NewCliException(fmt.Sprintf("Unknown completion action: %s. Valid actions: install, uninstall", action))
	},
	CompletionRequiresShell: func() error {
		return NewCliException("Completion command requires a shell parameter. Usage: completion <shell> [install|uninstall]")
	},
	CompletionInstallFailed: func(shell string, message string, cause error) error {
		return NewCliExceptionWithCause(fmt.Sprintf("Failed to install completion script for %s: %s", shell, message), cause)
	},
	CompletionUninstallFailed: func(shell string, message string, cause error) error {
		return NewCliExceptionWithCause(fmt.Sprintf("Failed to uninstall completion script for %s: %s", shell, message), cause)
	},
	UnableToDetermineProgramName: func() error { return NewCliException("Unable to determine program name from entry assembly") },
	ValueOutOfRange: func(min any, max any, cause error) error {
		return errors.New(fmt.Sprintf("Out of range (%v-%v)", min, max))
	},
	InvalidValueFormat: func(expectedTypeName string, cause error) error {
		return errors.New(fmt.Sprintf("Invalid format (expected %s)", expectedTypeName))
	},
}
