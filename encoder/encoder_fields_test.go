package encoder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFields(t *testing.T) {
	f := fields{}
	f.Init([]string{"a", "b", "c", "nested.a.b.c.d"}, false)
	require.Equal(t, []string{"a", "b", "c", "nested.a.b.c.d", "nested", "nested.a", "nested.a.b", "nested.a.b.c"}, f.fields)
}
