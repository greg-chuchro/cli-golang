package extensions

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
)

// interfaceToConcreteMap maps collection interface types to their concrete implementation types.
var interfaceToConcreteMap = map[reflect.Type]reflect.Type{
	reflect.TypeOf((*[]any)(nil)).Elem():          reflect.TypeOf([]any(nil)),
	reflect.TypeOf((*map[string]any)(nil)).Elem(): reflect.TypeOf(map[string]any(nil)),
}

// IsEnumerable determines whether the specified type is an enumerable type.
func IsEnumerable(t reflect.Type) bool {
	if t == nil {
		return false
	}
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array || t.Kind() == reflect.Map {
		return true
	}
	return false
}

// GetEnumerableElementType gets the element type for an enumerable, or the type itself if not an enumerable.
func GetEnumerableElementType(t reflect.Type) reflect.Type {
	if !IsEnumerable(t) {
		return t
	}
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		return t.Elem()
	case reflect.Map:
		return t.Elem()
	case reflect.Pointer:
		return GetEnumerableElementType(t.Elem())
	default:
		return reflect.TypeOf((*any)(nil)).Elem()
	}
}

// CreateInstance creates an instance of the specified type, with support for collection interfaces.
func CreateInstance(t reflect.Type) any {
	if t == nil {
		panic(exceptions.DataAccessErrors.CannotCreateInstance("nil"))
	}

	switch t.Kind() {
	case reflect.Bool:
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uint64(0)
	case reflect.Float32, reflect.Float64:
		return float64(0)
	case reflect.Complex64, reflect.Complex128:
		return complex128(0)
	case reflect.String:
		return ""
	case reflect.Pointer, reflect.Interface:
		if t.Kind() == reflect.Pointer {
			elem := t.Elem()
			if elem.Kind() == reflect.Struct || elem.Kind() == reflect.Ptr {
				return reflect.New(elem).Interface()
			}
		}
		if concrete := GetConcreteCollectionType(t); concrete != nil {
			return reflect.New(concrete).Elem().Interface()
		}
		return reflect.Zero(t).Interface()
	case reflect.Struct:
		return reflect.New(t).Elem().Interface()
	case reflect.Slice:
		return reflect.MakeSlice(t, 0, 0).Interface()
	case reflect.Map:
		return reflect.MakeMap(t).Interface()
	case reflect.Array:
		return reflect.New(t).Elem().Interface()
	default:
		if concrete := GetConcreteCollectionType(t); concrete != nil {
			return reflect.New(concrete).Elem().Interface()
		}
		panic(exceptions.DataAccessErrors.CannotCreateInstance(typeName(t)))
	}
}

// GetConcreteCollectionType resolves a concrete collection type for an interface type.
func GetConcreteCollectionType(interfaceType reflect.Type) reflect.Type {
	if interfaceType == nil {
		return nil
	}
	if concrete, ok := interfaceToConcreteMap[interfaceType]; ok {
		return concrete
	}
	return nil
}

// IsParsable determines whether the specified type can be parsed from a string.
func IsParsable(t reflect.Type) bool {
	if t == nil {
		return false
	}
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		return true
	default:
		return false
	}
}

// Parse parses the string value into an instance of the given type.
func Parse(t reflect.Type, value string) any {
	if t == nil {
		panic(exceptions.DataAccessErrors.TypeCannotBeParsed("nil"))
	}
	elem := t
	for elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
	}
	switch elem.Kind() {
	case reflect.Bool:
		return parseBool(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return parseInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return parseUint(value)
	case reflect.Float32, reflect.Float64:
		return parseFloat(value)
	case reflect.String:
		return value
	default:
		panic(exceptions.DataAccessErrors.TypeCannotBeParsed(typeName(t)))
	}
}

func parseBool(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1", "t", "yes":
		return true
	case "false", "0", "f", "no", "":
		return false
	default:
		return false
	}
}

func parseInt(value string) int64 {
	var result int64
	n, err := fmt.Sscanf(value, "%d", &result)
	if err != nil || n != 1 {
		panic(exceptions.DataAccessErrors.TypeCannotBeParsed(value))
	}
	return result
}

func parseUint(value string) uint64 {
	var result uint64
	n, err := fmt.Sscanf(value, "%d", &result)
	if err != nil || n != 1 {
		panic(exceptions.DataAccessErrors.TypeCannotBeParsed(value))
	}
	return result
}

func parseFloat(value string) float64 {
	var result float64
	n, err := fmt.Sscanf(value, "%f", &result)
	if err != nil || n != 1 {
		panic(exceptions.DataAccessErrors.TypeCannotBeParsed(value))
	}
	return result
}

func typeName(t reflect.Type) string {
	if t == nil {
		return "nil"
	}
	return t.String()
}
