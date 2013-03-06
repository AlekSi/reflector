package reflector_test

import (
	. "."
	"fmt"
)

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

func ExampleStructValueToMap() {
	type T struct {
		Uint8   uint8
		Float32 float32 `json:"f32"` // tag will be used
		String  string
		foo     int // not exported
	}
	s := T{8, 3.14, "str", 13}
	m := make(map[string]interface{})
	StructValueToMap(s, m, "json")
	fmt.Printf("%#v %#v %#v %#v", m["Uint8"], m["f32"], m["String"], m["foo"])
	// Output:
	// 0x8 3.14 "str" <nil>
}

func ExampleStructsToMaps() {
	type T struct {
		Uint8   uint8
		Float32 float32 `json:"f32"` // tag will be used
		String  string
		foo     int // not exported
	}
	s := []T{{8, 3.14, "str", 13}}
	var m []map[string]interface{}
	StructsToMaps(s, &m, "json")
	fmt.Printf("%#v %#v %#v %#v", m[0]["Uint8"], m[0]["f32"], m[0]["String"], m[0]["foo"])
	// Output:
	// 0x8 3.14 "str" <nil>
}

func ExampleMapToStructNoConvert() {
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
		Uint8   uint8   // type conversion via strconv
		Float32 float32 `json:"f32"` // tag will be used
		String  string  // not present in first map, will not be set
		foo     int     // not exported, will not be set
	}
	var s []T
	maps := []map[string]interface{}{
		{"Uint8": 8, "f32": 3, "foo": 13},
		{"Uint8": "9", "f32": "4", "String": "43", "foo": "13"},
	}
	MapsToStructs(maps, &s, Strconv, "json")
	fmt.Printf("%+v\n", s[0])
	fmt.Printf("%+v\n", s[1])
	// Output:
	// {Uint8:8 Float32:3 String: foo:0}
	// {Uint8:9 Float32:4 String:43 foo:0}
}
