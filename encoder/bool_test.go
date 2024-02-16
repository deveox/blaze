package encoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
)

func TestEncode_Bool(t *testing.T) {
	var tr bool = true
	EqualMarshaling(t, tr)

	var f bool = false
	EqualMarshaling(t, f)
}

func TestEncode_Ptr_Bool(t *testing.T) {
	var tr = true
	EqualMarshaling(t, &tr)

	var f = false
	EqualMarshaling(t, &f)
}

func TestEncode_BoolLongPtr(t *testing.T) {
	tr := true
	v1 := &tr
	v2 := &v1
	v3 := &v2
	v4 := &v3
	v5 := &v4
	EqualMarshaling(t, v5)

	f := false
	v1 = &f
	v2 = &v1
	v3 = &v2
	v4 = &v3
	v5 = &v4
	EqualMarshaling(t, v5)
}

// Benchmarks

func BenchmarkBool_Blaze(b *testing.B) {
	v := false
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = DEncoder.Marshal(v)
	}
	b.SetBytes(int64(len(bytes)))
}
func BenchmarkBool_Std(b *testing.B) {
	v := false
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = json.Marshal(v)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkBool_GoJson(b *testing.B) {
	v := false
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = gojson.Marshal(v)
	}
	b.SetBytes(int64(len(bytes)))
}
