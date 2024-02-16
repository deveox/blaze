package encoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
)

func newSlice() []interface{} {
	return []interface{}{
		42,
		"Hello, World!",
		3.14,
		true,
		Data{},
	}
}

func TestEncode_Slice(t *testing.T) {
	sl := newSlice()
	EqualMarshaling(t, sl)

	sl = []interface{}{}
	EqualMarshaling(t, sl)
}

func benchMatrix() [][][][]string {
	s := make([][][][]string, 0, 10)
	for i := 0; i < 10; i++ {
		v := make([][][]string, 0, 10)
		for j := 0; j < 10; j++ {
			vv := make([][]string, 0, 10)
			for k := 0; k < 10; k++ {
				vvv := make([]string, 0, 10)
				for l := 0; l < 10; l++ {
					vvv = append(vvv, "hello")
				}
				vv = append(vv, vvv)
			}
			v = append(v, vv)
		}
		s = append(s, v)
	}
	return s
}

func benchSlice() []int {
	s := make([]int, 0, 10000)
	for i := 0; i < 10000; i++ {
		s = append(s, i)
	}
	return s
}

func BenchmarkSlice_Matrix_Blaze(b *testing.B) {
	s := benchMatrix()

	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = DEncoder.Marshal(s)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkSlice_Matrix_Std(b *testing.B) {
	s := benchMatrix()
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = json.Marshal(s)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkSlice_Matrix_GoJson(b *testing.B) {
	s := benchMatrix()
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = gojson.Marshal(s)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkSlice_Simple_Blaze(b *testing.B) {
	s := benchSlice()

	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = DEncoder.Marshal(s)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkSlice_Simple_Std(b *testing.B) {
	s := benchSlice()
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = json.Marshal(s)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkSlice_Simple_GoJson(b *testing.B) {
	s := benchSlice()
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = gojson.Marshal(s)
	}
	b.SetBytes(int64(len(bytes)))
}
