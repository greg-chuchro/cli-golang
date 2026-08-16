package data_access

import (
	"reflect"
	"regexp"
	"strings"

	attr "github.com/greg-chuchro/cli-golang/cli_attributes"
)

// IClassMemberStringifier converts class member names to CLI-friendly string representations (e.g., kebab-case).
type IClassMemberStringifier interface {
	GetAlternativeNames(info reflect.StructField) []string
	GetRequiredNames(info reflect.StructField) []string
}

// ClassMemberStringifierBase converts class member names to CLI-friendly kebab-case format.
type ClassMemberStringifierBase struct{}

func (this *ClassMemberStringifierBase) GetAlternativeNames(info reflect.StructField) []string {
	return this.getAlternativeNames(info.Name, attr.GetCliNameAttributes(info))
}

func (this *ClassMemberStringifierBase) GetRequiredNames(info reflect.StructField) []string {
	return this.getRequiredNames(info.Name, attr.GetCliNameAttributes(info))
}

func (this *ClassMemberStringifierBase) getAlternativeNames(name string, cliNameAttributes []*attr.CliNameAttribute) []string {
	keys := []string{}
	if len(cliNameAttributes) == 1 {
		cliName := cliNameAttributes[0].Name
		if len(cliName) != 1 {
			keys = append(keys, string(cliName[0]))
		}
	}

	if len(cliNameAttributes) == 0 {
		kebabName := getKebabCase(name)
		if len(kebabName) != 1 {
			keys = append(keys, string(kebabName[0]))
		}
	}

	return keys
}

func (this *ClassMemberStringifierBase) getRequiredNames(name string, cliNameAttributes []*attr.CliNameAttribute) []string {
	keys := []string{}
	for _, a := range cliNameAttributes {
		keys = append(keys, a.Name)
	}

	if len(keys) == 0 {
		keys = append(keys, getKebabCase(name))
	}

	return keys
}

// ClassMemberStringifier converts class member names to CLI-friendly kebab-case format.
type ClassMemberStringifier struct {
	ClassMemberStringifierBase
}

func getKebabCase(value string) string {
	re1 := regexp.MustCompile("([a-z0-9])([A-Z])")
	value = re1.ReplaceAllString(value, "$1-$2")
	re2 := regexp.MustCompile("([a-zA-Z0-9])([A-Z][a-z])")
	value = re2.ReplaceAllString(value, "$1-$2")
	re3 := regexp.MustCompile("[. ]")
	value = re3.ReplaceAllString(value, "-")
	return strings.ToLower(value)
}
