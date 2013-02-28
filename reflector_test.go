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

func ExampleMapToStruct() {
	var s S
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13,
	}
	MapToStruct(m, &s, "json")
	fmt.Printf("%+v", s)
	// Output:
	// {Int:42 Uint8:8 Uintptr:195939070 Float32:3.14 String:str foo:0}
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
	var s []S
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13,
	}
	MapsToStructs([]map[string]interface{}{m}, &s, "json")
	fmt.Printf("%+v", s)
	// Output:
	// [{Int:42 Uint8:8 Uintptr:195939070 Float32:3.14 String:str foo:0}]
}

func BenchmarkMapsToStructs(b *testing.B) {
	var s []S
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapsToStructs([]map[string]interface{}{m}, &s, "json")
	}
	b.StopTimer()
	expected := []S{{42, 8, 0xbadcafe, 3.14, "str", 0}}
	if !reflect.DeepEqual(expected, s) {
		b.Fatalf("Expected %#v, got %#v", expected, s)
	}
}

func BenchmarkMapsToStructs2(b *testing.B) {
	var s []S
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapsToStructs2([]interface{}{m}, &s, "json")
	}
	b.StopTimer()
	expected := []S{{42, 8, 0xbadcafe, 3.14, "str", 0}}
	if !reflect.DeepEqual(expected, s) {
		b.Fatalf("Expected %#v, got %#v", expected, s)
	}
}
