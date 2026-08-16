package exceptions

import (
	"errors"
	"fmt"
)

// DataAccessException is thrown when data access operations fail.
type DataAccessException struct {
	message string
	cause   error
}

func (e *DataAccessException) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

func (e *DataAccessException) Unwrap() error {
	return e.cause
}

func NewDataAccessException(message string) *DataAccessException {
	return &DataAccessException{message: message}
}

func NewDataAccessExceptionWithCause(message string, cause error) *DataAccessException {
	return &DataAccessException{message: message, cause: cause}
}

// DataAccessErrors is a factory for common DataAccessException instances with consistent messaging.
var DataAccessErrors = struct {
	KeyNotFound            func(key any) error
	AmbiguousKey           func(key any) error
	UnexpectedArgument     func(argument any) error
	UnassignedParameter    func(parameterName string) error
	NoAccessConfigured     func() error
	InvalidListKey         func(keyTypeName string) error
	TypeCannotBeParsed     func(typeName string) error
	UnrecognizedMemberType func() error
	OperationNotSupported  func(operation string) error
	CannotCreateInstance   func(typeName string) error
}{
	KeyNotFound: func(key any) error {
		return NewDataAccessException(fmt.Sprintf("Key '%v' not found.", key))
	},
	AmbiguousKey: func(key any) error {
		return NewDataAccessException(fmt.Sprintf("Ambiguous key '%v', enable shadowing to ignore this error", key))
	},
	UnexpectedArgument: func(argument any) error {
		return NewDataAccessException(fmt.Sprintf("Unexpected argument: %v", argument))
	},
	UnassignedParameter: func(parameterName string) error {
		return NewDataAccessException(fmt.Sprintf("Unassigned parameter: %s", parameterName))
	},
	NoAccessConfigured: func() error {
		return errors.New("Neither AccessFields nor AccessProperties is set")
	},
	InvalidListKey: func(keyTypeName string) error {
		if keyTypeName == "" {
			keyTypeName = "null"
		}
		return errors.New(fmt.Sprintf("Key must be an integer or parsable string for list removal, got %s", keyTypeName))
	},
	TypeCannotBeParsed: func(typeName string) error {
		return errors.New(fmt.Sprintf("Type cannot be parsed: %s", typeName))
	},
	UnrecognizedMemberType: func() error {
		return errors.New("MemberInfo is not a recognized type")
	},
	OperationNotSupported: func(operation string) error {
		return errors.New(fmt.Sprintf("Operation '%s' is not supported", operation))
	},
	CannotCreateInstance: func(typeName string) error {
		return errors.New(fmt.Sprintf("Cannot create instance of type '%s'", typeName))
	},
}

// Free-function helpers mirroring the DataAccessErrors factory methods.
func KeyNotFound(key any) error {
	return DataAccessErrors.KeyNotFound(key)
}

func AmbiguousKey(key any) error {
	return DataAccessErrors.AmbiguousKey(key)
}

func UnexpectedArgument(argument any) error {
	return DataAccessErrors.UnexpectedArgument(argument)
}

func UnassignedParameter(parameterName string) error {
	return DataAccessErrors.UnassignedParameter(parameterName)
}

func NoAccessConfigured() error {
	return DataAccessErrors.NoAccessConfigured()
}

func InvalidListKey(keyTypeName string) error {
	return DataAccessErrors.InvalidListKey(keyTypeName)
}

func TypeCannotBeParsed(typeName string) error {
	return DataAccessErrors.TypeCannotBeParsed(typeName)
}

func UnrecognizedMemberType() error {
	return DataAccessErrors.UnrecognizedMemberType()
}

func OperationNotSupported(operation string) error {
	return DataAccessErrors.OperationNotSupported(operation)
}

func CannotCreateInstance(typeName string) error {
	return DataAccessErrors.CannotCreateInstance(typeName)
}
