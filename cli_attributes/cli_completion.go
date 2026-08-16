package cli_attributes

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/greg-chuchro/cli-golang/cli_attributes/completion_providers"
)

type CliCompletionAttribute struct {
	ProviderType reflect.Type
	Filter       string
}

func NewCliCompletionAttribute(providerType reflect.Type, filter string) *CliCompletionAttribute {
	if !providerType.Implements(reflect.TypeOf((*completion_providers.ICompletionProvider)(nil)).Elem()) {
		panic("provider type must implement ICompletionProvider")
	}

	return &CliCompletionAttribute{
		ProviderType: providerType,
		Filter:       filter,
	}
}

var completionProviderTypes = map[string]reflect.Type{
	"method":                    reflect.TypeOf((*completion_providers.MethodCompletionProvider)(nil)).Elem(),
	"methodcompletionprovider":  reflect.TypeOf((*completion_providers.MethodCompletionProvider)(nil)).Elem(),
	"file":                      reflect.TypeOf((*completion_providers.FileCompletionProvider)(nil)).Elem(),
	"filecompletionprovider":    reflect.TypeOf((*completion_providers.FileCompletionProvider)(nil)).Elem(),
	"directory":                 reflect.TypeOf((*completion_providers.DirectoryCompletionProvider)(nil)).Elem(),
	"directorycompletionprovider": reflect.TypeOf((*completion_providers.DirectoryCompletionProvider)(nil)).Elem(),
	"filesystem":                reflect.TypeOf((*completion_providers.FileSystemCompletionProvider)(nil)).Elem(),
	"filesystemcompletionprovider": reflect.TypeOf((*completion_providers.FileSystemCompletionProvider)(nil)).Elem(),
}

func ParseCliCompletionTag(tagValue string) ([]*CliCompletionAttribute, error) {
	if tagValue == "" {
		return nil, nil
	}

	parts := SplitTagOptions(tagValue)
	if len(parts) == 0 {
		return nil, nil
	}

	attrs := make([]*CliCompletionAttribute, 0, len(parts))
	for _, part := range parts {
		attr, err := parseSingleCompletionPart(TrimSpace(part))
		if err != nil {
			return nil, err
		}
		if attr != nil {
			attrs = append(attrs, attr)
		}
	}

	return attrs, nil
}

func parseSingleCompletionPart(part string) (*CliCompletionAttribute, error) {
	if part == "" {
		return nil, nil
	}

	if idx := strings.Index(part, ":"); idx != -1 {
		prefix := part[:idx]
		value := part[idx+1:]

		providerType, ok := completionProviderTypes[strings.ToLower(prefix)]
		if !ok {
			return nil, fmt.Errorf("unknown completion provider prefix: %s", prefix)
		}

		return NewCliCompletionAttribute(providerType, value), nil
	}

	providerType, ok := completionProviderTypes[strings.ToLower(part)]
	if !ok {
		return nil, fmt.Errorf("unknown completion provider type: %s", part)
	}

	return NewCliCompletionAttribute(providerType, ""), nil
}

func GetCliCompletionAttributes(field reflect.StructField) ([]*CliCompletionAttribute, error) {
	tag := field.Tag.Get(CliCompletionTag)
	if tag == "" {
		return nil, nil
	}
	return ParseCliCompletionTag(tag)
}

func MustGetCliCompletionAttributes(field reflect.StructField) []*CliCompletionAttribute {
	attrs, err := GetCliCompletionAttributes(field)
	if err != nil {
		panic(err)
	}
	return attrs
}
