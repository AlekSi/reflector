// Package reflector extends standard package reflect with useful utilities.
package reflector

import (
	"fmt"
	"reflect"
	"strconv"
)

// Converts value to kind. Panics if it can't be done.
type Converter func(value interface{}, kind reflect.Kind) interface{}

// Converter: requires value to be exactly of specified kind.
func NoConvert(value interface{}, kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Bool:
		return value.(bool)

	case reflect.Int:
		return int64(value.(int))
	case reflect.Int8:
		return int64(value.(int8))
	case reflect.Int16:
		return int64(value.(int16))
	case reflect.Int32:
		return int64(value.(int32))
	case reflect.Int64:
		return value.(int64)

	case reflect.Uint:
		return uint64(value.(uint))
	case reflect.Uint8:
		return uint64(value.(uint8))
	case reflect.Uint16:
		return uint64(value.(uint16))
	case reflect.Uint32:
		return uint64(value.(uint32))
	case reflect.Uint64:
		return value.(uint64)
	case reflect.Uintptr:
		return uint64(value.(uintptr))

	case reflect.Float32:
		return float64(value.(float32))
	case reflect.Float64:
		return value.(float64)

	case reflect.String:
		return value.(string)
	}

	panic(fmt.Errorf("NoConvert: can't convert %#v to %s", value, kind))
}

// Converter: uses strconv.Parse* functions.
func Strconv(value interface{}, kind reflect.Kind) (res interface{}) {
	e := fmt.Errorf("Strconv: can't convert %#v to %s", value, kind)
	s := fmt.Sprint(value)

	switch kind {
	case reflect.Bool:
		res, e = strconv.ParseBool(s)
		if e != nil {
			panic(e)
		}
		return

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res, e = strconv.ParseInt(s, 10, 64)
		if e != nil {
			panic(e)
		}
		return

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		res, e = strconv.ParseUint(s, 10, 64)
		if e != nil {
			panic(e)
		}
		return

	case reflect.Float32, reflect.Float64:
		res, e = strconv.ParseFloat(s, 64)
		if e != nil {
			panic(e)
		}
		return

	case reflect.String:
		return s
	}

	panic(e)
}

// Converts a struct to map.
// First argument is a pointer to struct.
// Second argument is a not-nil map which will be modified.
// Only exported struct fields are used.
// Tag may be used to change mapping between struct field and map key.
// Currently supports bool, ints, uints, floats, strings.
// Panics in case of error.
func StructToMap(structPointer interface{}, m map[string]interface{}, tag string) {
	structPointerType := reflect.TypeOf(structPointer)
	if structPointerType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("StructToMap: expected pointer to struct as first argument, got %s", structPointerType.Kind()))
	}

	structType := structPointerType.Elem()
	if structType.Kind() != reflect.Struct {
		panic(fmt.Errorf("StructToMap: expected pointer to struct as first argument, got pointer to %s", structType.Kind()))
	}

	s := reflect.ValueOf(structPointer).Elem()

	var name string
	for i := 0; i < structType.NumField(); i++ {
		f := s.Field(i)
		if !f.CanSet() {
			continue
		}

		stf := structType.Field(i)
		name = ""
		if tag != "" {
			name = stf.Tag.Get(tag)
		}
		if name == "" {
			name = stf.Name
		}

		m[name] = f.Interface()
	}
}

// Converts a map to struct using converter function.
// First argument is a map.
// Second argument is a not-nil pointer to struct which will be modified.
// Only exported struct fields are set. Omitted or extra values in map are ignored.
// Tag may be used to change mapping between struct field and map key.
// Currently supports bool, ints, uints, floats, strings.
// Panics in case of error.
func MapToStruct(m map[string]interface{}, structPointer interface{}, converter Converter, tag string) {
	structPointerType := reflect.TypeOf(structPointer)
	if structPointerType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("MapToStruct: expected pointer to struct as second argument, got %s", structPointerType.Kind()))
	}

	structType := structPointerType.Elem()
	if structType.Kind() != reflect.Struct {
		panic(fmt.Errorf("MapToStruct: expected pointer to struct as second argument, got pointer to %s", structType.Kind()))
	}
	s := reflect.ValueOf(structPointer).Elem()

	var name string
	defer func() {
		e := recover()
		if e == nil {
			return
		}

		panic(fmt.Errorf("MapToStruct, field %s: %s", name, e))
	}()

	for i := 0; i < structType.NumField(); i++ {
		f := s.Field(i)
		if !f.CanSet() {
			continue
		}

		stf := structType.Field(i)
		name = ""
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

		kind := f.Kind()
		switch kind {
		case reflect.Bool:
			f.SetBool(converter(v, kind).(bool))

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f.SetInt(converter(v, kind).(int64))

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			f.SetUint(converter(v, kind).(uint64))

		case reflect.Float32, reflect.Float64:
			f.SetFloat(converter(v, kind).(float64))

		case reflect.String:
			f.SetString(converter(v, kind).(string))

		default:
			// not implemented
		}
	}

	return
}

// Converts a slice of maps to a slice of structs. Uses MapToStruct().
// First argument is a slice of maps.
// Second argument is a pointer to (possibly nil) slice of structs which will be set.
func MapsToStructs(maps []map[string]interface{}, slicePointer interface{}, converter Converter, tag string) {
	slicePointerType := reflect.TypeOf(slicePointer)
	if slicePointerType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("MapsToStructs: expected pointer to slice of structs as second argument, got %s", slicePointerType.Kind()))
	}

	sliceType := slicePointerType.Elem()
	if sliceType.Kind() != reflect.Slice {
		panic(fmt.Errorf("MapsToStructs: expected pointer to slice of structs as second argument, got pointer to %s", sliceType.Kind()))
	}

	structType := sliceType.Elem()
	if structType.Kind() != reflect.Struct {
		panic(fmt.Errorf("MapsToStructs: expected pointer to slice of structs as second argument, got pointer to slice of %s", structType.Kind()))
	}

	slice := reflect.MakeSlice(sliceType, 0, len(maps))
	for _, m := range maps {
		str := reflect.New(structType)
		MapToStruct(m, str.Interface(), converter, tag)
		slice = reflect.Append(slice, str.Elem())
	}
	reflect.ValueOf(slicePointer).Elem().Set(slice)
}

// Variant of MapsToStructs() with relaxed signature.
func MapsToStructs2(maps []interface{}, slicePointer interface{}, converter Converter, tag string) {
	m := make([]map[string]interface{}, len(maps))
	for index, i := range maps {
		m[index] = i.(map[string]interface{})
	}
	MapsToStructs(m, slicePointer, converter, tag)
}
