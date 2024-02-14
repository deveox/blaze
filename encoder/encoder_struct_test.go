package encoder

import (
	"encoding/json"
	"testing"
	"time"

	gojson "github.com/goccy/go-json"
)

type DataInterface interface {
	InterfaceMethod()
}
type EmbeddedPrimitive string

func newData(depth int, sliceLen int, nested bool) *Data {
	res := &Data{
		String:            "test",
		Int16:             42,
		Float32:           3.14,
		Bool:              true,
		Time:              time.Now(),
		TimePtr:           nil,
		EmbeddedPrimitive: "embedded",
	}
	if nested {
		nested := newData(0, 0, false)
		res.Nested = &nested
		res.IfacePtr = nested
	}
	for i := 0; i < sliceLen; i++ {
		res.Slice = append(res.Slice, newData(0, 0, false))
	}
	if depth > 0 {
		res.Data = newData(depth-1, sliceLen, false)
	}
	return res
}

type Data struct {
	String            string        `json:"string" blaze:"keep"`
	Int16             int16         `json:"int16" blaze:"keep"`
	Float32           float32       `json:"float32" blaze:"keep"`
	Bool              bool          `json:"bool" blaze:"keep"`
	Time              time.Time     `json:"time" blaze:"keep"`
	TimePtr           *time.Time    `json:"timePtr" blaze:"keep"`
	Slice             []*Data       `json:"slice" blaze:"keep"`
	Nested            **Data        `json:"nested" blaze:"keep"`
	IfacePtr          DataInterface `json:"ifacePtr" blaze:"keep"`
	EmbeddedPrimitive `json:"embeddedPrimitive" blaze:"keep"`
	*Data             `json:"embedded" blaze:"keep"`
}

func (d *Data) InterfaceMethod() {}

func newDataEmpty(depth int, sliceLen int, nested bool) *DataEmpty {
	res := &DataEmpty{
		String:            "test",
		Int16:             42,
		Float32:           3.14,
		Bool:              true,
		Time:              time.Now(),
		TimePtr:           nil,
		EmbeddedPrimitive: "embedded",
	}
	if nested {
		nested := newDataEmpty(0, 0, false)
		res.Nested = &nested
		res.IfacePtr = nested
	}
	for i := 0; i < sliceLen; i++ {
		res.Slice = append(res.Slice, newDataEmpty(0, 0, false))
	}
	if depth > 0 {
		res.DataEmpty = newDataEmpty(depth-1, sliceLen, false)
	}
	return res
}

type DataEmpty struct {
	String            string        `json:"string,omitempty"`
	Int16             int16         `json:"int16,omitempty"`
	Float32           float32       `json:"float32,omitempty"`
	Bool              bool          `json:"bool,omitempty"`
	Time              time.Time     `json:"time,omitempty"`
	TimePtr           *time.Time    `json:"timePtr,omitempty"`
	Slice             []*DataEmpty  `json:"slice,omitempty"`
	Nested            **DataEmpty   `json:"nested,omitempty"`
	IfacePtr          DataInterface `json:"ifacePtr,omitempty"`
	EmbeddedPrimitive `json:"embeddedPrimitive,omitempty"`
	*DataEmpty        `json:"embedded,omitempty"`
}

func (d *DataEmpty) InterfaceMethod() {}

func TestEncode_Struct(t *testing.T) {
	// Omit empty
	v := newDataEmpty(2, 1, true)
	EqualMarshaling(t, v)

	// No omit empty
	v2 := newData(3, 10, true)
	EqualMarshaling(t, v2)
}

func TestEncode_Struct_Empty(t *testing.T) {
	// Omit empty
	str := DataEmpty{
		Time: time.Now(),
	}
	EqualMarshaling(t, str)

	// No omit empty
	str1 := Data{}
	EqualMarshaling(t, str1)
}

// Benchmarks
func BenchmarkStruct_Empty_Blaze(b *testing.B) {
	v := newDataEmpty(5, 10, true)
	bytes := []byte{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bytes, _ = Marshal(v)
		}
	})
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkStruct_Empty_Std(b *testing.B) {
	v := newDataEmpty(5, 10, true)
	bytes := []byte{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bytes, _ = json.Marshal(v)
		}
	})
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkStruct_Empty_GoJson(b *testing.B) {
	v := newDataEmpty(5, 10, true)
	bytes := []byte{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bytes, _ = gojson.Marshal(v)
		}
	})
	b.SetBytes(int64(len(bytes)))
}
