package types

import (
	"reflect"
	"strings"

	"github.com/deveox/blaze/scopes"
)

type StructField struct {
	Field *Field
	Idx   []int
}

func (e *StructField) Value(v reflect.Value) reflect.Value {
	for _, i := range e.Idx {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		}
		v = v.Field(i)
	}
	return v
}

type Field struct {
	Name      string
	ObjectKey []byte

	DB          bool
	ClientScope Operation
	AdminScope  Operation

	KeepEmpty bool
	Struct    *Struct
	Type      reflect.Type
	Short     bool
}

func (f *Field) CheckEncoderScope(context scopes.Context) bool {
	switch context {
	case scopes.CONTEXT_DB:
		return f.DB
	case scopes.CONTEXT_CLIENT:
		return f.ClientScope.CanRead()
	case scopes.CONTEXT_ADMIN:
		return f.AdminScope.CanRead()
	}
	return false
}

func (f *Field) CheckDecoderScope(context scopes.Context, scope scopes.Decoding) bool {
	switch context {
	case scopes.CONTEXT_DB:
		return f.DB
	case scopes.CONTEXT_CLIENT:
		return f.ClientScope.CanWrite(scope)
	case scopes.CONTEXT_ADMIN:
		return f.AdminScope.CanWrite(scope)
	}
	return false
}

func (f *Field) ParseTag(st reflect.StructTag) {
	jsonTag := st.Get(TAG_NAME_JSON)
	tag := st.Get(TAG_NAME_BLAZE)
	if jsonTag == "-" || tag == "-" {
		f.ClientScope = OPERATION_IGNORE
		f.AdminScope = OPERATION_IGNORE
		f.DB = false
		return
	}
	f.DB = true
	f.Name, _, _ = strings.Cut(jsonTag, ",")
loop:
	for {
		var v string
		v, tag, _ = strings.Cut(tag, ",")
		switch v {
		case "":
			break loop
		case TAG_SHORT:
			f.Short = true
		case TAG_KEEP:
			f.KeepEmpty = true
		case "omit":
			f.KeepEmpty = false
		case TAG_NO_DB:
			f.DB = false
		case TAG_NO_HTTP:
			f.ClientScope = OPERATION_IGNORE
			f.AdminScope = OPERATION_IGNORE
		default:
			s, after, _ := strings.Cut(v, ":")
			switch s {
			case TAG_SCOPE_CLIENT:
				f.ClientScope = tagPartToOperation(after)
			case TAG_SCOPE_ADMIN:
				f.AdminScope = tagPartToOperation(after)
			default:
				sc := tagPartToOperation(s)
				f.ClientScope = sc
				f.AdminScope = sc

			}
			continue
		}
		if tag == "" {
			break
		}
	}
}