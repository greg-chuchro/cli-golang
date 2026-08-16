package data_access

import (
	"reflect"

	ces "github.com/greg-chuchro/cli-golang/data_access/collection_element_stores"
)

// ICollectionElementStoreFactory creates collection element stores for multi-value conversion.
type ICollectionElementStoreFactory interface {
	CreateStore(enumerable any) ces.ICollectionElementStore[string, any]
}

// CollectionElementStoreFactory creates the appropriate collection element store for a given enumerable.
type CollectionElementStoreFactory struct{}

func NewCollectionElementStoreFactory() *CollectionElementStoreFactory {
	return &CollectionElementStoreFactory{}
}

func (this *CollectionElementStoreFactory) CreateStore(enumerable any) ces.ICollectionElementStore[string, any] {
	t := reflect.TypeOf(enumerable)
	if t == nil {
		return ces.NewEnumerableElementStore([]any{})
	}
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		slice := enumerable
		if t.Kind() == reflect.Array {
			s := reflect.ValueOf(enumerable).Slice(0, reflect.ValueOf(enumerable).Len()).Interface()
			slice = s
		}
		return ces.NewListElementStore(slice)
	default:
		return ces.NewEnumerableElementStore(enumerable)
	}
}
