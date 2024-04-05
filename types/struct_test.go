package types

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type NestedStruct struct {
	Name string `blaze:"keep"`
	Age  int    `json:",omitempty" blaze:"admin:read.write"`
}

type TestEmbed struct {
	Name    string `blaze:"client:read"`
	Age     uint   `json:"adminAge" blaze:"admin:write"`
	Nested  *NestedStruct
	Created time.Time `json:"created"`
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

	f, ok := s.GetField("name")
	require.True(t, ok, "name not found")
	require.NotNil(t, f, "name is nill")
	require.Equal(t, []byte(`"name":`), f.Field.ObjectKey)
	require.Equal(t, OPERATION_READ, f.Field.ClientScope)
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope)
	require.Equal(t, []int{0, 0}, f.Idx)

	f, ok = s.GetField("created")
	require.True(t, ok, "created not found")
	require.Equal(t, reflect.TypeOf(time.Time{}), f.Field.Type)
	require.Equal(t, []byte(`"created":`), f.Field.ObjectKey)

	f, ok = s.GetField("adminAge")
	require.True(t, ok, "adminAge not found")
	require.Equal(t, reflect.TypeOf(int(0)), f.Field.Type)
	require.Equal(t, []byte(`"adminAge":`), f.Field.ObjectKey)

	s2, ok := Cache.load(reflect.TypeOf(NestedStruct{}))
	require.True(t, ok, "nested struct not found")

	f, ok = s.GetField("nested")
	require.True(t, ok, "nested not found")
	require.Equal(t, s2, f.Field.Struct)

	f, ok = s.GetField("noneJson")
	require.True(t, ok, "noneJson is not found")
	require.Equal(t, OPERATION_IGNORE, f.Field.ClientScope)
	require.Equal(t, OPERATION_IGNORE, f.Field.AdminScope)
	require.False(t, f.Field.DBScope)

	f, ok = s.GetField("hardCamelCase")
	require.True(t, ok, "hardCamelCase not found")
	require.Equal(t, "hardCamelCase", f.Field.Name, "hardCamelCase name is wrong")

	f, ok = s.GetField("keep")
	require.True(t, ok, "keep not found")
	require.True(t, f.Field.KeepEmpty, "keep is not keep empty")

	f, ok = s.GetField("noDb")
	require.True(t, ok, "noDb not found")
	require.False(t, f.Field.DBScope, "noDb has DB true")

	f, ok = s.GetField("noHttp")
	require.True(t, ok, "noHttp not found")
	require.Equal(t, OPERATION_IGNORE, f.Field.ClientScope, "noHttp client scope is wrong")
	require.Equal(t, OPERATION_IGNORE, f.Field.AdminScope, "noHttp admin scope is wrong")

	f, ok = s.GetField("noAdmin")
	require.True(t, ok, "noAdmin not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "noAdmin client scope is wrong")
	require.Equal(t, OPERATION_IGNORE, f.Field.AdminScope, "noAdmin admin scope is wrong")

	f, ok = s.GetField("noClient")
	require.True(t, ok, "noClient not found")
	require.Equal(t, OPERATION_IGNORE, f.Field.ClientScope, "noClient client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "noClient admin scope is wrong")

	f, ok = s.GetField("all")
	require.True(t, ok, "all not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "all client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "all admin scope is wrong")

	f, ok = s.GetField("none")
	require.True(t, ok, "none is found")
	require.Equal(t, OPERATION_IGNORE, f.Field.ClientScope, "none client scope is wrong")
	require.Equal(t, OPERATION_IGNORE, f.Field.AdminScope, "none admin scope is wrong")
	require.False(t, f.Field.DBScope, "none has DB true")

	f, ok = s.GetField("readOnly")
	require.True(t, ok, "readOnly not found")
	require.Equal(t, OPERATION_READ, f.Field.ClientScope, "readOnly client scope is wrong")
	require.Equal(t, OPERATION_READ, f.Field.AdminScope, "readOnly admin scope is wrong")

	f, ok = s.GetField("updateOnly")
	require.True(t, ok, "updateOnly not found")
	require.Equal(t, OPERATION_UPDATE, f.Field.ClientScope, "updateOnly client scope is wrong")
	require.Equal(t, OPERATION_UPDATE, f.Field.AdminScope, "updateOnly admin scope is wrong")

	f, ok = s.GetField("createOnly")
	require.True(t, ok, "createOnly not found")
	require.Equal(t, OPERATION_CREATE, f.Field.ClientScope, "createOnly client scope is wrong")
	require.Equal(t, OPERATION_CREATE, f.Field.AdminScope, "createOnly admin scope is wrong")

	f, ok = s.GetField("readCreate")
	require.True(t, ok, "readCreate not found")
	require.Equal(t, OPERATION_READ_CREATE, f.Field.ClientScope, "readCreate client scope is wrong")
	require.Equal(t, OPERATION_READ_CREATE, f.Field.AdminScope, "readCreate admin scope is wrong")

	f, ok = s.GetField("readUpdate")
	require.True(t, ok, "readUpdate not found")
	require.Equal(t, OPERATION_READ_UPDATE, f.Field.ClientScope, "readUpdate client scope is wrong")
	require.Equal(t, OPERATION_READ_UPDATE, f.Field.AdminScope, "readUpdate admin scope is wrong")

	f, ok = s.GetField("writeOnly")
	require.True(t, ok, "writeOnly not found")
	require.Equal(t, OPERATION_WRITE, f.Field.ClientScope, "writeOnly client scope is wrong")
	require.Equal(t, OPERATION_WRITE, f.Field.AdminScope, "writeOnly admin scope is wrong")

	f, ok = s.GetField("readCreateUpdate")
	require.True(t, ok, "readCreateUpdate not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "readCreateUpdate client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "readCreateUpdate admin scope is wrong")

	f, ok = s.GetField("clientAll")
	require.True(t, ok, "clientAll not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "clientAll client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "clientAll admin scope is wrong")

	f, ok = s.GetField("clientNone")
	require.True(t, ok, "clientNone not found")
	require.Equal(t, OPERATION_IGNORE, f.Field.ClientScope, "clientNone client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "clientNone admin scope is wrong")

	f, ok = s.GetField("clientReadOnly")
	require.True(t, ok, "clientReadOnly not found")
	require.Equal(t, OPERATION_READ, f.Field.ClientScope, "clientReadOnly client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "clientReadOnly admin scope is wrong")

	f, ok = s.GetField("clientUpdateOnly")
	require.True(t, ok, "clientUpdateOnly not found")
	require.Equal(t, OPERATION_UPDATE, f.Field.ClientScope, "clientUpdateOnly client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "clientUpdateOnly admin scope is wrong")

	f, ok = s.GetField("clientCreateOnly")
	require.True(t, ok, "clientCreateOnly not found")
	require.Equal(t, OPERATION_CREATE, f.Field.ClientScope, "clientCreateOnly client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "clientCreateOnly admin scope is wrong")

	f, ok = s.GetField("clientReadCreate")
	require.True(t, ok, "clientReadCreate not found")
	require.Equal(t, OPERATION_READ_CREATE, f.Field.ClientScope, "clientReadCreate client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "clientReadCreate admin scope is wrong")

	f, ok = s.GetField("clientReadUpdate")
	require.True(t, ok, "clientReadUpdate not found")
	require.Equal(t, OPERATION_READ_UPDATE, f.Field.ClientScope, "clientReadUpdate client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "clientReadUpdate admin scope is wrong")

	f, ok = s.GetField("clientWriteOnly")
	require.True(t, ok, "clientWriteOnly not found")
	require.Equal(t, OPERATION_WRITE, f.Field.ClientScope, "clientWriteOnly client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "clientWriteOnly admin scope is wrong")

	f, ok = s.GetField("clientReadCreateUpdate")
	require.True(t, ok, "clientReadCreateUpdate not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "clientReadCreateUpdate client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "clientReadCreateUpdate admin scope is wrong")

	f, ok = s.GetField("adminAll")
	require.True(t, ok, "adminAll not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "adminAll client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "adminAll admin scope is wrong")

	f, ok = s.GetField("adminNone")
	require.True(t, ok, "adminNone not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "adminNone client scope is wrong")
	require.Equal(t, OPERATION_IGNORE, f.Field.AdminScope, "adminNone admin scope is wrong")

	f, ok = s.GetField("adminReadOnly")
	require.True(t, ok, "adminReadOnly not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "adminReadOnly client scope is wrong")
	require.Equal(t, OPERATION_READ, f.Field.AdminScope, "adminReadOnly admin scope is wrong")

	f, ok = s.GetField("adminUpdateOnly")
	require.True(t, ok, "adminUpdateOnly not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "adminUpdateOnly client scope is wrong")
	require.Equal(t, OPERATION_UPDATE, f.Field.AdminScope, "adminUpdateOnly admin scope is wrong")

	f, ok = s.GetField("adminCreateOnly")
	require.True(t, ok, "adminCreateOnly not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "adminCreateOnly client scope is wrong")
	require.Equal(t, OPERATION_CREATE, f.Field.AdminScope, "adminCreateOnly admin scope is wrong")

	f, ok = s.GetField("adminReadCreate")
	require.True(t, ok, "adminReadCreate not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "adminReadCreate client scope is wrong")
	require.Equal(t, OPERATION_READ_CREATE, f.Field.AdminScope, "adminReadCreate admin scope is wrong")

	f, ok = s.GetField("adminReadUpdate")
	require.True(t, ok, "adminReadUpdate not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "adminReadUpdate client scope is wrong")
	require.Equal(t, OPERATION_READ_UPDATE, f.Field.AdminScope, "adminReadUpdate admin scope is wrong")

	f, ok = s.GetField("adminWriteOnly")
	require.True(t, ok, "adminWriteOnly not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "adminWriteOnly client scope is wrong")
	require.Equal(t, OPERATION_WRITE, f.Field.AdminScope, "adminWriteOnly admin scope is wrong")

	f, ok = s.GetField("adminReadCreateUpdate")
	require.True(t, ok, "adminReadCreateUpdate not found")
	require.Equal(t, OPERATION_ALL, f.Field.ClientScope, "adminReadCreateUpdate client scope is wrong")
	require.Equal(t, OPERATION_ALL, f.Field.AdminScope, "adminReadCreateUpdate admin scope is wrong")

}
