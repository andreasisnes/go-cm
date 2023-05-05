package configurationmanager

import (
	"reflect"

	"github.com/spf13/cast"
)

func CastAndAssignValue(value any, to any) {
	value = CastValue(value, to)
	cValue := reflect.ValueOf(value).Convert(reflect.ValueOf(to).Type())
	reflect.ValueOf(to).Elem().Set(cValue.Elem())
}

func CastValue(value any, to any) any {
	switch toType := reflect.TypeOf(to); toType.Kind() {
	case reflect.Pointer:
		toElemType := toType.Elem()
		newValue := reflect.New(toElemType)
		newValue.Elem().Set(reflect.ValueOf(CastValue(value, newValue.Elem().Interface())))
		return newValue.Interface()
	case reflect.String:
		if res, ok := value.(string); ok {
			return res
		}
		return cast.ToString(value)
	case reflect.Int:
		if res, ok := value.(int); ok {
			return res
		}
		return cast.ToInt(value)
	case reflect.Int64:
		if res, ok := value.(int64); ok {
			return res
		}
		return cast.ToInt64(value)
	case reflect.Int32:
		if res, ok := value.(int32); ok {
			return res
		}
		return cast.ToInt32(value)
	case reflect.Int16:
		if res, ok := value.(int16); ok {
			return res
		}
		return cast.ToInt16(value)
	case reflect.Int8:
		if res, ok := value.(int8); ok {
			return res
		}
		return cast.ToInt8(value)
	case reflect.Uint:
		if res, ok := value.(uint); ok {
			return res
		}
		return cast.ToUint(value)
	case reflect.Uint64:
		if res, ok := value.(uint64); ok {
			return res
		}
		return cast.ToUint16(value)
	case reflect.Uint32:
		if res, ok := value.(uint32); ok {
			return res
		}
		return cast.ToUint32(value)
	case reflect.Uint16:
		if res, ok := value.(uint16); ok {
			return res
		}
		return cast.ToUint16(value)
	case reflect.Uint8:
		if res, ok := value.(uint8); ok {
			return res
		}
		return cast.ToUint8(value)
	case reflect.Float64:
		if res, ok := value.(float64); ok {
			return res
		}
		return cast.ToFloat64(value)
	case reflect.Float32:
		if res, ok := value.(float32); ok {
			return res
		}
		return cast.ToFloat32(value)
	case reflect.Bool:
		if res, ok := value.(bool); ok {
			return res
		}
		return cast.ToBool(value)
	// case time.Time:
	// 	if res, ok := from.(time.Time); ok {
	// 		return res
	// 	}
	// 	return cast.ToTime(from)
	// case time.Duration:
	// 	if res, ok := from.(time.Duration); ok {
	// 		return res
	// 	}
	// 	return cast.ToDuration(from)
	default:
		return value
	}
}

func pointer[T any](value T) *T {
	return &value
}
