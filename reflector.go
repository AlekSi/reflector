package reflector

import (
	"fmt"
	"reflect"
)

func MapToStruct(m map[string]interface{}, structPointer interface{}, tag string) {
	structPointerType := reflect.TypeOf(structPointer)
	if structPointerType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("Expected struct pointer as second argument, got %s", structPointerType.Kind()))
	}

	structType := structPointerType.Elem()
	if structType.Kind() != reflect.Struct {
		panic(fmt.Errorf("Expected struct pointer as second argument, got pointer to %s", structType.Kind()))
	}
	s := reflect.ValueOf(structPointer).Elem()

	for i := 0; i < structType.NumField(); i++ {
		f := s.Field(i)
		if !f.CanSet() {
			continue
		}

		var name string
		stf := structType.Field(i)
		if tag != "" {
			name = stf.Tag.Get(tag)
		}
		if name == "" {
			name = stf.Name
		}
		v, ok := m[name]
		if !ok {
			continue
		}

		switch f.Kind() {
		case reflect.Bool:
			f.SetBool(v.(bool))

		case reflect.Int:
			f.SetInt(int64(v.(int)))
		case reflect.Int8:
			f.SetInt(int64(v.(int8)))
		case reflect.Int16:
			f.SetInt(int64(v.(int16)))
		case reflect.Int32:
			f.SetInt(int64(v.(int32)))
		case reflect.Int64:
			f.SetInt(v.(int64))

		case reflect.Uint:
			f.SetUint(uint64(v.(uint)))
		case reflect.Uint8:
			f.SetUint(uint64(v.(uint8)))
		case reflect.Uint16:
			f.SetUint(uint64(v.(uint16)))
		case reflect.Uint32:
			f.SetUint(uint64(v.(uint32)))
		case reflect.Uint64:
			f.SetUint(v.(uint64))
		case reflect.Uintptr:
			f.SetUint(uint64(v.(uintptr)))

		case reflect.Float32:
			f.SetFloat(float64(v.(float32)))
		case reflect.Float64:
			f.SetFloat(v.(float64))

		case reflect.String:
			f.SetString(v.(string))

		default:
			// not implemented
		}
	}
	return
}
