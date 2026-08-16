package collection_element_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/data_access/extensions"
)

// ListElementStoreBase is the base for list stores that provide key-value access to list elements.
type ListElementStoreBase[TKey comparable, TValue any] struct {
	*CollectionElementStoreBase[TKey, TValue]
	List reflect.Value
}

func (this *ListElementStoreBase[TKey, TValue]) GetValueType(key TKey) reflect.Type {
	return this.List.Type().Elem()
}

func (this *ListElementStoreBase[TKey, TValue]) Add(value TValue) {
	this.List.Set(reflect.Append(this.List, reflect.ValueOf(value)))
}

func (this *ListElementStoreBase[TKey, TValue]) Append(value TValue) any {
	this.Add(value)
	return this.List.Interface()
}

func (this *ListElementStoreBase[TKey, TValue]) AddNew() TValue {
	elementType := this.List.Type().Elem()
	newElement := extensions.CreateInstance(elementType)
	this.List.Set(reflect.Append(this.List, reflect.ValueOf(newElement)))
	return newElement.(TValue)
}
