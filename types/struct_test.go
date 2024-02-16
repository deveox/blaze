package types

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type NestedStruct struct {
	Name string `blaze:"keep"`
	Age  int    `json:",omitempty" blaze:"admin:read.write"`
}

type TestEmbed struct {
	Name   string `blaze:"client:read"`
	Age    uint   `json:"adminAge" blaze:"admin:write"`
	Nested *NestedStruct
}

type TestStruct struct {
	*TestEmbed
	NoneJSON      string `json:"-"`
	CustomName    int    `json:"adminAge"`
	HARDCamelCase string

	Keep     string `blaze:"keep"`
	NoDB     string `blaze:"no-db"`
	NoHttp   string `blaze:"no-http"`
	NoAdmin  string `blaze:"admin:-"`
	NoClient string `blaze:"client:-"`

	All              string `blaze:"all"`
	None             string `blaze:"-"`
	ReadOnly         string `blaze:"read"`
	UpdateOnly       string `blaze:"update"`
	CreateOnly       string `blaze:"create"`
	ReadCreate       string `blaze:"read.create"`
	ReadUpdate       string `blaze:"read.update"`
	WriteOnly        string `blaze:"write"`
	ReadCreateUpdate string `blaze:"read.create.update"`

	ClientAll              string `blaze:"client:all"`
	ClientNone             string `blaze:"client:-"`
	ClientReadOnly         string `blaze:"client:read"`
	ClientUpdateOnly       string `blaze:"client:update"`
	ClientCreateOnly       string `blaze:"client:create"`
	ClientReadCreate       string `blaze:"client:read.create"`
	ClientReadUpdate       string `blaze:"client:read.update"`
	ClientWriteOnly        string `blaze:"client:write"`
	ClientReadCreateUpdate string `blaze:"client:read.create.update"`

	AdminAll              string `blaze:"admin:all"`
	AdminNone             string `blaze:"admin:-"`
	AdminReadOnly         string `blaze:"admin:read"`
	AdminUpdateOnly       string `blaze:"admin:update"`
	AdminCreateOnly       string `blaze:"admin:create"`
	AdminReadCreate       string `blaze:"admin:read.create"`
	AdminReadUpdate       string `blaze:"admin:read.update"`
	AdminWriteOnly        string `blaze:"admin:write"`
	AdminReadCreateUpdate string `blaze:"admin:read.create.update"`
}

