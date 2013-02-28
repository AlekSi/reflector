package reflector_test

import (
	. "."
	"fmt"
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
		if !ok || e.Error() != "Expected struct pointer as second argument, got struct" {
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
		if !ok || e.Error() != "Expected struct pointer as second argument, got pointer to int" {
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
	m := map[string]interface{}{"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapToStruct(m, &s, "json")
	}
}
