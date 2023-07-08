package configurationmanager

import (
	"reflect"
	"time"

	"github.com/spf13/cast"
)

func CastAndAssignValue(value any, to any) any {
	value = CastValue(value, to)
	cValue := reflect.ValueOf(value).Elem()
	reflect.ValueOf(to).Elem().Set(cValue)
	return value
}

func CastValue(value any, to any) any {
	switch toType := reflect.TypeOf(to); toType.Kind() {
	case reflect.Pointer:
		// Check if packed inside interface{}
		iValue := reflect.Indirect(reflect.ValueOf(to))
		iType := toType.Elem()
		if iValue.Kind() == reflect.Interface {
			iType = iValue.Elem().Type()
		}

		ptr := reflect.New(iType)
		h := ptr.Elem().Interface()
		f := CastValue(value, h)
		v := reflect.ValueOf(f)
		ptr.Elem().Set(v)
		return ptr.Interface()
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
		if _, ok := to.(time.Duration); ok {
			return cast.ToDuration(value)
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
	case reflect.Struct:
		if toType == reflect.TypeOf(time.Time{}) {
			if res, ok := value.(time.Time); ok {
				return res
			}
			return cast.ToTime(value)
		}
		return nil
	default:
		return value
	}
}

func pointer[T any](value T) *T {
	return &value
}
