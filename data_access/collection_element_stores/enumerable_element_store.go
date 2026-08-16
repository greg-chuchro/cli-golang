package collection_element_stores

import (
	"reflect"
	"strconv"

	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
	"github.com/greg-chuchro/cli-golang/data_access/extensions"
)

// EnumerableElementStore provides element store for slices using immutable append operations.
type EnumerableElementStore struct {
	*CollectionElementStoreBase[string, any]
	TargetSlice []any
}

func NewEnumerableElementStore(enumerable any) *EnumerableElementStore {
	s := &EnumerableElementStore{}
	v := reflect.ValueOf(enumerable)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			s.TargetSlice = append(s.TargetSlice, v.Index(i).Interface())
		}
	}
	s.CollectionElementStoreBase = NewCollectionElementStoreBase[string, any](s)
	return s
}

func (this *EnumerableElementStore) Get(key string) any {
	index, err := strconv.Atoi(key)
	if err != nil {
		panic(exceptions.InvalidListKey("string"))
	}
	return this.TargetSlice[index]
}

func (this *EnumerableElementStore) Set(key string, value any) {
	panic(exceptions.OperationNotSupported("Indexed set"))
}

func (this *EnumerableElementStore) ContainsKey(key string) bool {
	index, err := strconv.Atoi(key)
	if err != nil {
		return false
	}
	return index >= 0 && index < len(this.TargetSlice)
}

func (this *EnumerableElementStore) GetValueType(key string) reflect.Type {
	if len(this.TargetSlice) > 0 {
		return reflect.TypeOf(this.TargetSlice[0])
	}
	return reflect.TypeOf((*any)(nil)).Elem()
}

func (this *EnumerableElementStore) GetValueOrInitialize(key string) any {
	index, err := strconv.Atoi(key)
	if err != nil {
		panic(exceptions.InvalidListKey("string"))
	}
	if index < len(this.TargetSlice) && this.TargetSlice[index] != nil {
		return this.TargetSlice[index]
	}
	element := extensions.CreateInstance(this.GetValueType(key))
	return element
}

func (this *EnumerableElementStore) Append(value any) any {
	this.TargetSlice = append(this.TargetSlice, value)
	return this.TargetSlice
}
