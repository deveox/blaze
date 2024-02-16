package encoder

import (
	"testing"

	"github.com/deveox/blaze/scopes"
	"github.com/deveox/blaze/testdata"
	"github.com/stretchr/testify/require"
)

var dbEncoder = &Config{
	Scope: scopes.CONTEXT_DB,
}

func newScopedStruct() *testdata.ScopedStruct {
	return &testdata.ScopedStruct{
		Name:             "test",
		NoDB:             true,
		Read:             true,
		ReadCreate:       true,
		ReadUpdate:       true,
		Update:           true,
		Create:           true,
		Write:            true,
		NoClient:         true,
		ClientRead:       true,
		ClientReadCreate: true,
		ClientReadUpdate: true,
		ClientUpdate:     true,
		ClientCreate:     true,
		ClientWrite:      true,
		NoAdmin:          true,
		AdminRead:        true,
		AdminReadCreate:  true,
		AdminReadUpdate:  true,
		AdminUpdate:      true,
		AdminCreate:      true,
		AdminWrite:       true,
	}
}

func TestScope_DB(t *testing.T) {
	s := newScopedStruct()
	res, _ := dbEncoder.Marshal(s)
	res = AddIndent(res)
	expected := AddIndent([]byte(`{"name":"test","read":true,"readCreate":true,"readUpdate":true,"update":true,"create":true,"write":true,"noClient":true,"clientRead":true,"clientReadCreate":true,"clientReadUpdate":true,"clientUpdate":true,"clientCreate":true,"clientWrite":true,"noAdmin":true,"adminRead":true,"adminReadCreate":true,"adminReadUpdate":true,"adminUpdate":true,"adminCreate":true,"adminWrite":true}`))
	require.Equal(t, string(expected), string(res))
}

var clientEncoder = &Config{
	Scope: scopes.CONTEXT_CLIENT,
}

func TestScope_Client(t *testing.T) {
	s := newScopedStruct()
	res, _ := clientEncoder.Marshal(s)
	res = AddIndent(res)
	expected := AddIndent([]byte(`{"name":"test","noDb":true,"read":true,"readCreate":true,"readUpdate":true,"clientRead":true,"clientReadCreate":true,"clientReadUpdate":true,"noAdmin":true,"adminRead":true,"adminReadCreate":true,"adminReadUpdate":true,"adminUpdate":true,"adminCreate":true,"adminWrite":true}`))
	require.Equal(t, string(expected), string(res))
}

var adminEncoder = &Config{
	Scope: scopes.CONTEXT_ADMIN,
}

func TestScope_Admin(t *testing.T) {
	s := newScopedStruct()
	res, _ := adminEncoder.Marshal(s)
	res = AddIndent(res)
	expected := AddIndent([]byte(`{"name":"test","noDb":true,"read":true,"readCreate":true,"readUpdate":true,"noClient":true,"clientRead":true,"clientReadCreate":true,"clientReadUpdate":true,"clientUpdate":true,"clientCreate":true,"clientWrite":true,"adminRead":true,"adminReadCreate":true,"adminReadUpdate":true}`))
	require.Equal(t, string(expected), string(res))
}
