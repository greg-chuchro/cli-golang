package collection_element_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
	"github.com/greg-chuchro/cli-golang/data_access/extensions"
)

// DictionaryElementStore provides key-value access to dictionary elements.
type DictionaryElementStore struct {
	*DictionaryElementStoreBase[string, any]
}

func NewDictionaryElementStore(dictionary any) *DictionaryElementStore {
	s := &DictionaryElementStore{}
	s.DictionaryElementStoreBase = &DictionaryElementStoreBase[string, any]{}
	s.DictionaryElementStoreBase.Dictionary = addressable(dictionary)
	s.DictionaryElementStoreBase.CollectionElementStoreBase = NewCollectionElementStoreBase[string, any](s)
	return s
}

func (this *DictionaryElementStore) parseKey(key string) reflect.Value {
	keyType := this.Dictionary.Type().Key()
	return reflect.ValueOf(extensions.Parse(keyType, key))
}

func (this *DictionaryElementStore) Get(key string) any {
	parsedKey := this.parseKey(key)
	return this.Dictionary.MapIndex(parsedKey).Interface()
}

func (this *DictionaryElementStore) Set(key string, value any) {
	parsedKey := this.parseKey(key)
	this.Dictionary.SetMapIndex(parsedKey, reflect.ValueOf(value))
}

func (this *DictionaryElementStore) ContainsKey(key string) bool {
	parsedKey := this.parseKey(key)
	return this.Dictionary.MapIndex(parsedKey).IsValid()
}

func (this *DictionaryElementStore) GetValueOrInitialize(key string) any {
	parsedKey := this.parseKey(key)
	element := this.Dictionary.MapIndex(parsedKey)
	if !element.IsValid() {
		element = reflect.ValueOf(extensions.CreateInstance(this.Dictionary.Type().Elem()))
		this.Dictionary.SetMapIndex(parsedKey, element)
	}
	return element.Interface()
}

var _ = exceptions.InvalidListKey
