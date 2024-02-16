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
