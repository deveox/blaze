package decoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
)

func TestDecode_Float32(t *testing.T) {
	data := []byte("100.123")
	EqualUnmarshaling[float32](t, data)

	data = []byte("-100.123")
	EqualUnmarshaling[float32](t, data)

	data = []byte("0")
	EqualUnmarshaling[float32](t, data)

	data = []byte("0.003")
	EqualUnmarshaling[float32](t, data)

	data = []byte("0e6")
	EqualUnmarshaling[float32](t, data)

	data = []byte("0e-6")
	EqualUnmarshaling[float32](t, data)

	data = []byte("0e+6")
	EqualUnmarshaling[float32](t, data)

	data = []byte("0.003e6")
	EqualUnmarshaling[float32](t, data)

	data = []byte("0.003e-6")
	EqualUnmarshaling[float32](t, data)

	data = []byte("0.003e+6")
	EqualUnmarshaling[float32](t, data)
}

func TestDecode_Float64(t *testing.T) {
	data := []byte("100.123")
	EqualUnmarshaling[float64](t, data)

	data = []byte("-100.123")
	EqualUnmarshaling[float64](t, data)

	data = []byte("0")
	EqualUnmarshaling[float64](t, data)

	data = []byte("0.003")
	EqualUnmarshaling[float64](t, data)

	data = []byte("0e6")
	EqualUnmarshaling[float64](t, data)

	data = []byte("0e-6")
	EqualUnmarshaling[float64](t, data)

	data = []byte("0e+6")
	EqualUnmarshaling[float64](t, data)

	data = []byte("0.003e6")
	EqualUnmarshaling[float64](t, data)

	data = []byte("0.003e-6")
	EqualUnmarshaling[float64](t, data)

	data = []byte("0.003e+6")
	EqualUnmarshaling[float64](t, data)
}

var benchFloat64 = []byte("100.123e10")

func BenchmarkFloat64_Blaze(b *testing.B) {
	var v float64
	for i := 0; i < b.N; i++ {
		DDecoder.Unmarshal(benchFloat64, &v)
	}
	b.SetBytes(int64(len(benchFloat64)))
}

func BenchmarkFloat64_Std(b *testing.B) {
	var v float64
	for i := 0; i < b.N; i++ {
		json.Unmarshal(benchFloat64, &v)
	}
	b.SetBytes(int64(len(benchFloat64)))
}

func BenchmarkFloat64_GoJson(b *testing.B) {
	var v float64
	for i := 0; i < b.N; i++ {
		gojson.Unmarshal(benchFloat64, &v)
	}
	b.SetBytes(int64(len(benchFloat64)))
}

func BenchmarkFloat64_JsonIter(b *testing.B) {
	var v float64
	for i := 0; i < b.N; i++ {
		jsoniter.Unmarshal(benchFloat64, &v)
	}
	b.SetBytes(int64(len(benchFloat64)))
}
