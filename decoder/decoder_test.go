package decoder

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

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
	if err := Unmarshal(data, &v2); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, v, v2)
}

func EqualTo[T any](t *testing.T, data []byte, expected T) {
	var v T
	if err := Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, expected, v)
}

func UnmarshalNull[T any](t *testing.T, byDefault T) {
	data := []byte("null")
	v := byDefault
	if err := Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	var expected T
	require.Equal(t, expected, v)
}
