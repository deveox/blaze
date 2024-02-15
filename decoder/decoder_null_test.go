package decoder

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode_Null(t *testing.T) {
	UnmarshalNull[interface{}](t, "test")

	UnmarshalNull[string](t, "test")
	str := "test"
	UnmarshalNull[*string](t, &str)
	UnmarshalNull[int](t, 22)
	UnmarshalNull[float32](t, 22)
	UnmarshalNull[float64](t, 22)
	UnmarshalNull[bool](t, true)
	UnmarshalNull[[]int](t, []int{1, 2, 3})
	UnmarshalNull[map[string]int](t, map[string]int{"a": 1})

	type A struct {
		A int
	}
	UnmarshalNull[A](t, A{A: 1})

	v := &A{A: 1}
	if err := json.Unmarshal([]byte("null"), v); err != nil {
		t.Fatal(err)
	}

	v2 := &A{A: 1}
	if err := Unmarshal([]byte("null"), v2); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, &A{}, v2)
}

func TestDecode_Null_Slice(t *testing.T) {
	data := []byte(`[null, null, null]`)
	var v []int
	err := Unmarshal(data, &v)
	require.NoError(t, err)
	require.Equal(t, []int{0, 0, 0}, v)
}
