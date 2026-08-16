package exceptions

import "fmt"

type CliAttributeErrors struct{}

func (e *CliAttributeErrors) InvalidCompletionProvider(providerTypeName string) *CliAttributeException {
	msg := fmt.Sprintf("Provider type '%s' must implement ICompletionProvider", providerTypeName)
	return NewCliAttributeException(msg)
}

func (e *CliAttributeErrors) CompletionMethodNotFound(methodName string, typeName string) *CliAttributeException {
	msg := fmt.Sprintf("Completion method '%s' not found on type '%s'", methodName, typeName)
	return NewCliAttributeException(msg)
}

func (e *CliAttributeErrors) InvalidCompletionMethodSignature(methodName string, typeName string) *CliAttributeException {
	msg := fmt.Sprintf("Completion method '%s' on type '%s' must have signature: []string %s(string partialInput)", methodName, typeName, methodName)
	return NewCliAttributeException(msg)
}

func (e *CliAttributeErrors) InvalidCompletionMethodReturnType(methodName string, typeName string) *CliAttributeException {
	msg := fmt.Sprintf("Completion method '%s' on type '%s' must return []string", methodName, typeName)
	return NewCliAttributeException(msg)
}
