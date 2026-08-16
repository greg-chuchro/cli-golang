package collection_element_stores

import "reflect"

type ArrayElementStoreBase[TKey comparable, TValue any] struct {
	*CollectionElementStoreBase[TKey, TValue]
	Array reflect.Value
}

func (this *ArrayElementStoreBase[TKey, TValue]) GetValueType(key TKey) reflect.Type {
	return this.Array.Type().Elem()
}
