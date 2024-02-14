package types

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/deveox/blaze/scopes"

	"github.com/stretchr/testify/require"
)

type NestedStruct struct {
	Name string `blaze:"keep"`
	Age  int    `json:",omitempty" blaze:"admin:read.write"`
}

type TestEmbed struct {
	Name   string `blaze:"client:read"`
	Age    int    `json:"adminAge" blaze:"admin:write"`
	Nested *NestedStruct
}

type TestStruct struct {
	Name string `blaze:"update"`
	Age  int    `json:"adminAge" blaze:"create"`
	Role string `json:"role" blaze:"read,keep,admin:create"`
	*TestEmbed
}

func TestNewStruct(t *testing.T) {
	fmt.Println("st")
	typ := reflect.TypeOf(TestStruct{})
	fmt.Println("typ", typ)
	s, err := Cache.Get(typ)
	require.NoError(t, err)
	require.Len(t, s.Fields, 4)

	require.Equal(t, "name", s.Fields[0].Name)
	require.Equal(t, scopes.OPERATION_UPDATE.String(), s.Fields[0].ClientScope.String(), "name has wrong client scope")
	require.Equal(t, scopes.OPERATION_UPDATE.String(), s.Fields[0].AdminScope.String(), "name has wrong admin scope")
	require.Equal(t, false, s.Fields[0].KeepEmpty)

	require.Equal(t, "adminAge", s.Fields[1].Name)
	require.Equal(t, scopes.OPERATION_CREATE.String(), s.Fields[1].ClientScope.String(), "adminAge has wrong client scope")
	require.Equal(t, scopes.OPERATION_CREATE.String(), s.Fields[1].AdminScope.String(), "adminAge has wrong admin scope")
	require.Equal(t, false, s.Fields[1].KeepEmpty)

	require.Equal(t, "role", s.Fields[2].Name)
	require.Equal(t, scopes.OPERATION_READ.String(), s.Fields[2].ClientScope.String(), "role has wrong client scope")
	require.Equal(t, scopes.OPERATION_CREATE.String(), s.Fields[2].AdminScope.String(), "role has wrong admin scope")
	require.Equal(t, true, s.Fields[2].KeepEmpty)

	require.Equal(t, "testEmbed", s.Fields[3].Name)
	require.Equal(t, scopes.OPERATION_ALL.String(), s.Fields[3].ClientScope.String(), "testEmbed has wrong client scope")
	require.Equal(t, scopes.OPERATION_ALL.String(), s.Fields[3].AdminScope.String(), "testEmbed has wrong admin scope")
	require.Equal(t, false, s.Fields[3].KeepEmpty)

	em, ok := Cache.load(reflect.TypeOf(TestEmbed{}))
	require.True(t, ok, "TestEmbed struct not found in cache")
	require.Len(t, em.Fields, 3)

	require.Equal(t, "name", em.Fields[0].Name)
	require.Equal(t, scopes.OPERATION_READ.String(), em.Fields[0].ClientScope.String(), "name has wrong client scope")
	require.Equal(t, scopes.OPERATION_ALL.String(), em.Fields[0].AdminScope.String(), "name has wrong admin scope")
	require.Equal(t, false, em.Fields[0].KeepEmpty)

	require.Equal(t, "adminAge", em.Fields[1].Name)
	require.Equal(t, scopes.OPERATION_ALL.String(), em.Fields[1].ClientScope.String(), "adminAge has wrong client scope")
	require.Equal(t, scopes.OPERATION_WRITE.String(), em.Fields[1].AdminScope.String(), "adminAge has wrong admin scope")
	require.Equal(t, false, em.Fields[1].KeepEmpty)

	require.Equal(t, "nested", em.Fields[2].Name)
	require.Equal(t, scopes.OPERATION_ALL.String(), em.Fields[2].ClientScope.String(), "nested has wrong client scope")
	require.Equal(t, scopes.OPERATION_ALL.String(), em.Fields[2].AdminScope.String(), "nested has wrong admin scope")
	require.Equal(t, false, em.Fields[2].KeepEmpty)

	nested, ok := Cache.load(reflect.TypeOf(NestedStruct{}))

	require.True(t, ok, "NestedStruct struct not found in cache")
	require.Len(t, nested.Fields, 2)

	require.Equal(t, "name", nested.Fields[0].Name)
	require.Equal(t, scopes.OPERATION_ALL.String(), nested.Fields[0].ClientScope.String(), "name has wrong client scope")
	require.Equal(t, scopes.OPERATION_ALL.String(), nested.Fields[0].AdminScope.String(), "name has wrong admin scope")
	require.Equal(t, true, nested.Fields[0].KeepEmpty)

	require.Equal(t, "age", nested.Fields[1].Name)
	require.Equal(t, scopes.OPERATION_ALL.String(), nested.Fields[1].ClientScope.String(), "age has wrong client scope")
	require.Equal(t, scopes.OPERATION_ALL.String(), nested.Fields[1].AdminScope.String(), "age has wrong admin scope")
	require.Equal(t, false, nested.Fields[1].KeepEmpty)
}
