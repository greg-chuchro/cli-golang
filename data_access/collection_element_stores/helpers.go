package collection_element_stores

import "reflect"

// addressable returns a settable reflect.Value for the given collection.
// If a pointer is supplied, its element is used so that mutation is possible.
func addressable(value any) reflect.Value {
	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	return v
}
