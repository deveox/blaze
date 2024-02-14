package encoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
)

func TestEncode_Uint(t *testing.T) {
	var v uint = 0
	EqualMarshaling(t, v)
}

func TestEncode_Ptr_Uint(t *testing.T) {
	var v uint = 0
	EqualMarshaling(t, &v)
}

func TestEncode_Uint8(t *testing.T) {
	var v uint8 = 255
	EqualMarshaling(t, v)
}

func TestEncode_Ptr_Uint8(t *testing.T) {
	var v uint8 = 255
	EqualMarshaling(t, &v)
}

func TestEncode_Uint16(t *testing.T) {
	var v uint16 = 65535
	EqualMarshaling(t, v)
}

func TestEncode_Ptr_Uint16(t *testing.T) {
	var v uint16 = 65535
	EqualMarshaling(t, &v)
}

func TestEncode_Uint32(t *testing.T) {
	var v uint32 = 4294967295
	EqualMarshaling(t, v)
}

func TestEncode_Ptr_Uint32(t *testing.T) {
	var v uint32 = 4294967295
	EqualMarshaling(t, &v)
}

func TestEncode_Uint64(t *testing.T) {
	var v uint64 = 18446744073709551615
	EqualMarshaling(t, v)
}

func TestEncode_Ptr_Uint64(t *testing.T) {
	var v uint64 = 18446744073709551615
	EqualMarshaling(t, &v)
}

// Benchmarks

func BenchmarkUint_Blaze(b *testing.B) {
	uin := uint(102547890)
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = Marshal(uin)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkUint_Std(b *testing.B) {
	uin := uint(102547890)
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = json.Marshal(uin)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkUint_GoJson(b *testing.B) {
	uin := uint(102547890)
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = gojson.Marshal(uin)
	}
	b.SetBytes(int64(len(bytes)))
}
