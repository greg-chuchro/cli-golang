package collection_element_stores

import (
	"fmt"

	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
)

// SetElementStore provides key-value access to set elements. For sets, the key and value are the same - the element itself.
type SetElementStore struct {
	*GenericCollectionElementStoreBase[string, any]
}

func NewSetElementStore(set any) *SetElementStore {
	s := &SetElementStore{}
	s.GenericCollectionElementStoreBase = NewGenericCollectionElementStoreBase[string, any](set)
	return s
}

func (this *SetElementStore) Get(key string) any {
	for i := 0; i < this.Collection.Len(); i++ {
		item := this.Collection.Index(i).Interface()
		if equals(item, key) {
			return item
		}
	}
	panic(exceptions.InvalidListKey(key))
}

func (this *SetElementStore) Set(key string, value any) {
	panic(exceptions.OperationNotSupported("Indexed set on Set - use Add/Remove instead"))
}

func (this *SetElementStore) ContainsKey(key string) bool {
	return this.contains(key)
}

func (this *SetElementStore) GetValueOrInitialize(key string) any {
	if !this.contains(key) {
		this.Add(key)
	}
	return this.Get(key)
}

func (this *SetElementStore) Remove(key string) {
	if !this.remove(key) {
		panic(exceptions.InvalidListKey(key))
	}
}

func equals(a any, b string) bool {
	if bs, ok := a.(string); ok {
		return bs == b
	}
	return fmt.Sprintf("%v", a) == b
}
