package cli

import (
	"os"

	ics "github.com/greg-chuchro/cli-golang/cli/data_access/interface_component_stores"
	"github.com/greg-chuchro/cli-golang/cli/exceptions"
	"github.com/greg-chuchro/cli-golang/cli/parsing"
	da "github.com/greg-chuchro/cli-golang/data_access"
)

// OptionDeserializerConfiguration configures OptionDeserializer behavior.
type OptionDeserializerConfiguration struct {
	SkipUnknown bool
}

// OptionDeserializer deserializes command-line options into object properties and fields.
type OptionDeserializer struct{}

// Deserialize deserializes command-line arguments from the environment into the specified object.
func (this *OptionDeserializer) Deserialize(target any, options *OptionDeserializerConfiguration) {
	args := os.Args[1:]
	this.DeserializeArgs(target, args, options)
}

// DeserializeArgs deserializes the specified arguments into the object's properties and fields.
func (this *OptionDeserializer) DeserializeArgs(target any, args []string, options *OptionDeserializerConfiguration) {
	store := ics.NewCliComponentStoreFactory().CreateOptionStore(target)
	this.DeserializeStore(store, args, options)
}

// DeserializeStore deserializes the specified arguments into the key-value store.
func (this *OptionDeserializer) DeserializeStore(store da.IKeyValueStore[string, any], args []string, options *OptionDeserializerConfiguration) {
	reader := parsing.NewOptionReader(args, store)
	this.deserializeReader(reader, options)
}

func (this *OptionDeserializer) deserializeReader(reader *parsing.OptionReader, options *OptionDeserializerConfiguration) {
	if options == nil {
		options = &OptionDeserializerConfiguration{}
	}
	store := reader.Store
	for _, tuple := range reader.ReadOptions() {
		if tuple.Attr.HasFlag(parsing.OptionFlagValueUnassigned) {
			if tuple.Attr.HasFlag(parsing.OptionFlagUnknown) {
				if options.SkipUnknown {
					continue
				}
				panic(exceptions.CliErrors.NotAnOption(tuple.Option))
			}
			if tuple.Attr.HasFlag(parsing.OptionFlagAmbigousValue) {
				panic(exceptions.CliErrors.AmbiguousSyntax(tuple.Option))
			}
			panic(exceptions.CliErrors.UnexpectedValue(tuple.Option))
		}
		store.Set(tuple.Option, tuple.Value)
	}
}
