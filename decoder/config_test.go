package decoder

import (
	"testing"

	"github.com/deveox/blaze/scopes"
	"github.com/deveox/blaze/testdata"
	"github.com/stretchr/testify/require"
)

var dbDecoder = &Config{
	ContextScope: scopes.CONTEXT_DB,
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
