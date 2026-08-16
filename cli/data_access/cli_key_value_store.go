package data_access

import (
	"reflect"

	da "github.com/greg-chuchro/cli-golang/data_access"
)

// ICliReadOnlyKeyValueStore provides read-only access to a CLI key-value store with accessor grouping.
type ICliReadOnlyKeyValueStore[TValue any] interface {
	da.IReadOnlyKeyValueStore[string, TValue]
	GetAccessorKeysPairs() []AccessorKeysPair[reflect.StructField]
}

// ICliKeyValueStore provides read-write access to a CLI key-value store with accessor grouping.
type ICliKeyValueStore[TValue any] interface {
	da.IKeyValueStore[string, TValue]
	GetAccessorKeysPairs() []AccessorKeysPair[reflect.StructField]
	IsMultiValue(key string) bool
}

// ICliFunctionExecutor provides a function executor for CLI with accessor grouping.
type ICliFunctionExecutor[TValue any] interface {
	da.IFunctionExecutor[string, TValue]
	ICliKeyValueStore[TValue]
	SetParameterValues()
}
