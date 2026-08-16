package interface_component_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/cli/data_access"
	"github.com/greg-chuchro/cli-golang/cli/interface_components"
)

// OptionStore provides access to CLI options with metadata for help generation.
type OptionStore struct {
	store data_access.ICliKeyValueStore[any]
}

func NewOptionStore(store data_access.ICliKeyValueStore[any]) *OptionStore {
	return &OptionStore{store: store}
}

func (this *OptionStore) Get(key string) any {
	return this.store.Get(key)
}

func (this *OptionStore) Set(key string, value any) {
	this.store.Set(key, value)
}

func (this *OptionStore) GetValueOrInitialize(key string) any {
	return this.store.GetValueOrInitialize(key)
}

func (this *OptionStore) ContainsKey(key string) bool {
	return this.store.ContainsKey(key)
}

func (this *OptionStore) GetValueType(key string) reflect.Type {
	return this.store.GetValueType(key)
}

func (this *OptionStore) GetOptions() []interface_components.Option {
	result := make([]interface_components.Option, 0)
	for _, pair := range this.store.GetAccessorKeysPairs() {
		key := pair.Keys[0]
		result = append(result, interface_components.Option{
			IsMultiValue: this.store.IsMultiValue(key),
			Keys:         pair.Keys,
			MemberInfo:   pair.Accessor,
			ValueType:    this.GetValueType(key),
			Value:        this.toString(this.Get(key)),
		})
	}
	return result
}

func (this *OptionStore) toString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return reflect.ValueOf(v).String()
}
