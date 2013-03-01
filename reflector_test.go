package reflector_test

import (
	. "."
	"fmt"
	"reflect"
	"testing"
)

type S struct {
	Int     int
	Uint8   uint8
	Uintptr uintptr
	Float32 float32 `json:"f32"`
	String  string
	foo     int
}

func TestMapToStructBad1(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "Expected pointer to struct as second argument, got struct" {
			t.Error(r)
		}
	}()
	MapToStruct(map[string]interface{}{}, S{}, "")
	t.Fatal("should panic")
}

func TestMapToStructBad2(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "Expected pointer to struct as second argument, got pointer to int" {
			t.Error(r)
		}
	}()
	var i int
	MapToStruct(map[string]interface{}{}, &i, "")
	t.Fatal("should panic")
}

func TestMapToStructWrongType(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "Field Uint8: interface conversion: interface is int, not uint8" {
			t.Error(r)
		}
	}()
	type T struct {
		Uint8 uint8
	}
	var s T
	m := map[string]interface{}{"Uint8": 8}
	MapToStruct(m, &s, "")
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
	MapToStruct(m, &s, "json")
	fmt.Printf("%+v", s)
	// Output:
	// {Uint8:8 Float32:3.14 String: foo:0}
}

func BenchmarkMapToStruct(b *testing.B) {
	var s S
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapToStruct(m, &s, "json")
	}
	b.StopTimer()
	expected := S{42, 8, 0xbadcafe, 3.14, "str", 0}
	if !reflect.DeepEqual(expected, s) {
		b.Fatalf("Expected %#v, got %#v", expected, s)
	}
}

func TestMapsToStructsBad1(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "Expected pointer to slice of structs as second argument, got slice" {
			t.Error(r)
		}
	}()
	var s []S
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13,
	}
	MapsToStructs([]map[string]interface{}{m}, s, "json")
	t.Fatal("should panic")
}

func TestMapsToStructsBad2(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "Expected pointer to slice of structs as second argument, got pointer to int" {
			t.Error(r)
		}
	}()
	var s *int
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13,
	}
	MapsToStructs([]map[string]interface{}{m}, s, "json")
	t.Fatal("should panic")
}

func TestMapsToStructsBad3(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "Expected pointer to slice of structs as second argument, got pointer to slice of int" {
			t.Error(r)
		}
	}()
	var s *[]int
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13,
	}
	MapsToStructs([]map[string]interface{}{m}, s, "json")
	t.Fatal("should panic")
}

func ExampleMapsToStructs() {
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
	MapsToStructs(maps, &s, "json")
	fmt.Printf("%+v\n", s[0])
	fmt.Printf("%+v\n", s[1])
	// Output:
	// {Uint8:8 Float32:0 String: foo:0}
	// {Uint8:0 Float32:3.14 String:str foo:0}
}

func BenchmarkMapsToStructs(b *testing.B) {
	var s []S
	maps := []map[string]interface{}{
		{"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe)},
		{"f32": float32(3.14), "String": "str", "foo": 13},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapsToStructs(maps, &s, "json")
	}
	b.StopTimer()
	expected := []S{{42, 8, 0xbadcafe, 0, "", 0}, {0, 0, 0, 3.14, "str", 0}}
	if !reflect.DeepEqual(expected, s) {
		b.Fatalf("Expected %#v, got %#v", expected, s)
	}
}

func BenchmarkMapsToStructs2(b *testing.B) {
	var s []S
	maps := []interface{}{
		map[string]interface{}{"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe)},
		map[string]interface{}{"f32": float32(3.14), "String": "str", "foo": 13},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapsToStructs2(maps, &s, "json")
	}
	b.StopTimer()
	expected := []S{{42, 8, 0xbadcafe, 0, "", 0}, {0, 0, 0, 3.14, "str", 0}}
	if !reflect.DeepEqual(expected, s) {
		b.Fatalf("Expected %#v, got %#v", expected, s)
	}
}
