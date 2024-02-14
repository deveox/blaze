package decoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
)

type Simple struct {
	Int    int    `json:"int,omitempty"`
	String string `json:"string,omitempty"`
	Bool   bool   `json:"bool,omitempty"`
	Slice  []int  `json:"slice,omitempty"`
}

func newSimpleStruct() []byte {
	return []byte(`{"int":100,"string":"hello","bool":true}`)
}

func TestDecode_Struct_Simple(t *testing.T) {
	s := newSimpleStruct()
	EqualUnmarshaling[Simple](t, s)
}

func BenchmarkStruct_Simple_Blaze(b *testing.B) {
	var simple Simple
	for i := 0; i < b.N; i++ {
		Unmarshal(newSimpleStruct(), &simple)
	}
	b.SetBytes(int64(len(newSimpleStruct())))
}

func BenchmarkStruct_Simple_Std(b *testing.B) {
	var simple Simple
	for i := 0; i < b.N; i++ {
		json.Unmarshal(newSimpleStruct(), &simple)
	}
	b.SetBytes(int64(len(newSimpleStruct())))
}

func BenchmarkStruct_Simple_GoJson(b *testing.B) {
	var simple Simple
	for i := 0; i < b.N; i++ {
		gojson.Unmarshal(newSimpleStruct(), &simple)
	}
	b.SetBytes(int64(len(newSimpleStruct())))
}

func BenchmarkStruct_Simple_JsonIter(b *testing.B) {
	var simple Simple
	for i := 0; i < b.N; i++ {
		jsoniter.Unmarshal(newSimpleStruct(), &simple)
	}
	b.SetBytes(int64(len(newSimpleStruct())))
}

type EmbeddedPrimitive string

type Embedded struct {
	EmbeddedPrimitive `json:"primitive"`
	*Simple
}
type Complex struct {
	*Embedded
	Primitive string     `json:"primitive,omitempty"`
	Nested    *Simple    `json:"nested,omitempty"`
	Array     *[]*Simple `json:"array,omitempty"`
}

func newComplexStruct() []byte {
	b := newSimpleStruct()
	b[len(b)-1] = ','
	b = append(b, []byte(`"primitive":"embedded",`)...)
	b = append(b, []byte(`"nested":`)...)
	b = append(b, newSimpleStruct()...)
	b = append(b, ',')
	b = append(b, []byte(`"array":[`)...)
	for i := 0; i < 10; i++ {
		b = append(b, newSimpleStruct()...)
		b = append(b, ',')
	}
	b[len(b)-1] = ']'
	b = append(b, '}')
	return b
}

func TestDecode_Struct_Complex(t *testing.T) {
	s := newComplexStruct()
	EqualUnmarshaling[Complex](t, s)
}

var benchStructComplex = newComplexStruct()

func BenchmarkStruct_Complex_Blaze(b *testing.B) {
	var complex Complex
	data := newComplexStruct()
	for i := 0; i < b.N; i++ {
		err := Unmarshal(data, &complex)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(newComplexStruct())))
}

func BenchmarkStruct_Complex_Std(b *testing.B) {
	var complex Complex
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(benchStructComplex, &complex)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(benchStructComplex)))
}

func BenchmarkStruct_Complex_GoJson(b *testing.B) {
	var complex Complex
	for i := 0; i < b.N; i++ {
		err := gojson.Unmarshal(benchStructComplex, &complex)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(benchStructComplex)))
}

func BenchmarkStruct_Complex_JsonIter(b *testing.B) {
	var complex Complex
	for i := 0; i < b.N; i++ {
		err := jsoniter.Unmarshal(benchStructComplex, &complex)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(benchStructComplex)))
}
