package completion

import (
	"github.com/greg-chuchro/cli-golang/cli/interface_components"
)

// argEnumerator walks a slice of strings.
type argEnumerator struct {
	tokens []string
	index  int
}

func newArgEnumerator(tokens []string) *argEnumerator {
	return &argEnumerator{tokens: tokens, index: -1}
}

func (this *argEnumerator) MoveNext() bool {
	this.index++
	return this.index < len(this.tokens)
}

func (this *argEnumerator) Current() string {
	if this.index < 0 || this.index >= len(this.tokens) {
		return ""
	}
	return this.tokens[this.index]
}

func contains(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}

func findParameter(params []interface_components.Parameter, name string) *interface_components.Parameter {
	for i := range params {
		for _, k := range params[i].Keys {
			if k == name {
				return &params[i]
			}
		}
	}
	return nil
}

func findOption(options []interface_components.Option, name string) *interface_components.Option {
	for i := range options {
		for _, k := range options[i].Keys {
			if k == name {
				return &options[i]
			}
		}
	}
	return nil
}
