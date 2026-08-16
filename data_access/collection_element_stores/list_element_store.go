package collection_element_stores

import (
	"reflect"
	"strconv"

	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
	"github.com/greg-chuchro/cli-golang/data_access/extensions"
)

// ListElementStore provides key-value access to list elements by index.
type ListElementStore struct {
	*ListElementStoreBase[string, any]
}

func NewListElementStore(list any) *ListElementStore {
	s := &ListElementStore{}
	s.ListElementStoreBase = &ListElementStoreBase[string, any]{}
	s.ListElementStoreBase.List = addressable(list)
	s.ListElementStoreBase.CollectionElementStoreBase = NewCollectionElementStoreBase[string, any](s)
	return s
}

func (this *ListElementStore) Get(key string) any {
	index, err := strconv.Atoi(key)
	if err != nil {
		panic(exceptions.InvalidListKey("string"))
	}
	return this.List.Index(index).Interface()
}

func (this *ListElementStore) Set(key string, value any) {
	index, err := strconv.Atoi(key)
	if err != nil {
		panic(exceptions.InvalidListKey("string"))
	}
	this.List.Index(index).Set(reflect.ValueOf(value))
}

func (this *ListElementStore) ContainsKey(key string) bool {
	index, err := strconv.Atoi(key)
	if err != nil {
		return false
	}
	return index >= 0 && index < this.List.Len()
}

func (this *ListElementStore) GetValueOrInitialize(key string) any {
	index, err := strconv.Atoi(key)
	if err != nil {
		panic(exceptions.InvalidListKey("string"))
	}
	element := this.List.Index(index).Interface()
	if element == nil {
		element = extensions.CreateInstance(this.List.Type().Elem())
		this.List.Index(index).Set(reflect.ValueOf(element))
	}
	return element
}

func (this *ListElementStore) Remove(key string) {
	index, err := strconv.Atoi(key)
	if err != nil {
		panic(exceptions.InvalidListKey("string"))
	}
	prefix := this.List.Slice(0, index)
	suffix := this.List.Slice(index+1, this.List.Len())
	this.List.Set(reflect.AppendSlice(prefix, suffix))
}
