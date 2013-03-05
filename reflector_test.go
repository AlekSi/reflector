package reflector_test

import (
	. "."
	"fmt"
	"reflect"
	"testing"
)

type T struct {
	Int     int
	Uint8   uint8
	Uintptr uintptr
	Float32 float32 `json:"f32"`
	String  string
	foo     int
}

func TestStructToMapBad1(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "StructToMap: expected pointer to struct as first argument, got struct" {
			t.Error(r)
		}
	}()
	m := make(map[string]interface{})
	StructToMap(T{}, m, "")
	t.Fatal("should panic")
}

func TestStructToMapBad2(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "StructToMap: expected pointer to struct as first argument, got pointer to int" {
			t.Error(r)
		}
	}()
	var i int
	m := make(map[string]interface{})
	StructToMap(&i, m, "")
	t.Fatal("should panic")
}

func ExampleStructToMap() {
	type T struct {
		Uint8   uint8
		Float32 float32 `json:"f32"` // tag will be used
		String  string
		foo     int // not exported
	}
	s := T{8, 3.14, "str", 13}
	m := make(map[string]interface{})
	StructToMap(&s, m, "json")
	fmt.Printf("%#v %#v %#v %#v", m["Uint8"], m["f32"], m["String"], m["foo"])
	// Output:
	// 0x8 3.14 "str" <nil>
}

func BenchmarkStructToMap(b *testing.B) {
	s := T{42, 8, 0xbadcafe, 3.14, "str", 13}
	m := make(map[string]interface{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StructToMap(&s, m, "json")
	}
	b.StopTimer()
	expected := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str",
	}
	if !reflect.DeepEqual(expected, m) {
		b.Fatalf("%#v\n%#v", expected, m)
	}
}

func TestMapToStructBad1(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "MapToStruct: expected pointer to struct as second argument, got struct" {
			t.Error(r)
		}
	}()
	MapToStruct(map[string]interface{}{}, T{}, NoConvert, "")
	t.Fatal("should panic")
}

func TestMapToStructBad2(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "MapToStruct: expected pointer to struct as second argument, got pointer to int" {
			t.Error(r)
		}
	}()
	var i int
	MapToStruct(map[string]interface{}{}, &i, NoConvert, "")
	t.Fatal("should panic")
}

func TestMapToStructWrongType(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "MapToStruct, field Uint8: interface conversion: interface is int, not uint8" {
			t.Error(r)
		}
	}()
	type T struct {
		Uint8 uint8
	}
	var s T
	m := map[string]interface{}{"Uint8": 8}
	MapToStruct(m, &s, NoConvert, "")
	t.Fatal("should panic")
}

func ExampleMapToStruct() {
	type T struct {
		Uint8   uint8   // no automatic type conversion
		Float32 float32 `json:"f32"` // tag will be used
		String  string  // not present in map, will not be set
		foo     int     // not exported, will not be set
	}
	var s T
	m := map[string]interface{}{"Uint8": uint8(8), "f32": float32(3.14), "foo": 13}
	MapToStruct(m, &s, NoConvert, "json")
	fmt.Printf("%+v", s)
	// Output:
	// {Uint8:8 Float32:3.14 String: foo:0}
}

func BenchmarkMapToStruct(b *testing.B) {
	var s T
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapToStruct(m, &s, NoConvert, "json")
	}
	b.StopTimer()
	expected := T{42, 8, 0xbadcafe, 3.14, "str", 0}
	if !reflect.DeepEqual(expected, s) {
		b.Fatalf("Expected %#v, got %#v", expected, s)
	}
}

func TestMapsToStructsBad1(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "MapsToStructs: expected pointer to slice of structs as second argument, got slice" {
			t.Error(r)
		}
	}()
	var s []T
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13,
	}
	MapsToStructs([]map[string]interface{}{m}, s, NoConvert, "json")
	t.Fatal("should panic")
}

func TestMapsToStructsBad2(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "MapsToStructs: expected pointer to slice of structs as second argument, got pointer to int" {
			t.Error(r)
		}
	}()
	var s *int
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13,
	}
	MapsToStructs([]map[string]interface{}{m}, s, NoConvert, "json")
	t.Fatal("should panic")
}

func TestMapsToStructsBad3(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "MapsToStructs: expected pointer to slice of structs as second argument, got pointer to slice of int" {
			t.Error(r)
		}
	}()
	var s *[]int
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13,
	}
	MapsToStructs([]map[string]interface{}{m}, s, NoConvert, "json")
	t.Fatal("should panic")
}

func ExampleMapsToStructsNoConvert() {
	type T struct {
		Uint8   uint8   // no automatic type conversion
		Float32 float32 `json:"f32"` // tag will be used
		String  string  // not present in first map, will not be set
		foo     int     // not exported, will not be set
	}
	var s []T
	maps := []map[string]interface{}{
		{"Uint8": uint8(8)},
		{"f32": float32(3.14), "String": "str", "foo": 13},
	}
	MapsToStructs(maps, &s, NoConvert, "json")
	fmt.Printf("%+v\n", s[0])
	fmt.Printf("%+v\n", s[1])
	// Output:
	// {Uint8:8 Float32:0 String: foo:0}
	// {Uint8:0 Float32:3.14 String:str foo:0}
}

func ExampleMapsToStructsStrconv() {
	type T struct {
		Uint8   uint8
		Float32 float32 `json:"f32"`
		String  string
		foo     int
	}
	var s []T
	maps := []map[string]interface{}{
		{"Uint8": 8, "f32": 3, "String": 42, "foo": 13},
		{"Uint8": "9", "f32": "4", "String": "43", "foo": "13"},
	}
	MapsToStructs(maps, &s, Strconv, "json")
	fmt.Printf("%+v\n", s[0])
	fmt.Printf("%+v\n", s[1])
	// Output:
	// {Uint8:8 Float32:3 String:42 foo:0}
	// {Uint8:9 Float32:4 String:43 foo:0}
}

func BenchmarkMapsToStructs(b *testing.B) {
	var s []T
	maps := []map[string]interface{}{
		{"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe)},
		{"f32": float32(3.14), "String": "str", "foo": 13},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapsToStructs(maps, &s, NoConvert, "json")
	}
	b.StopTimer()
	expected := []T{{42, 8, 0xbadcafe, 0, "", 0}, {0, 0, 0, 3.14, "str", 0}}
	if !reflect.DeepEqual(expected, s) {
		b.Fatalf("Expected %#v, got %#v", expected, s)
	}
}

func BenchmarkMapsToStructs2(b *testing.B) {
	var s []T
	maps := []interface{}{
		map[string]interface{}{"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe)},
		map[string]interface{}{"f32": float32(3.14), "String": "str", "foo": 13},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapsToStructs2(maps, &s, NoConvert, "json")
	}
	b.StopTimer()
	expected := []T{{42, 8, 0xbadcafe, 0, "", 0}, {0, 0, 0, 3.14, "str", 0}}
	if !reflect.DeepEqual(expected, s) {
		b.Fatalf("Expected %#v, got %#v", expected, s)
	}
}
