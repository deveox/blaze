package encoder

import (
	"encoding/json"
	"strings"
	"testing"

	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
)

func TestEncode_String(t *testing.T) {
	str := newStr()
	EqualMarshaling(t, str)
}

func TestEncode_String_Empty(t *testing.T) {
	str := ""
	EqualMarshaling(t, str)
}

func TestEncode_String_Escaped(t *testing.T) {
	str := newStrEsc()

	EqualMarshaling(t, str)
}

func TestEncode_Ptr_String(t *testing.T) {
	EqualMarshaling(t, &benchStr)
}

func newStrEsc() string {
	s := strings.Builder{}
	for i := 0; i < 100; i++ {
		s.WriteString("hello \"world\r\n")
	}
	return s.String()
}

func newStr() string {
	s := strings.Builder{}
	for i := 0; i < 100; i++ {
		s.WriteString("hello world")
	}
	return s.String()
}

var benchStr = newStr()

var benchStrEsc = newStrEsc()

// Benchmarks

func BenchmarkString_Blaze_Simple(b *testing.B) {
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = Marshal(benchStr)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkString_Std_Simple(b *testing.B) {
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = json.Marshal(benchStr)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkString_GoJson_Simple(b *testing.B) {
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = gojson.Marshal(benchStr)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkString_JsonIter_Simple(b *testing.B) {
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = jsoniter.Marshal(benchStr)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkString_Blaze_Escaped(b *testing.B) {
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = Marshal(benchStrEsc)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkString_Std_Escaped(b *testing.B) {
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = json.Marshal(benchStrEsc)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkString_GoJson_Escaped(b *testing.B) {
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = gojson.Marshal(benchStrEsc)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkString_JsonIter_Escaped(b *testing.B) {
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = jsoniter.Marshal(benchStrEsc)
	}
	b.SetBytes(int64(len(bytes)))
}
