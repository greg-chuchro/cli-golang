package exceptions

type CliAttributeException struct {
	message string
}

func NewCliAttributeException(message string) *CliAttributeException {
	return &CliAttributeException{message: message}
}

func (e *CliAttributeException) Error() string {
	return e.message
}
