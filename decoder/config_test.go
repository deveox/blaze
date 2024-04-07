package decoder

import (
	"testing"

	"github.com/deveox/blaze/internal/testdata"
	"github.com/deveox/blaze/scopes"
	"github.com/stretchr/testify/require"
)

var dbDecoder = &Config{
	Scope: scopes.CONTEXT_DB,
}

func newScopedStructJSON() []byte {
	return []byte(`{
		"name": "test",
		"noDb": true,
		"read": true,
		"readCreate": true,
		"readUpdate": true,
		"update": true,
		"create": true,
		"write": true,
		"noClient": true,
		"clientRead": true,
		"clientReadCreate": true,
		"clientReadUpdate": true,
		"clientUpdate": true,
		"clientCreate": true,
		"clientWrite": true,
		"noAdmin": true,
		"adminRead": true,
		"adminReadCreate": true,
		"adminReadUpdate": true,
		"adminUpdate": true,
		"adminCreate": true,
		"adminWrite": true
	}`)
}

func TestScope_DB(t *testing.T) {
	data := newScopedStructJSON()
	var s testdata.ScopedStruct
	err := dbDecoder.Unmarshal(data, &s)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, testdata.ScopedStruct{
		Name:             "test",
		NoDB:             false,
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
	}, s)
}

var adminDecoder = &Config{}

func TestScope_Admin(t *testing.T) {
	data := newScopedStructJSON()
	var s testdata.ScopedStruct
	err := adminDecoder.Unmarshal(data, &s)
	require.NoError(t, err)
	require.Equal(t, testdata.ScopedStruct{
		Name:             "test",
		NoDB:             true,
		Read:             false,
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
		NoAdmin:          false,
		AdminRead:        false,
		AdminReadCreate:  true,
		AdminReadUpdate:  true,
		AdminUpdate:      true,
		AdminCreate:      true,
		AdminWrite:       true,
	}, s)
}

func TestScope_Admin_Create(t *testing.T) {
	data := newScopedStructJSON()
	var s testdata.ScopedStruct
	err := adminDecoder.UnmarshalScoped(data, &s, scopes.DECODE_CREATE)
	require.NoError(t, err)
	require.Equal(t, testdata.ScopedStruct{
		Name:             "test",
		NoDB:             true,
		Read:             false,
		ReadCreate:       true,
		ReadUpdate:       false,
		Update:           false,
		Create:           true,
		Write:            true,
		NoClient:         true,
		ClientRead:       true,
		ClientReadCreate: true,
		ClientReadUpdate: true,
		ClientUpdate:     true,
		ClientCreate:     true,
		ClientWrite:      true,
		NoAdmin:          false,
		AdminRead:        false,
		AdminReadCreate:  true,
		AdminReadUpdate:  false,
		AdminUpdate:      false,
		AdminCreate:      true,
		AdminWrite:       true,
	}, s)

}

func TestScope_Admin_Update(t *testing.T) {
	data := newScopedStructJSON()
	var s testdata.ScopedStruct
	err := adminDecoder.UnmarshalScoped(data, &s, scopes.DECODE_UPDATE)
	require.NoError(t, err)
	require.Equal(t, testdata.ScopedStruct{
		Name:             "test",
		NoDB:             true,
		Read:             false,
		ReadCreate:       false,
		ReadUpdate:       true,
		Update:           true,
		Create:           false,
		Write:            true,
		NoClient:         true,
		ClientRead:       true,
		ClientReadCreate: true,
		ClientReadUpdate: true,
		ClientUpdate:     true,
		ClientCreate:     true,
		ClientWrite:      true,
		NoAdmin:          false,
		AdminRead:        false,
		AdminReadCreate:  false,
		AdminReadUpdate:  true,
		AdminUpdate:      true,
		AdminCreate:      false,
		AdminWrite:       true,
	}, s)

}

var clientDecoder = &Config{
	Scope: scopes.CONTEXT_CLIENT,
}

func TestScope_Client(t *testing.T) {
	data := newScopedStructJSON()
	var s testdata.ScopedStruct
	err := clientDecoder.Unmarshal(data, &s)
	require.NoError(t, err)
	require.Equal(t, testdata.ScopedStruct{
		Name:             "test",
		NoDB:             true,
		Read:             false,
		ReadCreate:       true,
		ReadUpdate:       true,
		Update:           true,
		Create:           true,
		Write:            true,
		NoClient:         false,
		ClientRead:       false,
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
	}, s)
}

func TestScope_Client_Create(t *testing.T) {
	data := newScopedStructJSON()
	var s testdata.ScopedStruct
	err := clientDecoder.UnmarshalScoped(data, &s, scopes.DECODE_CREATE)
	require.NoError(t, err)
	require.Equal(t, testdata.ScopedStruct{
		Name:             "test",
		NoDB:             true,
		Read:             false,
		ReadCreate:       true,
		ReadUpdate:       false,
		Update:           false,
		Create:           true,
		Write:            true,
		NoClient:         false,
		ClientRead:       false,
		ClientReadCreate: true,
		ClientReadUpdate: false,
		ClientUpdate:     false,
		ClientCreate:     true,
		ClientWrite:      true,
		NoAdmin:          true,
		AdminRead:        true,
		AdminReadCreate:  true,
		AdminReadUpdate:  true,
		AdminUpdate:      true,
		AdminCreate:      true,
		AdminWrite:       true,
	}, s)

}

func TestScope_Client_Update(t *testing.T) {
	data := newScopedStructJSON()
	var s testdata.ScopedStruct
	err := clientDecoder.UnmarshalScoped(data, &s, scopes.DECODE_UPDATE)
	require.NoError(t, err)
	require.Equal(t, testdata.ScopedStruct{
		Name:             "test",
		NoDB:             true,
		Read:             false,
		ReadCreate:       false,
		ReadUpdate:       true,
		Update:           true,
		Create:           false,
		Write:            true,
		NoClient:         false,
		ClientRead:       false,
		ClientReadCreate: false,
		ClientReadUpdate: true,
		ClientUpdate:     true,
		ClientCreate:     false,
		ClientWrite:      true,
		NoAdmin:          true,
		AdminRead:        true,
		AdminReadCreate:  true,
		AdminReadUpdate:  true,
		AdminUpdate:      true,
		AdminCreate:      true,
		AdminWrite:       true,
	}, s)
}
