package encoder

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// Default encoder for tests
var DEncoder = &Config{}

func EqualString(t *testing.T, v interface{}, expected string) {
	bytes, err := DEncoder.Marshal(v)
	require.NoError(t, err)
	require.Equal(t, expected, string(bytes))

}

func EqualMarshaling(t *testing.T, v interface{}) {
	stdBytes, err := json.Marshal(v)
	require.NoError(t, err)
	bytes, err := DEncoder.Marshal(v)
	require.NoError(t, err)
	stdBytes = AddIndent(stdBytes)
	bytes = AddIndent(bytes)
	require.Equal(t, string(stdBytes), string(bytes))
}

func EqualMap(t *testing.T, v interface{}) {
	stdBytes, err := json.Marshal(v)
	require.NoError(t, err)
	bytes, err := DEncoder.Marshal(v)
	require.NoError(t, err)

	var v2 interface{}
	err = json.Unmarshal(stdBytes, &v2)
	require.NoError(t, err)
	var v3 interface{}
	err = json.Unmarshal(bytes, &v3)
	require.NoError(t, err)
	require.Equal(t, v2, v3)
}
