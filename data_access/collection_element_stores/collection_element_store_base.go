package collection_element_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
)

// CollectionElementStoreBase is the base for element stores that provide key-value access to enumerable elements.
// Concrete stores embed this and provide the abstract data-access methods; the collection mutating methods
// throw OperationNotSupported by default, mirroring the C# virtual defaults.
type CollectionElementStoreBase[TKey comparable, TValue any] struct {
	store interface {
		Get(key TKey) TValue
		Set(key TKey, value TValue)
		ContainsKey(key TKey) bool
		GetValueType(key TKey) reflect.Type
		GetValueOrInitialize(key TKey) TValue
	}
}

func NewCollectionElementStoreBase[TKey comparable, TValue any](store interface {
	Get(key TKey) TValue
	Set(key TKey, value TValue)
	ContainsKey(key TKey) bool
	GetValueType(key TKey) reflect.Type
	GetValueOrInitialize(key TKey) TValue
}) *CollectionElementStoreBase[TKey, TValue] {
	return &CollectionElementStoreBase[TKey, TValue]{store: store}
}

func (this *CollectionElementStoreBase[TKey, TValue]) Get(key TKey) TValue {
	return this.store.Get(key)
}

func (this *CollectionElementStoreBase[TKey, TValue]) Set(key TKey, value TValue) {
	this.store.Set(key, value)
}

func (this *CollectionElementStoreBase[TKey, TValue]) ContainsKey(key TKey) bool {
	return this.store.ContainsKey(key)
}

func (this *CollectionElementStoreBase[TKey, TValue]) GetValueType(key TKey) any {
	return this.store.GetValueType(key)
}

func (this *CollectionElementStoreBase[TKey, TValue]) GetValueOrInitialize(key TKey) TValue {
	return this.store.GetValueOrInitialize(key)
}

func (this *CollectionElementStoreBase[TKey, TValue]) Add(value TValue) {
	panic(exceptions.OperationNotSupported("Add"))
}

func (this *CollectionElementStoreBase[TKey, TValue]) Append(value TValue) any {
	panic(exceptions.OperationNotSupported("Append"))
}

func (this *CollectionElementStoreBase[TKey, TValue]) AddNew() TValue {
	panic(exceptions.OperationNotSupported("AddNew"))
}

func (this *CollectionElementStoreBase[TKey, TValue]) Remove(key TKey) {
	panic(exceptions.OperationNotSupported("Remove"))
}
