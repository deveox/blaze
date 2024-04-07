package types

import (
	"reflect"
	"strings"

	"github.com/deveox/blaze/scopes"
	"github.com/deveox/gu/stringer"
)

// GetDBName is a callback function that returns the database name of a field.
// By default, it searches for the "column" tag in the field's gorm tag. If it is not found, it converts the [Field.Name] to snake case.
var GetDBName = func(f reflect.StructField, fi *StructField) string {
	gorm := f.Tag.Get("gorm")
	if gorm != "" {
		for {
			var s string
			s, gorm, _ = strings.Cut(gorm, ";")
			col, name, _ := strings.Cut(s, ":")
			if col == "column" && name != "" {
				return name
			}
			if gorm == "" {
				break
			}
		}
	}
	return stringer.ToSnakeCase(f.Name)
}

// StructField represents a field instance in a particular struct.
type StructField struct {
	// Link to the actual field definition. Each field stored once and shared between all structs.
	Field *Field
	// Shows if the field is anonymous (doesn't have name).
	Anonymous bool
	// Reflect path to the field in the struct. Usually there's only one index, but in case of anonymous structs, there can be more.
	Idx []int
}

func (e *StructField) PostgreSQLType() string {
	switch e.Field.Kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "::INTEGER"
	case reflect.Int64, reflect.Uint64:
		return "::BIGINT"
	case reflect.Float32, reflect.Float64:
		return "::NUMERIC"
	case reflect.Bool:
		return "::BOOLEAN"
	default:
		return "::TEXT"
	}
}

// Value returns the [reflect.Value] of the field in the given struct.
// Accepts a struct [reflect.Value].
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

// Field represents a meta info about field in a struct.
type Field struct {
	// The native go name of the field.
	TitleCase string
	// The name of the field in the JSON representation. If not set in `json` tag, it will be converted from [Field.TitleCase] to camelCase.
	Name string
	// Precomputed key for the field in the json object.
	ObjectKey []byte

	// Defines if the field should be included in the database JSON.
	DBScope bool
	// Defines if the field should be included in the client marshaling/unmarshaling
	ClientScope Operation
	// Defines if the field should be included in the admin marshaling/unmarshaling
	AdminScope Operation

	// Defines if the field should be kept in the JSON object even if it's empty.
	KeepEmpty bool
	// If the field is a struct, then it will point to the struct definition. Otherwise, it will be nil.
	Struct *Struct
	Type   reflect.Type
	Kind   reflect.Kind
	// Defines if the field should be marshaled as a short version.
	Short bool
	// The database name of the field. Populated by [GetDBName] function.
	DBName string
}

// CheckEncoderScope checks if the field can be encoded in the given context.
func (f *Field) CheckEncoderScope(context scopes.Context) bool {
	switch context {
	case scopes.CONTEXT_DB:
		return f.DBScope
	case scopes.CONTEXT_CLIENT:
		return f.ClientScope.CanRead()
	case scopes.CONTEXT_ADMIN:
		return f.AdminScope.CanRead()
	}
	return false
}

// CheckDecoderScope checks if the field can be decoded in the given context.
func (f *Field) CheckDecoderScope(context scopes.Context, scope scopes.Decoding) bool {
	switch context {
	case scopes.CONTEXT_DB:
		return f.DBScope
	case scopes.CONTEXT_CLIENT:
		return f.ClientScope.CanWrite(scope)
	case scopes.CONTEXT_ADMIN:
		return f.AdminScope.CanWrite(scope)
	}
	return false
}

// ParseTag parses the struct tag and populates the field with the data.
func (f *Field) ParseTag(st reflect.StructTag) {
	jsonTag := st.Get(TAG_NAME_JSON)
	tag := st.Get(TAG_NAME_BLAZE)
	if jsonTag == "-" || tag == "-" {
		f.ClientScope = OPERATION_IGNORE
		f.AdminScope = OPERATION_IGNORE
		f.DBScope = false
		return
	}
	f.DBScope = true
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
			f.DBScope = false
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
