package decoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
)

func TestDecode_Bool(t *testing.T) {
	data := []byte("true")
	EqualUnmarshaling[bool](t, data)

	data = []byte("false")
	EqualUnmarshaling[bool](t, data)

	EqualUnmarshaling[*bool](t, data)
}

func BenchmarkBool_Blaze(b *testing.B) {
	var boolean bool
	for i := 0; i < b.N; i++ {
		Unmarshal([]byte("true"), &boolean)
	}
	b.SetBytes(int64(len("true")))
}

func BenchmarkBool_Std(b *testing.B) {
	var boolean bool
	for i := 0; i < b.N; i++ {
		json.Unmarshal([]byte("true"), &boolean)
	}
	b.SetBytes(int64(len("true")))
}

func BenchmarkBool_GoJson(b *testing.B) {
	var boolean bool
	for i := 0; i < b.N; i++ {
		gojson.Unmarshal([]byte("true"), &boolean)
	}
	b.SetBytes(int64(len("true")))
}

func BenchmarkBool_JsonIter(b *testing.B) {
	var boolean bool
	for i := 0; i < b.N; i++ {
		jsoniter.Unmarshal([]byte("true"), &boolean)
	}
	b.SetBytes(int64(len("true")))
}
