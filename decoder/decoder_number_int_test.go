package decoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
)

func TestDecode_Int(t *testing.T) {
	data := []byte("100")
	EqualUnmarshaling[int](t, data)

	data = []byte("-100")
	EqualUnmarshaling[int](t, data)

	data = []byte("0")
	EqualUnmarshaling[int](t, data)

}

func TestDecode_Int8(t *testing.T) {
	data := []byte("100")
	EqualUnmarshaling[int8](t, data)

	data = []byte("-100")
	EqualUnmarshaling[int8](t, data)

	data = []byte("0")
	EqualUnmarshaling[int8](t, data)
}

func TestDecode_Int16(t *testing.T) {
	data := []byte("100")
	EqualUnmarshaling[int16](t, data)

	data = []byte("-100")
	EqualUnmarshaling[int16](t, data)

	data = []byte("0")
	EqualUnmarshaling[int16](t, data)

}

func TestDecode_Int32(t *testing.T) {
	data := []byte("100")
	EqualUnmarshaling[int32](t, data)

	data = []byte("-100")
	EqualUnmarshaling[int32](t, data)

	data = []byte("0")
	EqualUnmarshaling[int32](t, data)

}

func TestDecode_Int64(t *testing.T) {
	data := []byte("100")
	EqualUnmarshaling[int64](t, data)

	data = []byte("-100")
	EqualUnmarshaling[int64](t, data)

	data = []byte("0")
	EqualUnmarshaling[int64](t, data)

}

var benchInt64 = []byte("1000000000000")

func BenchmarkInt64_Blaze(b *testing.B) {
	var integer int64
	for i := 0; i < b.N; i++ {
		DDecoder.Unmarshal(benchInt64, &integer)
	}
	b.SetBytes(int64(len(benchInt64)))
}

func BenchmarkInt64_Std(b *testing.B) {
	var integer int64
	for i := 0; i < b.N; i++ {
		json.Unmarshal(benchInt64, &integer)
	}
	b.SetBytes(int64(len(benchInt64)))
}

func BenchmarkInt64_GoJson(b *testing.B) {
	var integer int64
	for i := 0; i < b.N; i++ {
		gojson.Unmarshal(benchInt64, &integer)
	}
	b.SetBytes(int64(len(benchInt64)))
}

func BenchmarkInt64_JsonIter(b *testing.B) {
	var integer int64
	for i := 0; i < b.N; i++ {
		jsoniter.Unmarshal(benchInt64, &integer)
	}
	b.SetBytes(int64(len(benchInt64)))
}
