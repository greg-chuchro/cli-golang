package cli_attributes

import "strings"

const CliNameTag = "cli"
const CliCompletionTag = "cli_completion"

func SplitTagOptions(tagValue string) []string {
	parts := strings.Split(tagValue, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

type MethodAttribute struct {
	NameAttributes      []*CliNameAttribute
	CompletionAttribute *CliCompletionAttribute
}

func NewMethodAttribute(nameAttributes []*CliNameAttribute, completionAttribute *CliCompletionAttribute) *MethodAttribute {
	return &MethodAttribute{
		NameAttributes:      nameAttributes,
		CompletionAttribute: completionAttribute,
	}
}

type ParameterAttribute struct {
	NameAttributes      []*CliNameAttribute
	CompletionAttribute *CliCompletionAttribute
}

func NewParameterAttribute(nameAttributes []*CliNameAttribute, completionAttribute *CliCompletionAttribute) *ParameterAttribute {
	return &ParameterAttribute{
		NameAttributes:      nameAttributes,
		CompletionAttribute: completionAttribute,
	}
}
