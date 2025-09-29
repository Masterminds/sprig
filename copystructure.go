package sprig

import (
	"fmt"
	"reflect"
)

// deepcopy performs a deep copy of the given interface{}.
func deepcopy(src interface{}) (interface{}, error) {
	if src == nil {
		return make(map[string]interface{}), nil
	}
	return copyValue(reflect.ValueOf(src))
}

// copyValue handles copying using reflection for non-map types
func copyValue(original reflect.Value) (interface{}, error) {
	switch original.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128, reflect.String, reflect.Array:
		return original.Interface(), nil

	case reflect.Interface:
		if original.IsNil() {
			return original.Interface(), nil
		}
		return copyValue(original.Elem())

	case reflect.Map:
		if original.IsNil() {
			return original.Interface(), nil
		}
		copied := reflect.MakeMap(original.Type())

		var err error
		var child interface{}
		iter := original.MapRange()
		for iter.Next() {
			key := iter.Key()
			value := iter.Value()

			if value.Kind() == reflect.Interface && value.IsNil() {
				copied.SetMapIndex(key, value)
				continue
			}

			child, err = copyValue(value)
			if err != nil {
				return nil, err
			}
			copied.SetMapIndex(key, reflect.ValueOf(child))
		}
		return copied.Interface(), nil

	case reflect.Pointer:
		if original.IsNil() {
			return original.Interface(), nil
		}
		copied, err := copyValue(original.Elem())
		if err != nil {
			return nil, err
		}
		ptr := reflect.New(original.Type().Elem())
		ptr.Elem().Set(reflect.ValueOf(copied))
		return ptr.Interface(), nil

	case reflect.Slice:
		if original.IsNil() {
			return original.Interface(), nil
		}
		copied := reflect.MakeSlice(original.Type(), original.Len(), original.Cap())
		for i := 0; i < original.Len(); i++ {
			val, err := copyValue(original.Index(i))
			if err != nil {
				return nil, err
			}
			copied.Index(i).Set(reflect.ValueOf(val))
		}
		return copied.Interface(), nil

	case reflect.Struct:
		copied := reflect.New(original.Type()).Elem()
		for i := 0; i < original.NumField(); i++ {
			elem, err := copyValue(original.Field(i))
			if err != nil {
				return nil, err
			}
			copied.Field(i).Set(reflect.ValueOf(elem))
		}
		return copied.Interface(), nil

	case reflect.Func, reflect.Chan, reflect.UnsafePointer:
		if original.IsNil() {
			return original.Interface(), nil
		}
		return original.Interface(), nil

	default:
		return original.Interface(), fmt.Errorf("unsupported type %v", original)
	}
}
