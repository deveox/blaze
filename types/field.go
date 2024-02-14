package types

import (
	"reflect"
	"strings"

	"github.com/deveox/blaze/scopes"
)

type Field struct {
	Name      string
	ObjectKey []byte

	HTTPScope   Operation
	DBScope     Operation
	ClientScope Operation
	AdminScope  Operation

	KeepEmpty bool
	Anonymous bool
	Struct    *Struct
	Type      reflect.Type
	Idx       int
}

func (f *Field) CheckScope(context scopes.Context, user scopes.User, operation scopes.Operation) bool {
	if context == scopes.CONTEXT_DB {
		ok := f.CheckOperationScope(operation, f.DBScope)
		if !ok {
			return false
		}
	} else {
		ok := f.CheckOperationScope(operation, f.HTTPScope)
		if !ok {
			return false
		}
	}
	if user == scopes.USER_CLIENT {
		return f.CheckOperationScope(operation, f.ClientScope)
	}
	return f.CheckOperationScope(operation, f.AdminScope)
}

func (f *Field) CheckOperationScope(scope scopes.Operation, op Operation) bool {
	switch scope {
	case scopes.OPERATION_READ:
		switch op {
		case OPERATION_READ, OPERATION_READ_CREATE, OPERATION_READ_UPDATE, OPERATION_ALL:
			return true
		}
		return false
	case scopes.OPERATION_WRITE:
		switch op {
		case OPERATION_WRITE, OPERATION_ALL:
			return true
		}
		return false
	case scopes.OPERATION_CREATE:
		switch op {
		case OPERATION_CREATE, OPERATION_READ_CREATE, OPERATION_ALL, OPERATION_WRITE:
			return true
		}
		return false
	case scopes.OPERATION_UPDATE:
		switch op {
		case OPERATION_UPDATE, OPERATION_READ_UPDATE, OPERATION_ALL, OPERATION_WRITE:
			return true
		}
		return false
	default:
		return true
	}
}

func (f *Field) ParseTag(st reflect.StructTag) bool {
	jsonTag := st.Get(TAG)
	tag := st.Get(TAG_BLAZE)
	if jsonTag == "-" || tag == "-" {
		return false
	}
	f.Name, _, _ = strings.Cut(jsonTag, ",")
loop:
	for {
		var v string
		v, tag, _ = strings.Cut(tag, ",")
		switch v {
		case "":
			break loop
		case "keep":
			f.KeepEmpty = true
		case "omit":
			f.KeepEmpty = false
		default:
			s, after, _ := strings.Cut(v, ":")
			switch s {
			case TAG_SCOPE_CLIENT:
				f.ClientScope = tagPartToOperation(after)
			case TAG_SCOPE_ADMIN:
				f.AdminScope = tagPartToOperation(after)
			case TAG_SCOPE_DB:
				f.DBScope = tagPartToOperation(after)
			case TAG_SCOPE_HTTP:
				f.HTTPScope = tagPartToOperation(after)
			default:
				sc := tagPartToOperation(s)
				f.ClientScope = sc
				f.AdminScope = sc
				f.DBScope = sc
				f.HTTPScope = sc
			}
			continue
		}
		if tag == "" {
			break
		}
	}
	return true
}
