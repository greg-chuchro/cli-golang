package data_access

import (
	"reflect"
)

// IValueConverter converts values between CLI string representation and internal object types.
type IValueConverter[TValue any] interface {
	CanConvert(targetType reflect.Type) bool
	ConvertFrom(value any, targetType reflect.Type) TValue
	ConvertToOrUpdate(value TValue, targetType reflect.Type, currentValue any) any
}

// ICompositeValueConverter extends IValueConverter with multi-value metadata capabilities.
type ICompositeValueConverter[TValue any] interface {
	IValueConverter[TValue]
	IsMultiValue(t reflect.Type) bool
	GetValueType(t reflect.Type) reflect.Type
}

// AccessorKeysPair represents a pairing of a reflection accessor with its associated CLI key names.
type AccessorKeysPair[TAccessor any] struct {
	Accessor TAccessor
	Keys     []string
}

func NewAccessorKeysPair[TAccessor any](accessor TAccessor, keys []string) *AccessorKeysPair[TAccessor] {
	return &AccessorKeysPair[TAccessor]{Accessor: accessor, Keys: keys}
}
