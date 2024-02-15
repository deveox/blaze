package decoder

import (
	"testing"

	"github.com/deveox/blaze/scopes"
)

var dbDecoder = &Config{
	ContextScope: scopes.CONTEXT_DB,
}

func TestScope_DB(t *testing.T) {
	data := []byte(`{"name":"test", "no_db":true, "no_client":true, "no_admin":true}`)
	var s scopeStruct
	err := dbDecoder.Unmarshal(data, &s)
	if err != nil {
		t.Fatal(err)
	}

}
