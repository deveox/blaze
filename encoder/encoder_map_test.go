package encoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
)

func TestEncode_Map_String(t *testing.T) {
	// Zero Map
	var m map[string]interface{}
	EqualMarshaling(t, m)

	// Empty Map
	m = map[string]interface{}{}
	EqualMarshaling(t, m)

	// Map with values
	m = map[string]interface{}{
		"42":        false,
		"stringKey": 20,
		"true":      []int{1, 2, 3},
		"struct": map[string]interface{}{
			"key1": "value1",
			"key2": 42,
			"key3": true,
		},
	}
	EqualMap(t, m)

}

func getMap() map[string]interface{} {
	v := []int{1, 2, 3}
	return map[string]interface{}{
		"42":        false,
		"stringKey": 20,
		"true":      &v,
		"Data{}": map[string]interface{}{
			"key1": "value1",
			"key2": 42,
			"key3": true,
		},
		"&v": "value",
	}
}

func TestEncode_Map_Any(t *testing.T) {
	// Zero Map
	var m map[string]interface{}
	EqualMarshaling(t, m)

	// Empty Map
	m = map[string]interface{}{}
	EqualMarshaling(t, m)

	// Map with values
	m = getMap()
	EqualMap(t, m)
}

func TestEncode_Ptr_Map(t *testing.T) {
	// Zero Map
	var m map[string]interface{}
	EqualMarshaling(t, &m)

	// Empty Map
	m = map[string]interface{}{}
	EqualMarshaling(t, &m)

	// Map with values
	m = getMap()
	EqualMap(t, &m)
}

// Benchmarks

func getBenchMap(n int) map[string]string {
	m := make(map[string]string)
	for i := 0; i < n; i++ {
		m["LongKeyLongKeyLongKeyLongKeyLongKeyLongKeyLongKeyLongKeyLongKeyLongKeyLongKeyLongKeyLongKeyLongKey"] = "LongstringLongstringLongstringLongstringLongstringLongstringLongstringLongstringLongstring"
	}
	return m
}

func BenchmarkMap_Simple_Blaze(b *testing.B) {
	s := getBenchMap(100)
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = Marshal(s)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkMap_Simple_Std(b *testing.B) {
	s := getBenchMap(100)
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = json.Marshal(s)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkMap_Simple_GoJson(b *testing.B) {
	s := getBenchMap(100)
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = gojson.Marshal(s)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkMap_Simple_Big_Blaze(b *testing.B) {
	s := getBenchMap(10_000)
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = Marshal(s)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkMap_Simple_Big_Std(b *testing.B) {
	s := getBenchMap(10_000)
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = json.Marshal(s)
	}
	b.SetBytes(int64(len(bytes)))
}

func BenchmarkMap_Simple_Big_GoJson(b *testing.B) {
	s := getBenchMap(10_000)
	bytes := []byte{}
	for i := 0; i < b.N; i++ {
		bytes, _ = gojson.Marshal(s)
	}
	b.SetBytes(int64(len(bytes)))
}
