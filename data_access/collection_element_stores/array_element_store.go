package collection_element_stores

import (
	"reflect"
	"strconv"

	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
	"github.com/greg-chuchro/cli-golang/data_access/extensions"
)

// ArrayElementStore provides key-value access to array/slice elements by index.
type ArrayElementStore struct {
	*ArrayElementStoreBase[string, any]
}

func NewArrayElementStore(array any) *ArrayElementStore {
	s := &ArrayElementStore{}
	s.ArrayElementStoreBase = &ArrayElementStoreBase[string, any]{}
	s.ArrayElementStoreBase.Array = reflect.ValueOf(array)
	s.ArrayElementStoreBase.CollectionElementStoreBase = NewCollectionElementStoreBase[string, any](s)
	return s
}

func (this *ArrayElementStore) Get(key string) any {
	index, err := strconv.Atoi(key)
	if err != nil {
		panic(exceptions.InvalidListKey("string"))
	}
	return this.Array.Index(index).Interface()
}

func (this *ArrayElementStore) Set(key string, value any) {
	index, err := strconv.Atoi(key)
	if err != nil {
		panic(exceptions.InvalidListKey("string"))
	}
	this.Array.Index(index).Set(reflect.ValueOf(value))
}

func (this *ArrayElementStore) ContainsKey(key string) bool {
	index, err := strconv.Atoi(key)
	if err != nil {
		return false
	}
	return index >= 0 && index < this.Array.Len()
}

func (this *ArrayElementStore) GetValueOrInitialize(key string) any {
	index, err := strconv.Atoi(key)
	if err != nil {
		panic(exceptions.InvalidListKey("string"))
	}
	element := this.Array.Index(index).Interface()
	if element == nil {
		element = extensions.CreateInstance(this.Array.Type().Elem())
		this.Array.Index(index).Set(reflect.ValueOf(element))
	}
	return element
}
