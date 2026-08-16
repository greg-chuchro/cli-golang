package parsing

import (
	da "github.com/greg-chuchro/cli-golang/data_access"
)

// OptionReader reads CLI option/argument pairs against a key-value store.
type OptionReader struct {
	*OptionReaderBase
	Store da.IKeyValueStore[string, any]
}

func NewOptionReader(args []string, store da.IKeyValueStore[string, any]) *OptionReader {
	return &OptionReader{
		OptionReaderBase: NewOptionReaderBase(args),
		Store:            store,
	}
}

func (this *OptionReader) ReadOptions() []OptionTuple {
	return this.Read(this.Store)
}
