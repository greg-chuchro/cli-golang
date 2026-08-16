package data_access

import (
	"fmt"
	"reflect"

	"github.com/greg-chuchro/cli-golang/data_access/extensions"
)

// ReadOnlyPassThroughConverter is a pass-through converter for read-only contexts (submodules).
type ReadOnlyPassThroughConverter[TValue any] struct{}

func NewReadOnlyPassThroughConverter[TValue any]() *ReadOnlyPassThroughConverter[TValue] {
	return &ReadOnlyPassThroughConverter[TValue]{}
}

func (this *ReadOnlyPassThroughConverter[TValue]) CanConvert(targetType reflect.Type) bool {
	return true
}

func (this *ReadOnlyPassThroughConverter[TValue]) ConvertFrom(value any, targetType reflect.Type) TValue {
	return value.(TValue)
}

func (this *ReadOnlyPassThroughConverter[TValue]) ConvertToOrUpdate(value TValue, targetType reflect.Type, currentValue any) any {
	panic("not supported")
}

func (this *ReadOnlyPassThroughConverter[TValue]) IsMultiValue(t reflect.Type) bool {
	return t.Kind() == reflect.Slice || t.Kind() == reflect.Array || t.Kind() == reflect.Map
}

func (this *ReadOnlyPassThroughConverter[TValue]) GetValueType(t reflect.Type) reflect.Type {
	return extensions.GetEnumerableElementType(t)
}

// CompositeValueConverter decorates a base value converter with multi-value (enumerable) handling.
type CompositeValueConverter[TValue any] struct {
	baseConverter            IValueConverter[TValue]
	collectionElementFactory ICollectionElementStoreFactory
}

func NewCompositeValueConverter[TValue any](baseConverter IValueConverter[TValue], collectionElementFactory ICollectionElementStoreFactory) *CompositeValueConverter[TValue] {
	return &CompositeValueConverter[TValue]{baseConverter: baseConverter, collectionElementFactory: collectionElementFactory}
}

func (this *CompositeValueConverter[TValue]) CanConvert(targetType reflect.Type) bool {
	return this.baseConverter.CanConvert(targetType) ||
		(this.IsMultiValue(targetType) && this.baseConverter.CanConvert(targetType.Elem()))
}

func (this *CompositeValueConverter[TValue]) ConvertFrom(value any, targetType reflect.Type) TValue {
	if this.IsMultiValue(targetType) {
		if value == nil {
			return *new(TValue)
		}
		return this.convertEnumerableToString(value, targetType)
	}
	return this.baseConverter.ConvertFrom(value, targetType)
}

func (this *CompositeValueConverter[TValue]) ConvertToOrUpdate(value TValue, targetType reflect.Type, currentValue any) any {
	if this.IsMultiValue(targetType) {
		return this.appendToEnumerable(value, targetType, currentValue)
	}
	return this.baseConverter.ConvertToOrUpdate(value, targetType, currentValue)
}

func (this *CompositeValueConverter[TValue]) IsMultiValue(t reflect.Type) bool {
	return !this.baseConverter.CanConvert(t) && (t.Kind() == reflect.Slice || t.Kind() == reflect.Array || t.Kind() == reflect.Map)
}

func (this *CompositeValueConverter[TValue]) GetValueType(t reflect.Type) reflect.Type {
	return extensions.GetEnumerableElementType(t)
}

func (this *CompositeValueConverter[TValue]) convertEnumerableToString(value any, targetType reflect.Type) TValue {
	elementType := targetType.Elem()
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	elements := make([]string, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		elements = append(elements, fmt.Sprintf("%v", this.baseConverter.ConvertFrom(rv.Index(i).Interface(), elementType)))
	}
	return any("[" + joinComma(elements) + "]").(TValue)
}

func (this *CompositeValueConverter[TValue]) appendToEnumerable(value TValue, targetType reflect.Type, currentValue any) any {
	elementType := targetType.Elem()
	item := this.baseConverter.ConvertToOrUpdate(value, elementType, nil)
	store := this.collectionElementFactory.CreateStore(currentValue)
	return store.Append(item)
}

func joinComma(elements []string) string {
	result := ""
	for i, e := range elements {
		if i > 0 {
			result += ", "
		}
		result += e
	}
	return result
}
