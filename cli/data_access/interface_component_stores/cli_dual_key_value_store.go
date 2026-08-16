package interface_component_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/cli/data_access"
)

// cliDualKeyValueStore combines two CLI key-value stores with primary and secondary lookup priority.
type cliDualKeyValueStore struct {
	primary   data_access.ICliKeyValueStore[any]
	secondary data_access.ICliKeyValueStore[any]
}

func (this *cliDualKeyValueStore) Get(key string) any {
	if this.primary.ContainsKey(key) {
		return this.primary.Get(key)
	}
	return this.secondary.Get(key)
}

func (this *cliDualKeyValueStore) Set(key string, value any) {
	if this.primary.ContainsKey(key) {
		this.primary.Set(key, value)
	} else {
		this.secondary.Set(key, value)
	}
}

func (this *cliDualKeyValueStore) GetValueOrInitialize(key string) any {
	if this.primary.ContainsKey(key) {
		return this.primary.GetValueOrInitialize(key)
	}
	return this.secondary.GetValueOrInitialize(key)
}

func (this *cliDualKeyValueStore) ContainsKey(key string) bool {
	return this.primary.ContainsKey(key) || this.secondary.ContainsKey(key)
}

func (this *cliDualKeyValueStore) GetValueType(key string) reflect.Type {
	if this.primary.ContainsKey(key) {
		return this.primary.GetValueType(key)
	}
	return this.secondary.GetValueType(key)
}

func (this *cliDualKeyValueStore) GetAccessorKeysPairs() []data_access.AccessorKeysPair[reflect.StructField] {
	return append(this.primary.GetAccessorKeysPairs(), this.secondary.GetAccessorKeysPairs()...)
}

func (this *cliDualKeyValueStore) IsMultiValue(key string) bool {
	if this.primary.ContainsKey(key) {
		return this.primary.IsMultiValue(key)
	}
	return this.secondary.IsMultiValue(key)
}
