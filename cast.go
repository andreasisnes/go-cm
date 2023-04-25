package configurationmanager

import (
	"reflect"

	"github.com/spf13/cast"
)

func castAndTryAssignValue(from any, to any) any {
	toType := reflect.TypeOf(to)
	depth := 0
	for toType.Kind() == reflect.Pointer {
		toType = toType.Elem()
		depth += 1
	}

	value := castValue(toType.Kind(), from)
	for i := 0; i < depth-1; i++ {
		value = pointer(value)
	}

	toValue := reflect.ValueOf(to)
	if toValue.Kind() == reflect.Pointer {
		toValue = toValue.Elem()
		if toValue.Kind() == reflect.Pointer {
			toValue.SetPointer(ptr)
		} else {
			toValue.Set(reflect.ValueOf(value))
		}
	}

	return value
}

func castValue(tType reflect.Kind, from any) (res any) {
	switch tType {
	case reflect.String:
		if res, ok := from.(string); ok {
			return res
		}
		return cast.ToString(from)
	case reflect.Int:
		if res, ok := from.(int); ok {
			return res
		}
		return cast.ToInt(from)
	case reflect.Int64:
		if res, ok := from.(int64); ok {
			return res
		}
		return cast.ToInt64(from)
	case reflect.Int32:
		if res, ok := from.(int32); ok {
			return res
		}
		return cast.ToInt32(from)
	case reflect.Int16:
		if res, ok := from.(int16); ok {
			return res
		}
		return cast.ToInt16(from)
	case reflect.Int8:
		if res, ok := from.(int8); ok {
			return res
		}
		return cast.ToInt8(from)
	case reflect.Uint:
		if res, ok := from.(uint); ok {
			return res
		}
		return cast.ToUint(from)
	case reflect.Uint64:
		if res, ok := from.(uint64); ok {
			return res
		}
		return cast.ToUint16(from)
	case reflect.Uint32:
		if res, ok := from.(uint32); ok {
			return res
		}
		return cast.ToUint32(from)
	case reflect.Uint16:
		if res, ok := from.(uint16); ok {
			return res
		}
		return cast.ToUint16(from)
	case reflect.Uint8:
		if res, ok := from.(uint8); ok {
			return res
		}
		return cast.ToUint8(from)
	case reflect.Float64:
		if res, ok := from.(float64); ok {
			return res
		}
		return cast.ToFloat64(from)
	case reflect.Float32:
		if res, ok := from.(float32); ok {
			return res
		}
		return cast.ToFloat32(from)
	case reflect.Bool:
		if res, ok := from.(bool); ok {
			return res
		}
		return cast.ToBool(from)
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
		return from
	}
}

func pointer[T any](value T) *T {
	return &value
}
