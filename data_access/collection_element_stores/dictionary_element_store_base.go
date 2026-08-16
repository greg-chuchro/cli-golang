package collection_element_stores

import (
	"reflect"
)

// DictionaryElementStoreBase is the base for dictionary stores that provide key-value access to dictionary elements.
type DictionaryElementStoreBase[TKey comparable, TValue any] struct {
	*CollectionElementStoreBase[TKey, TValue]
	Dictionary reflect.Value
}

func (this *DictionaryElementStoreBase[TKey, TValue]) GetValueType(key TKey) reflect.Type {
	return this.Dictionary.Type().Elem()
}

func (this *DictionaryElementStoreBase[TKey, TValue]) Remove(key TKey) {
	this.Dictionary.SetMapIndex(reflect.ValueOf(key), reflect.Value{})
}
