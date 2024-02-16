package decoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
)

func TestDecode_Map(t *testing.T) {
	data := []byte(`{"hello":"world"}`)
	EqualUnmarshaling[map[string]string](t, data)
	EqualUnmarshaling[**map[string]string](t, data)
	data = []byte(`{"2":2,"3":3}`)
	EqualUnmarshaling[map[int]int](t, data)
	EqualUnmarshaling[map[int8]int](t, data)
	EqualUnmarshaling[map[int16]int](t, data)
	EqualUnmarshaling[map[int32]int](t, data)
	EqualUnmarshaling[map[int64]int](t, data)
	EqualUnmarshaling[map[uint]int](t, data)
	EqualUnmarshaling[map[uint8]int](t, data)
	EqualUnmarshaling[map[uint16]int](t, data)
	EqualUnmarshaling[map[uint32]int](t, data)
	EqualUnmarshaling[map[uint64]int](t, data)

	data = []byte(`{"2.2":2,"3.3":3}`)
	EqualTo[map[float32]int](t, data, map[float32]int{2.2: 2, 3.3: 3})
	EqualTo[map[float64]int](t, data, map[float64]int{2.2: 2, 3.3: 3})

}

var benchMap = []byte(`{"hello":"world"}`)

func BenchmarkMap_Blaze(b *testing.B) {
	var m map[string]string
	for i := 0; i < b.N; i++ {
		DDecoder.Unmarshal(benchMap, &m)
	}
	b.SetBytes(int64(len(benchMap)))
}

func BenchmarkMap_Std(b *testing.B) {
	var m map[string]string
	for i := 0; i < b.N; i++ {
		json.Unmarshal(benchMap, &m)
	}
	b.SetBytes(int64(len(benchMap)))
}

func BenchmarkMap_GoJson(b *testing.B) {
	var m map[string]string
	for i := 0; i < b.N; i++ {
		gojson.Unmarshal(benchMap, &m)
	}
	b.SetBytes(int64(len(benchMap)))
}

func BenchmarkMap_JsonIter(b *testing.B) {
	var m map[string]string
	for i := 0; i < b.N; i++ {
		jsoniter.Unmarshal(benchMap, &m)
	}
	b.SetBytes(int64(len(benchMap)))
}
