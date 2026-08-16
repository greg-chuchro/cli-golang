package cli_attributes

import (
	"reflect"
)

type CliNameAttribute struct {
	Name string
}

func NewCliNameAttribute(name string) *CliNameAttribute {
	return &CliNameAttribute{Name: name}
}

func ParseCliNameTag(tagValue string) []*CliNameAttribute {
	if tagValue == "" {
		return nil
	}

	parts := SplitTagOptions(tagValue)
	attrs := make([]*CliNameAttribute, 0, len(parts))
	for _, part := range parts {
		name := TrimSpace(part)
		if name != "" {
			attrs = append(attrs, NewCliNameAttribute(name))
		}
	}
	return attrs
}

func GetCliNameAttributes(field reflect.StructField) []*CliNameAttribute {
	tag := field.Tag.Get(CliNameTag)
	if tag == "" {
		return nil
	}
	return ParseCliNameTag(tag)
}

func (a *CliNameAttribute) String() string {
	return a.Name
}
