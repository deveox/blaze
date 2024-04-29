package decoder

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type StructWithRaw struct {
	Name string
	Raw  json.RawMessage
	Age  int
}

func TestDecode_RawMessage(t *testing.T) {
	raw := json.RawMessage(`{"key":"value"}`)
	data := []byte(`{"name":"test","raw":`)
	data = append(data, raw...)
	data = append(data, []byte(`,"age":10}`)...)
	var s StructWithRaw
	err := DDecoder.Unmarshal(data, &s)
	require.NoError(t, err)
	require.Equal(t, "test", s.Name)
	require.Equal(t, raw, s.Raw)
	require.Equal(t, 10, s.Age)
}

func TestDecode_Slice_Struct(t *testing.T) {
	data := []byte(`[{"name":"test1","age":10},{"name":"test2","age":20}]`)
	type Test struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var s []Test
	err := DDecoder.Unmarshal(data, &s)
	require.NoError(t, err)
	require.Len(t, s, 2)
	require.Equal(t, "test1", s[0].Name)
	require.Equal(t, 10, s[0].Age)
	require.Equal(t, "test2", s[1].Name)
	require.Equal(t, 20, s[1].Age)
}
