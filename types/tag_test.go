package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagToOperation(t *testing.T) {
	s := "read.update.create"
	res := tagPartToOperation(s)
	require.Equal(t, OPERATION_ALL, res)

	s = "read.update"
	res = tagPartToOperation(s)
	require.Equal(t, OPERATION_READ_UPDATE, res)

	s = "read.create"
	res = tagPartToOperation(s)
	require.Equal(t, OPERATION_READ_CREATE, res)

	s = "read"
	res = tagPartToOperation(s)
	require.Equal(t, OPERATION_READ, res)

	s = "update.create"
	res = tagPartToOperation(s)
	require.Equal(t, OPERATION_WRITE, res)

	s = "update"
	res = tagPartToOperation(s)
	require.Equal(t, OPERATION_UPDATE, res)

	s = "create"
	res = tagPartToOperation(s)
	require.Equal(t, OPERATION_CREATE, res)

	s = "-"
	res = tagPartToOperation(s)
	require.Equal(t, OPERATION_IGNORE, res)

	s = "all"
	res = tagPartToOperation(s)
	require.Equal(t, OPERATION_ALL, res)

	s = "write"
	res = tagPartToOperation(s)
	require.Equal(t, OPERATION_WRITE, res)

}
