package decoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
)

func encodedSlice(size int) []byte {
	res := make([]byte, 0, size*13+2)
	res = append(res, '[')
	for i := 0; i < size; i++ {
		res = append(res, []byte(`"hello world"`)...)
		res = append(res, ',')
	}
	res[len(res)-1] = ']'
	return res
}

func encodedSliceInt(size int) []byte {
	res := make([]byte, 0, size)

	res = append(res, '[')
	for i := 0; i < size; i++ {
		b := make([]byte, 0, size)
		b = append(b, '[')
		for j := 0; j < size; j++ {
			c := make([]byte, 0, size)
			c = append(c, '[')
			for k := 0; k < size; k++ {
				d := make([]byte, 0, size)
				d = append(d, '[')
				for l := 0; l < size; l++ {
					d = append(d, []byte(`10000`)...)
					d = append(d, ',')
				}
				d[len(d)-1] = ']'
				c = append(c, d...)
				c = append(c, ',')
			}
			c[len(c)-1] = ']'
			b = append(b, c...)
			b = append(b, ',')
		}
		b[len(b)-1] = ']'
		res = append(res, b...)
		res = append(res, ',')
	}
	res[len(res)-1] = ']'
	return res
}

func TestDecode_Slice(t *testing.T) {
	slice := encodedSlice(10)

	EqualUnmarshaling[[]string](t, slice)

	EqualUnmarshaling[**[]string](t, slice)
}

var benchSlice = encodedSlice(100)

func BenchmarkSlice_String_Blaze(b *testing.B) {
	var slice []string
	data := encodedSlice(100)
	for i := 0; i < b.N; i++ {
		Unmarshal(data, &slice)
	}
	b.SetBytes(int64(len(data)))
}

func BenchmarkSlice_String_Std(b *testing.B) {
	var slice []string
	for i := 0; i < b.N; i++ {
		json.Unmarshal(benchSlice, &slice)
	}
	b.SetBytes(int64(len(benchSlice)))
}

func BenchmarkSlice_String_GoJson(b *testing.B) {
	var slice []string
	for i := 0; i < b.N; i++ {
		gojson.Unmarshal(benchSlice, &slice)
	}
	b.SetBytes(int64(len(benchSlice)))
}

func BenchmarkSlice_String_JsonIter(b *testing.B) {
	var slice []string
	for i := 0; i < b.N; i++ {
		jsoniter.Unmarshal(benchSlice, &slice)
	}
	b.SetBytes(int64(len(benchSlice)))
}

var benchSliceInt = encodedSliceInt(3)

func BenchmarkMatrix_Int_Blaze(b *testing.B) {
	var slice [][][][]int
	data := encodedSliceInt(3)
	for i := 0; i < b.N; i++ {
		err := Unmarshal(data, &slice)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(data)))
}

func BenchmarkMatrix_Int_Std(b *testing.B) {
	var slice [][][][]int
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(benchSliceInt, &slice)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(benchSliceInt)))
}

func BenchmarkMatrix_Int_GoJson(b *testing.B) {
	var slice [][][][]int
	for i := 0; i < b.N; i++ {
		gojson.Unmarshal(benchSliceInt, &slice)
	}
	b.SetBytes(int64(len(benchSliceInt)))
}

func BenchmarkMatrix_Int_JsonIter(b *testing.B) {
	var slice [][][][]int
	for i := 0; i < b.N; i++ {
		jsoniter.Unmarshal(benchSliceInt, &slice)
	}
	b.SetBytes(int64(len(benchSliceInt)))
}
