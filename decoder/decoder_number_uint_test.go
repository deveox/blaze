package decoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
)

func TestDecode_Uint(t *testing.T) {
	data := []byte("100")
	EqualUnmarshaling[uint](t, data)

	data = []byte("0")
	EqualUnmarshaling[uint](t, data)

}

func TestDecode_Uint8(t *testing.T) {
	data := []byte("100")
	EqualUnmarshaling[uint8](t, data)

	data = []byte("0")
	EqualUnmarshaling[uint8](t, data)
}

func TestDecode_Uint16(t *testing.T) {
	data := []byte("100")
	EqualUnmarshaling[uint16](t, data)

	data = []byte("0")
	EqualUnmarshaling[uint16](t, data)

}

func TestDecode_Uint32(t *testing.T) {
	data := []byte("100")
	EqualUnmarshaling[uint32](t, data)

	data = []byte("0")
	EqualUnmarshaling[uint32](t, data)

}

func TestDecode_Uint64(t *testing.T) {
	data := []byte("100")
	EqualUnmarshaling[uint64](t, data)

	data = []byte("0")
	EqualUnmarshaling[uint64](t, data)

}

var benchUint64 = []byte("1000000000000")

func BenchmarkUint64_Blaze(b *testing.B) {
	var v uint64
	for i := 0; i < b.N; i++ {
		DDecoder.Unmarshal(benchUint64, &v)
	}
	b.SetBytes(int64(len(benchUint64)))
}

func BenchmarkUint64_Std(b *testing.B) {
	var v uint64
	for i := 0; i < b.N; i++ {
		json.Unmarshal(benchUint64, &v)
	}
	b.SetBytes(int64(len(benchUint64)))
}

func BenchmarkUint64_GoJson(b *testing.B) {
	var v uint64
	for i := 0; i < b.N; i++ {
		gojson.Unmarshal(benchUint64, &v)
	}
	b.SetBytes(int64(len(benchUint64)))
}

func BenchmarkUint64_JsonIter(b *testing.B) {
	var v uint64
	for i := 0; i < b.N; i++ {
		jsoniter.Unmarshal(benchUint64, &v)
	}
	b.SetBytes(int64(len(benchUint64)))
}
