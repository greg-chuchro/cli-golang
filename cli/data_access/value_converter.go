package data_access

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/greg-chuchro/cli-golang/cli/exceptions"
	"github.com/greg-chuchro/cli-golang/data_access/extensions"
)

// ValueConverter converts values between CLI string representation and internal object types.
type ValueConverter struct {
	FormatProvider string
}

func NewValueConverter() *ValueConverter {
	return &ValueConverter{FormatProvider: "R"}
}

func (this *ValueConverter) CanConvert(targetType reflect.Type) bool {
	if targetType == nil {
		return false
	}
	underlying := targetType
	if targetType.Kind() == reflect.Pointer {
		underlying = targetType.Elem()
	}
	return extensions.IsParsable(underlying) || underlying.Kind() == reflect.String
}

func (this *ValueConverter) ConvertFrom(value any, targetType reflect.Type) any {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

func (this *ValueConverter) ConvertToOrUpdate(value any, targetType reflect.Type, currentValue any) any {
	if value == nil {
		return nil
	}

	s, ok := value.(string)
	if !ok {
		s = fmt.Sprintf("%v", value)
	}

	underlying := targetType
	if targetType.Kind() == reflect.Pointer {
		underlying = targetType.Elem()
	}

	if underlying.Kind() == reflect.String {
		return s
	}

	tryParse := func(t reflect.Type) (any, bool) {
		switch t.Kind() {
		case reflect.Bool:
			b, err := strconv.ParseBool(strings.ToLower(s))
			return b, err == nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(s, 10, 64)
			return reflect.ValueOf(n).Convert(t).Interface(), err == nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, err := strconv.ParseUint(s, 10, 64)
			return reflect.ValueOf(n).Convert(t).Interface(), err == nil
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(s, 64)
			return reflect.ValueOf(f).Convert(t).Interface(), err == nil
		default:
			return nil, false
		}
	}

	if parsed, ok := tryParse(underlying); ok {
		if targetType.Kind() == reflect.Pointer {
			pv := reflect.New(targetType)
			pv.Elem().Set(reflect.ValueOf(parsed))
			return pv.Interface()
		}
		return parsed
	}

	panic(exceptions.CliErrors.InvalidValueFormat(targetType.String(), fmt.Errorf("cannot convert %q", value)))
}
