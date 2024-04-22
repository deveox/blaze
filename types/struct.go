package types

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/deveox/blaze/scopes"

	"github.com/deveox/gu/mirror"
	"github.com/deveox/gu/stringer"
)

// Struct represents a type meta information of a struct.
type Struct struct {
	Type reflect.Type
	// Fields is a list of fields in the struct.
	Fields      []*StructField
	byCamelName map[string]*StructField
}

// GetField returns a field by its name. If the field is not found, the second return value is [false].
func (c *Struct) GetField(name string) (*StructField, bool) {
	v, ok := c.byCamelName[name]
	return v, ok
}

type NestedField struct {
	Path   string
	GoPath string
	Field  *StructField
}

// GetNestedFields returns a list of all fields in the struct and its nested structs. Paths are dot-separated, e.g. "user.address.phoneNumber".
func (c *Struct) GetNestedFields(ctx scopes.Context, decoding scopes.Decoding) []NestedField {
	fields := make([]NestedField, 0, len(c.Fields))
	for _, f := range c.Fields {
		if decoding != scopes.DECODE_ANY && !f.Field.CheckDecoderScope(ctx, decoding) {
			continue
		}
		if !f.Field.CheckEncoderScope(ctx) {
			continue
		}
		if f.Field.Struct != nil && !f.Anonymous {
			ff := f.Field.Struct.GetNestedFields(ctx, decoding)
			for _, nf := range ff {
				fields = append(fields, NestedField{
					Path:   fmt.Sprintf("%s.%s", f.Field.Name, nf.Path),
					GoPath: fmt.Sprintf("%s.%s", f.Field.TitleCase, nf.GoPath),
					Field:  nf.Field,
				})
			}
		} else {
			fields = append(fields, NestedField{
				Path:   f.Field.Name,
				GoPath: f.Field.TitleCase,
				Field:  f,
			})
		}
	}
	return fields
}

// GetFieldDBPath returns a field by its path. The path can be a dot-separated string of field names, e.g. "user.address.phoneNumber".
// The second return value is the database name of the field. It will use the [sep] as a separator, e.g. `"user"->"address"->"phone_number"`.
// The third return value is [true] if the field is found.
func (c *Struct) GetFieldDBPath(path string, sep string) (*StructField, string, bool) {
	if !strings.Contains(path, ".") {
		s, ok := c.GetField(path)
		if ok {
			return s, s.Field.DBName, true
		}
		return s, "", false
	}
	parts := strings.Split(path, ".")
	return c.getFieldDBPath(parts, sep)
}

func (c *Struct) getFieldDBPath(parts []string, sep string) (*StructField, string, bool) {
	f, ok := c.GetField(parts[0])
	if !ok {
		return f, "", false
	}
	if len(parts) == 1 {
		return f, f.Field.DBName, true
	}
	if f.Field.Struct == nil {
		return nil, "", false
	}
	if f.Anonymous {
		return f.Field.Struct.getFieldDBPath(parts[1:], sep)
	}
	f2, db, ok := f.Field.Struct.getFieldDBJSONPath(parts[1:], sep)
	if !ok {
		return f2, "", false
	}
	return f2, f.Field.DBName + sep + db, true
}

func (c *Struct) getFieldDBJSONPath(parts []string, sep string) (*StructField, string, bool) {
	f, ok := c.GetField(parts[0])
	if !ok {
		return f, "", false
	}
	name := "'" + f.Field.Name + "'"
	if len(parts) == 1 {
		return f, name, true
	}
	if f.Field.Struct == nil {
		return nil, "", false
	}
	f2, db, ok := f.Field.Struct.getFieldDBJSONPath(parts[1:], sep)
	if !ok {
		return f2, "", false
	}
	return f2, name + sep + db, true
}

// GetDecoderField returns a field by its name. The field must be accessible in the given scope.
func (c *Struct) GetDecoderField(name string, context scopes.Context, scope scopes.Decoding) (*StructField, bool) {
	for _, f := range c.Fields {
		if f.Field.Name == name {
			return f, f.Field.CheckDecoderScope(context, scope)
		}
	}
	return nil, false
}

func newStruct(t reflect.Type) *Struct {
	return &Struct{
		Type: t,
	}
}

func (s *Struct) init() {
	n := s.Type.NumField()
	s.Fields = make([]*StructField, 0, n)
	for i := 0; i < n; i++ {
		f := s.Type.Field(i)
		s.initField(f)
	}
	s.byCamelName = make(map[string]*StructField, len(s.Fields))
	for _, f := range s.Fields {
		s.byCamelName[f.Field.Name] = f
	}
}

func (s *Struct) addField(f *StructField) {
	for i, ff := range s.Fields {
		if ff.Field.Name == f.Field.Name {
			s.Fields[i] = f
			return
		}
	}
	s.Fields = append(s.Fields, f)
}

func (c *Struct) initField(f reflect.StructField) {

	// Ignore unexported fields
	if !f.IsExported() {
		return
	}
	ft := mirror.DerefType(f.Type)

	res := &StructField{
		Anonymous: f.Anonymous,
		Embedded:  f.Anonymous,
		Field:     &Field{Type: ft, Kind: ft.Kind(), TitleCase: f.Name},
		Idx:       f.Index,
	}
	res.Field.ParseTag(f.Tag)

	if res.Field.Name == "" {
		res.Field.Name = stringer.ToCamelCase(f.Name)
	} else {
		res.Embedded = false
	}

	if res.Field.Kind == reflect.Struct && res.Field.Type != reflect.TypeFor[time.Time]() {
		if ft != c.Type {
			s := Cache.Get(ft)
			res.Field.Struct = s
			if res.Embedded {
				for _, f := range s.Fields {
					c.addField(&StructField{Field: f.Field, Anonymous: f.Anonymous, Embedded: true, Idx: append(res.Idx, f.Idx...)})
				}
				// Do not add the struct as a field if it's embedded
				return
			}
		} else {
			if res.Embedded {
				// Ignore self embedding
				return
			}
			res.Field.Struct = c
		}
	}

	res.Field.ObjectKey = []byte(`"` + res.Field.Name + `":`)
	res.Field.DBName = `"` + GetDBName(f, res) + `"`
	c.addField(res)
}