func TestNewStruct(t *testing.T) {
	typ := reflect.TypeOf(TestStruct{})
	s := Cache.Get(typ)

	f, ef, ok := s.GetField("name")
	require.True(t, ok, "name not found")
	require.NotNil(t, ef, "name is not embedded")
	require.NotNil(t, f, "name is nill")
	require.Equal(t, []byte(`"name":`), f.ObjectKey)
	require.Equal(t, OPERATION_READ, f.ClientScope)
	require.Equal(t, OPERATION_ALL, f.AdminScope)

	f, ef, ok = s.GetField("adminAge")
	require.True(t, ok, "adminAge not found")
	require.Nil(t, ef, "adminAge didn't override embedded")
	require.Equal(t, reflect.TypeOf(int(0)), f.Type)
	require.Equal(t, []byte(`"adminAge":`), f.ObjectKey)
	require.Equal(t, 2, f.Idx)

	s2, ok := Cache.load(reflect.TypeOf(NestedStruct{}))
	require.True(t, ok, "nested struct not found")

	f, ef, ok = s.GetField("nested")
	require.True(t, ok, "nested not found")
	require.NotNil(t, ef, "nested is not embedded")
	require.Equal(t, s2, f.Struct)

	f, ef, ok = s.GetField("noneJson")
	require.True(t, ok, "noneJson is not found")
	require.Nil(t, ef, "noneJson is embedded")
	require.Equal(t, OPERATION_IGNORE, f.ClientScope)
	require.Equal(t, OPERATION_IGNORE, f.AdminScope)
	require.False(t, f.DB)

	f, ef, ok = s.GetField("hardCamelCase")
	require.True(t, ok, "hardCamelCase not found")
	require.Nil(t, ef, "hardCamelCase is embedded")
	require.Equal(t, "hardCamelCase", f.Name, "hardCamelCase name is wrong")

	f, ef, ok = s.GetField("keep")
	require.True(t, ok, "keep not found")
	require.Nil(t, ef, "keep is embedded")
	require.True(t, f.KeepEmpty, "keep is not keep empty")

	f, ef, ok = s.GetField("noDb")
	require.True(t, ok, "noDb not found")
	require.Nil(t, ef, "noDb is embedded")
	require.False(t, f.DB, "noDb has DB true")

	f, ef, ok = s.GetField("noHttp")
	require.True(t, ok, "noHttp not found")
	require.Nil(t, ef, "noHttp is embedded")
	require.Equal(t, OPERATION_IGNORE, f.ClientScope, "noHttp client scope is wrong")
	require.Equal(t, OPERATION_IGNORE, f.AdminScope, "noHttp admin scope is wrong")

	f, ef, ok = s.GetField("noAdmin")
	require.True(t, ok, "noAdmin not found")
	require.Nil(t, ef, "noAdmin is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "noAdmin client scope is wrong")
	require.Equal(t, OPERATION_IGNORE, f.AdminScope, "noAdmin admin scope is wrong")

	f, ef, ok = s.GetField("noClient")
	require.True(t, ok, "noClient not found")
	require.Nil(t, ef, "noClient is embedded")
	require.Equal(t, OPERATION_IGNORE, f.ClientScope, "noClient client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "noClient admin scope is wrong")

	f, ef, ok = s.GetField("all")
	require.True(t, ok, "all not found")
	require.Nil(t, ef, "all is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "all client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "all admin scope is wrong")

	f, ef, ok = s.GetField("none")
	require.True(t, ok, "none is found")
	require.Nil(t, ef, "none is embedded")
	require.Equal(t, OPERATION_IGNORE, f.ClientScope, "none client scope is wrong")
	require.Equal(t, OPERATION_IGNORE, f.AdminScope, "none admin scope is wrong")
	require.False(t, f.DB, "none has DB true")

	f, ef, ok = s.GetField("readOnly")
	require.True(t, ok, "readOnly not found")
	require.Nil(t, ef, "readOnly is embedded")
	require.Equal(t, OPERATION_READ, f.ClientScope, "readOnly client scope is wrong")
	require.Equal(t, OPERATION_READ, f.AdminScope, "readOnly admin scope is wrong")

	f, ef, ok = s.GetField("updateOnly")
	require.True(t, ok, "updateOnly not found")
	require.Nil(t, ef, "updateOnly is embedded")
	require.Equal(t, OPERATION_UPDATE, f.ClientScope, "updateOnly client scope is wrong")
	require.Equal(t, OPERATION_UPDATE, f.AdminScope, "updateOnly admin scope is wrong")

	f, ef, ok = s.GetField("createOnly")
	require.True(t, ok, "createOnly not found")
	require.Nil(t, ef, "createOnly is embedded")
	require.Equal(t, OPERATION_CREATE, f.ClientScope, "createOnly client scope is wrong")
	require.Equal(t, OPERATION_CREATE, f.AdminScope, "createOnly admin scope is wrong")

	f, ef, ok = s.GetField("readCreate")
	require.True(t, ok, "readCreate not found")
	require.Nil(t, ef, "readCreate is embedded")
	require.Equal(t, OPERATION_READ_CREATE, f.ClientScope, "readCreate client scope is wrong")
	require.Equal(t, OPERATION_READ_CREATE, f.AdminScope, "readCreate admin scope is wrong")

	f, ef, ok = s.GetField("readUpdate")
	require.True(t, ok, "readUpdate not found")
	require.Nil(t, ef, "readUpdate is embedded")
	require.Equal(t, OPERATION_READ_UPDATE, f.ClientScope, "readUpdate client scope is wrong")
	require.Equal(t, OPERATION_READ_UPDATE, f.AdminScope, "readUpdate admin scope is wrong")

	f, ef, ok = s.GetField("writeOnly")
	require.True(t, ok, "writeOnly not found")
	require.Nil(t, ef, "writeOnly is embedded")
	require.Equal(t, OPERATION_WRITE, f.ClientScope, "writeOnly client scope is wrong")
	require.Equal(t, OPERATION_WRITE, f.AdminScope, "writeOnly admin scope is wrong")

	f, ef, ok = s.GetField("readCreateUpdate")
	require.True(t, ok, "readCreateUpdate not found")
	require.Nil(t, ef, "readCreateUpdate is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "readCreateUpdate client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "readCreateUpdate admin scope is wrong")

	f, ef, ok = s.GetField("clientAll")
	require.True(t, ok, "clientAll not found")
	require.Nil(t, ef, "clientAll is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "clientAll client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "clientAll admin scope is wrong")

	f, ef, ok = s.GetField("clientNone")
	require.True(t, ok, "clientNone not found")
	require.Nil(t, ef, "clientNone is embedded")
	require.Equal(t, OPERATION_IGNORE, f.ClientScope, "clientNone client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "clientNone admin scope is wrong")

	f, ef, ok = s.GetField("clientReadOnly")
	require.True(t, ok, "clientReadOnly not found")
	require.Nil(t, ef, "clientReadOnly is embedded")
	require.Equal(t, OPERATION_READ, f.ClientScope, "clientReadOnly client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "clientReadOnly admin scope is wrong")

	f, ef, ok = s.GetField("clientUpdateOnly")
	require.True(t, ok, "clientUpdateOnly not found")
	require.Nil(t, ef, "clientUpdateOnly is embedded")
	require.Equal(t, OPERATION_UPDATE, f.ClientScope, "clientUpdateOnly client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "clientUpdateOnly admin scope is wrong")

	f, ef, ok = s.GetField("clientCreateOnly")
	require.True(t, ok, "clientCreateOnly not found")
	require.Nil(t, ef, "clientCreateOnly is embedded")
	require.Equal(t, OPERATION_CREATE, f.ClientScope, "clientCreateOnly client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "clientCreateOnly admin scope is wrong")

	f, ef, ok = s.GetField("clientReadCreate")
	require.True(t, ok, "clientReadCreate not found")
	require.Nil(t, ef, "clientReadCreate is embedded")
	require.Equal(t, OPERATION_READ_CREATE, f.ClientScope, "clientReadCreate client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "clientReadCreate admin scope is wrong")

	f, ef, ok = s.GetField("clientReadUpdate")
	require.True(t, ok, "clientReadUpdate not found")
	require.Nil(t, ef, "clientReadUpdate is embedded")
	require.Equal(t, OPERATION_READ_UPDATE, f.ClientScope, "clientReadUpdate client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "clientReadUpdate admin scope is wrong")

	f, ef, ok = s.GetField("clientWriteOnly")
	require.True(t, ok, "clientWriteOnly not found")
	require.Nil(t, ef, "clientWriteOnly is embedded")
	require.Equal(t, OPERATION_WRITE, f.ClientScope, "clientWriteOnly client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "clientWriteOnly admin scope is wrong")

	f, ef, ok = s.GetField("clientReadCreateUpdate")
	require.True(t, ok, "clientReadCreateUpdate not found")
	require.Nil(t, ef, "clientReadCreateUpdate is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "clientReadCreateUpdate client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "clientReadCreateUpdate admin scope is wrong")

	f, ef, ok = s.GetField("adminAll")
	require.True(t, ok, "adminAll not found")
	require.Nil(t, ef, "adminAll is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "adminAll client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "adminAll admin scope is wrong")

	f, ef, ok = s.GetField("adminNone")
	require.True(t, ok, "adminNone not found")
	require.Nil(t, ef, "adminNone is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "adminNone client scope is wrong")
	require.Equal(t, OPERATION_IGNORE, f.AdminScope, "adminNone admin scope is wrong")

	f, ef, ok = s.GetField("adminReadOnly")
	require.True(t, ok, "adminReadOnly not found")
	require.Nil(t, ef, "adminReadOnly is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "adminReadOnly client scope is wrong")
	require.Equal(t, OPERATION_READ, f.AdminScope, "adminReadOnly admin scope is wrong")

	f, ef, ok = s.GetField("adminUpdateOnly")
	require.True(t, ok, "adminUpdateOnly not found")
	require.Nil(t, ef, "adminUpdateOnly is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "adminUpdateOnly client scope is wrong")
	require.Equal(t, OPERATION_UPDATE, f.AdminScope, "adminUpdateOnly admin scope is wrong")

	f, ef, ok = s.GetField("adminCreateOnly")
	require.True(t, ok, "adminCreateOnly not found")
	require.Nil(t, ef, "adminCreateOnly is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "adminCreateOnly client scope is wrong")
	require.Equal(t, OPERATION_CREATE, f.AdminScope, "adminCreateOnly admin scope is wrong")

	f, ef, ok = s.GetField("adminReadCreate")
	require.True(t, ok, "adminReadCreate not found")
	require.Nil(t, ef, "adminReadCreate is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "adminReadCreate client scope is wrong")
	require.Equal(t, OPERATION_READ_CREATE, f.AdminScope, "adminReadCreate admin scope is wrong")

	f, ef, ok = s.GetField("adminReadUpdate")
	require.True(t, ok, "adminReadUpdate not found")
	require.Nil(t, ef, "adminReadUpdate is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "adminReadUpdate client scope is wrong")
	require.Equal(t, OPERATION_READ_UPDATE, f.AdminScope, "adminReadUpdate admin scope is wrong")

	f, ef, ok = s.GetField("adminWriteOnly")
	require.True(t, ok, "adminWriteOnly not found")
	require.Nil(t, ef, "adminWriteOnly is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "adminWriteOnly client scope is wrong")
	require.Equal(t, OPERATION_WRITE, f.AdminScope, "adminWriteOnly admin scope is wrong")

	f, ef, ok = s.GetField("adminReadCreateUpdate")
	require.True(t, ok, "adminReadCreateUpdate not found")
	require.Nil(t, ef, "adminReadCreateUpdate is embedded")
	require.Equal(t, OPERATION_ALL, f.ClientScope, "adminReadCreateUpdate client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.AdminScope, "adminReadCreateUpdate admin scope is wrong")

}
