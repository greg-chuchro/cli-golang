package parsing

import (
	"math/big"
	"strings"

	"github.com/greg-chuchro/cli-golang/cli/exceptions"
	da "github.com/greg-chuchro/cli-golang/data_access"
)

// OptionFlags describes the parsing state of a single option token.
type OptionFlags int

const (
	OptionFlagNone            OptionFlags = 0
	OptionFlagShort           OptionFlags = 1
	OptionFlagPlus            OptionFlags = 2
	OptionFlagValueUnassigned OptionFlags = 4
	OptionFlagNotAnOption     OptionFlags = 8 + OptionFlagValueUnassigned
	OptionFlagUnknown         OptionFlags = 16 + OptionFlagValueUnassigned
	OptionFlagAmbigousValue   OptionFlags = 32 + OptionFlagValueUnassigned
)

func (f OptionFlags) HasFlag(flag OptionFlags) bool {
	return f&flag == flag
}

// OptionTuple is a parsed (option, value, flags) result.
type OptionTuple struct {
	Option string
	Value  string
	Attr   OptionFlags
}

// OptionReaderBase reads CLI option/argument pairs from an argument slice.
type OptionReaderBase struct {
	tokens       []string
	currentIndex int
}

func NewOptionReaderBase(args []string) *OptionReaderBase {
	return &OptionReaderBase{tokens: args, currentIndex: -1}
}

func (this *OptionReaderBase) MoveNext() bool {
	this.currentIndex++
	return this.currentIndex < len(this.tokens)
}

func (this *OptionReaderBase) Current() string {
	if this.currentIndex < 0 || this.currentIndex >= len(this.tokens) {
		return ""
	}
	return this.tokens[this.currentIndex]
}

// Read yields parsed (option, value, flags) tuples using the store to resolve option types.
func (this *OptionReaderBase) Read(store da.IKeyValueStore[string, any]) []OptionTuple {
	result := []OptionTuple{}
	isNumber := func(input string) bool {
		_, ok := new(big.Int).SetString(input, 10)
		return ok
	}

	extractOptionValuePair := func(arg string, optionAttr OptionFlags) (string, string, bool) {
		split := strings.SplitN(arg, "=", 2)
		value := ""
		assigned := false
		if len(split) == 2 {
			value = split[1]
			assigned = true
		}
		option := ""
		if optionAttr.HasFlag(OptionFlagShort) {
			if len(split[0]) > 1 {
				option = split[0][1:]
			}
		} else if len(split[0]) > 2 {
			option = split[0][2:]
		}
		return option, value, assigned
	}

	for this.MoveNext() {
		arg := this.Current()

		if len(arg) == 0 {
			panic(exceptions.CliErrors.EmptyArgument())
		}

		optionAttr := OptionFlagNone
		switch arg[0] {
		case '-':
			if !(len(arg) > 1 && arg[1] == '-') {
				optionAttr |= OptionFlagShort
			}
		case '+':
			optionAttr |= OptionFlagPlus
			if len(arg) < 2 || arg[1] != '+' {
				optionAttr |= OptionFlagShort
			}
		default:
			optionAttr |= OptionFlagNotAnOption
			result = append(result, OptionTuple{Option: arg, Value: "", Attr: optionAttr})
			continue
		}

		option, value, assigned := extractOptionValuePair(arg, optionAttr)
		if !assigned {
			if this.MoveNext() {
				next := this.Current()
				if next == "" || ((len(next) > 0 && next[0] != '-' && next[0] != '+') || isNumber(next)) {
					value = next
				} else {
					optionAttr |= OptionFlagAmbigousValue
					this.currentIndex--
				}
			} else {
				optionAttr |= OptionFlagValueUnassigned
			}
		}

		if optionAttr.HasFlag(OptionFlagShort) {
			for _, r := range option {
				so := string(r)
				if store.ContainsKey(so) {
					// recognized
				} else {
					optionAttr |= OptionFlagUnknown
				}
				result = append(result, OptionTuple{Option: so, Value: value, Attr: optionAttr})
			}
		} else {
			if store.ContainsKey(option) {
				// recognized
			} else {
				optionAttr |= OptionFlagUnknown
			}
			result = append(result, OptionTuple{Option: option, Value: value, Attr: optionAttr})
		}
	}

	return result
}
