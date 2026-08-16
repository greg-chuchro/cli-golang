package interface_component_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/cli/data_access"
	"github.com/greg-chuchro/cli-golang/cli/interface_components"
)

// SubmoduleStore provides access to CLI submodules (nested objects) with metadata.
type SubmoduleStore struct {
	store data_access.ICliKeyValueStore[any]
}

func NewSubmoduleStore(store data_access.ICliKeyValueStore[any]) *SubmoduleStore {
	return &SubmoduleStore{store: store}
}

func (this *SubmoduleStore) Get(key string) any {
	return this.store.Get(key)
}

func (this *SubmoduleStore) Set(key string, value any) {
	this.store.Set(key, value)
}

func (this *SubmoduleStore) GetValueOrInitialize(key string) any {
	return this.store.GetValueOrInitialize(key)
}

func (this *SubmoduleStore) ContainsKey(key string) bool {
	return this.store.ContainsKey(key)
}

func (this *SubmoduleStore) GetValueType(key string) reflect.Type {
	return this.store.GetValueType(key)
}

func (this *SubmoduleStore) GetSubmodules() []interface_components.Submodule {
	result := make([]interface_components.Submodule, 0)
	for _, pair := range this.store.GetAccessorKeysPairs() {
		result = append(result, interface_components.Submodule{
			Keys:       pair.Keys,
			MemberInfo: pair.Accessor,
		})
	}
	return result
}
