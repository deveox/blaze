package encoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
)

// Float32

func TestEncode_Float32(t *testing.T) {
	// Zero Float
	var v float32 = 0
	EqualMarshaling(t, v)

	// Small Float
	v = 3.14
	EqualMarshaling(t, v)

	// Large Float
	v = 3.14159265358979323846
	EqualMarshaling(t, v)
}

func TestEncode_Ptr_Float32(t *testing.T) {
	// Zero Float
	var v float32 = 0
	EqualMarshaling(t, &v)

	// Small Float
	v = 3.14
	EqualMarshaling(t, &v)

	// Large Float
	v = 3.14159265358979323846
	EqualMarshaling(t, &v)
}

// Float64

func TestEncode_Float64(t *testing.T) {
	// Zero Float
	var v float64 = 0
	EqualMarshaling(t, v)

	// Small Float
	v = 3.14
	EqualMarshaling(t, v)

	// Large Float
	v = 3.14159265358979323846
	EqualMarshaling(t, v)
}

func TestEncode_Ptr_Float64(t *testing.T) {
	// Zero Float
	var v float64 = 0
	EqualMarshaling(t, &v)

	// Small Float
	v = 3.14
	EqualMarshaling(t, &v)

	// Large Float
	v = 3.14159265358979323846
	EqualMarshaling(t, &v)
}

// Benchmarks

func BenchmarkFloat64_Blaze(b *testing.B) {
	f := 100000000003.14159265358979323846
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = Marshal(f)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkFloat64_Std(b *testing.B) {
	f := 100000000003.14159265358979323846
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = json.Marshal(f)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkFloat64_GoJson(b *testing.B) {
	f := 100000000003.14159265358979323846
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = gojson.Marshal(f)
	}
	b.SetBytes(int64(len(bytes)))
}
