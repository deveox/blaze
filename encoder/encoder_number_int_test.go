package encoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
)

func TestEncode_Int(t *testing.T) {
	var integer int = 0
	EqualMarshaling(t, integer)
}

func TestEncode_Ptr_Int(t *testing.T) {
	var integer int = 0
	EqualMarshaling(t, &integer)
}

func TestEncode_Int8(t *testing.T) {
	var integer int8 = 127
	EqualMarshaling(t, integer)
}

func TestEncode_Ptr_Int8(t *testing.T) {
	var integer int8 = 127
	EqualMarshaling(t, &integer)
}

func TestEncode_Int16(t *testing.T) {
	var integer int16 = 32767
	EqualMarshaling(t, integer)
}

func TestEncode_Ptr_Int16(t *testing.T) {
	var integer int16 = 32767
	EqualMarshaling(t, &integer)
}

func TestEncode_Int32(t *testing.T) {
	var integer int32 = 2147483647
	EqualMarshaling(t, integer)
}

func TestEncode_Int64(t *testing.T) {
	var integer int64 = 9223372036854775807
	EqualMarshaling(t, integer)
}

func TestEncode_Ptr_Int64(t *testing.T) {
	var integer int64 = 9223372036854775807
	EqualMarshaling(t, &integer)
}

// Benchmarks

func BenchmarkInt_Blaze(b *testing.B) {
	in := 102547890
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = Marshal(in)
	}
	b.SetBytes(int64(len(bytes)))
}
func BenchmarkInt_Std(b *testing.B) {
	in := 102547890
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = json.Marshal(in)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkInt_GoJson(b *testing.B) {
	in := 102547890
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = gojson.Marshal(in)
	}
	b.SetBytes(int64(len(bytes)))
}
