package decoder

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

// Default decoder for tests
var DDecoder = &Config{}

func EqualUnmarshaling[T any](t *testing.T, data []byte) {
	var v T
	EqualUnmarshalingTo(t, data, v)
}

func EqualUnmarshalingTo[T any](t *testing.T, data []byte, byDefault T) {
	v := byDefault
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}

	v2 := byDefault
	if err := DDecoder.Unmarshal(data, &v2); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, v, v2)
}

func EqualTo[T any](t *testing.T, data []byte, expected T) {
	var v T
	if err := DDecoder.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, expected, v)
}

func UnmarshalNull[T any](t *testing.T, byDefault T) {
	data := []byte("null")
	v := byDefault
	if err := DDecoder.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	rv := reflect.ValueOf(v)
	if rv.IsValid() {
		require.True(t, rv.IsZero())
	}
}

type WithChangesNested struct {
	Name   string
	Age    int
	Nested *WithChanges
}
type WithChanges struct {
	Name   string
	Age    int `blaze:"read"`
	Nested *WithChangesNested
}

func TestUnmarshal_WithChanges(t *testing.T) {
	data := []byte(`{"name":"test","age":10,"nested":{"name":"test","nested":{"name":"test","age":10}}}`)
	var v WithChanges
	changes, err := DDecoder.UnmarshalWithChanges(data, &v)
	require.NoError(t, err)
	require.Equal(t, []string{"name", "nested", "nested.name", "nested.nested", "nested.nested.name"}, changes)
	data = []byte(`{"name":"test","age":10,"nested":{"name":"test","age":10, "nested":{"age":10}}}`)
	changes, err = DDecoder.UnmarshalWithChanges(data, &v)
	require.NoError(t, err)
	require.Equal(t, []string{"name", "nested", "nested.name", "nested.age"}, changes)

}
