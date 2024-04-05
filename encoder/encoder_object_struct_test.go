package encoder

import (
	"encoding/json"
	"testing"
	"time"

	gojson "github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
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
		EmbeddedStruct: EmbeddedStruct{
			Id:      42,
			Created: time.Now(),
			Updated: time.Now(),
		},
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

type EmbeddedStruct struct {
	Id      uint64    `json:"id,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}
type DataEmpty struct {
	EmbeddedStruct
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
	*DataEmpty        `json:",omitempty"`
}

func (d *DataEmpty) InterfaceMethod() {}

type EmbeddedIgnored struct {
	Field string `json:"-"`
}
type WithEmbeddedIgnored struct {
	EmbeddedIgnored
	IgnoredField IgnoredField `json:"ignoredField,omitempty"`
	Hello        int          `json:"hello"`
}

type IgnoredField string

func (i IgnoredField) MarshalBlaze(e *Encoder) error {
	return nil
}
func (i IgnoredField) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

func TestEncode_Struct_Ignored(t *testing.T) {
	s := WithEmbeddedIgnored{
		EmbeddedIgnored: EmbeddedIgnored{Field: "123"},
		IgnoredField:    IgnoredField("ignored"),
		Hello:           123,
	}
	EqualString(t, s, `{"hello":123}`)
}
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
		EmbeddedStruct: EmbeddedStruct{
			Created: time.Now(),
			Updated: time.Now(),
		},
	}
	EqualMarshaling(t, str)

	// No omit empty
	str1 := Data{}
	EqualMarshaling(t, str1)
}

type PartialEmbedded string

type PartialStructEmbedded struct {
	ShortEmbedded string `blaze:"short"`
}

type PartialStruct struct {
	Short           string `blaze:"short"`
	NotShort        string
	Ignored         string
	PartialEmbedded `blaze:"short"`
	Nested          *PartialStruct `blaze:"short"`
	*PartialStructEmbedded
}

func newPartialStruct() *PartialStruct {
	return &PartialStruct{
		Short:    "short",
		NotShort: "not short",
		Ignored:  "ignored",
		PartialStructEmbedded: &PartialStructEmbedded{
			ShortEmbedded: "short embedded",
		},
		PartialEmbedded: "embedded",
		Nested: &PartialStruct{
			Short:    "short nested",
			NotShort: "not short nested",
			Ignored:  "ignored nested",
		},
	}
}

func TestEncode_Partial_Short(t *testing.T) {
	v := newPartialStruct()
	wanted := []byte(`{"short":"short","partialEmbedded":"embedded","nested":{"short":"short nested"},"shortEmbedded":"short embedded"}`)
	bytes, err := DEncoder.MarshalPartial(v, nil, true)
	require.NoError(t, err)
	require.Equal(t, string(wanted), string(bytes))
}

func TestEncode_Partial_Fields(t *testing.T) {
	v := newPartialStruct()
	wanted := []byte(`{"notShort":"not short","partialEmbedded":"embedded","nested":{"notShort":"not short nested"},"shortEmbedded":"short embedded"}`)
	bytes, err := DEncoder.MarshalPartial(v, []string{"notShort", "partialEmbedded", "shortEmbedded", "nested.notShort"}, false)
	require.NoError(t, err)
	require.Equal(t, string(wanted), string(bytes))
}

func TestEncode_Partial(t *testing.T) {
	v := newPartialStruct()
	wanted := []byte(`{"short":"short","notShort":"not short","partialEmbedded":"embedded","nested":{"short":"short nested","notShort":"not short nested"},"shortEmbedded":"short embedded"}`)
	bytes, err := DEncoder.MarshalPartial(v, []string{"notShort", "nested.notShort"}, true)
	require.NoError(t, err)
	require.Equal(t, string(wanted), string(bytes))
}

// Benchmarks
func BenchmarkStruct_Empty_Blaze(b *testing.B) {
	v := newDataEmpty(5, 10, true)
	bytes := []byte{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bytes, _ = DEncoder.Marshal(v)
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
