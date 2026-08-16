package collection_element_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
	"github.com/greg-chuchro/cli-golang/data_access/extensions"
)

// GenericCollectionElementStoreBase is the base for collection element stores that use reflection to access
// collection methods (Add, Remove, Contains, Count). Suitable for collections implementing those members.
type GenericCollectionElementStoreBase[TKey comparable, TValue any] struct {
	*CollectionElementStoreBase[TKey, TValue]
	Collection     reflect.Value
	addMethod      reflect.Value
	removeMethod   reflect.Value
	containsMethod reflect.Value
	countProperty  reflect.Value
}

func NewGenericCollectionElementStoreBase[TKey comparable, TValue any](collection any) *GenericCollectionElementStoreBase[TKey, TValue] {
	s := &GenericCollectionElementStoreBase[TKey, TValue]{}
	s.Collection = reflect.ValueOf(collection)
	collectionType := s.Collection.Type()
	if collectionType.Kind() == reflect.Pointer {
		collectionType = collectionType.Elem()
	}
	s.addMethod = s.Collection.MethodByName("Add")
	s.removeMethod = s.Collection.MethodByName("Remove")
	s.containsMethod = s.Collection.MethodByName("Contains")
	s.countProperty = s.Collection.FieldByName("Count")
	s.CollectionElementStoreBase = NewCollectionElementStoreBase[TKey, TValue](s)
	return s
}

func (this *GenericCollectionElementStoreBase[TKey, TValue]) GetValueType(key TKey) reflect.Type {
	return this.Collection.Type().Elem()
}

func (this *GenericCollectionElementStoreBase[TKey, TValue]) Add(value TValue) {
	if this.addMethod.IsValid() {
		this.addMethod.Call([]reflect.Value{reflect.ValueOf(value)})
	}
}

func (this *GenericCollectionElementStoreBase[TKey, TValue]) Append(value TValue) any {
	this.Add(value)
	return this.Collection.Interface()
}

func (this *GenericCollectionElementStoreBase[TKey, TValue]) AddNew() TValue {
	var zero TKey
	elementType := this.GetValueType(zero)
	newElement := extensions.CreateInstance(elementType)
	this.Add(newElement.(TValue))
	return newElement.(TValue)
}

func (this *GenericCollectionElementStoreBase[TKey, TValue]) contains(item any) bool {
	if this.containsMethod.IsValid() && item != nil {
		result := this.containsMethod.Call([]reflect.Value{reflect.ValueOf(item)})
		return len(result) > 0 && result[0].Bool()
	}
	return false
}

func (this *GenericCollectionElementStoreBase[TKey, TValue]) remove(item any) bool {
	if this.removeMethod.IsValid() && item != nil {
		result := this.removeMethod.Call([]reflect.Value{reflect.ValueOf(item)})
		return len(result) > 0 && result[0].Bool()
	}
	return false
}

var _ = exceptions.OperationNotSupported
