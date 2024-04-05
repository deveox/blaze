package types

import (
	"reflect"
	"strings"

	"github.com/deveox/blaze/scopes"

	"github.com/deveox/gu/mirror"
	"github.com/deveox/gu/stringer"
)

type Struct struct {
	Type        reflect.Type
	Fields      []*StructField
	ByCamelName map[string]*StructField
}

func (c *Struct) GetField(name string) (*StructField, bool) {
	for _, f := range c.Fields {
		if f.Field.Name == name {
			return f, true
		}
	}
	return nil, false
}

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
		if f.Field.Struct != nil {
			return f, "", false
		}
		return f, f.Field.DBName, true
	}
	if f.Field.Struct == nil {
		return nil, "", false
	}
	f2, db, ok := f.Field.Struct.getFieldDBPath(parts[1:], sep)
	if !ok {
		return f2, "", false
	}
	if f.Anonymous {
		return f2, db, true
	}
	return f2, f.Field.DBName + sep + db, true
}

func (c *Struct) GetDecoderField(name string, context scopes.Context, scope scopes.Decoding) (*StructField, bool) {
	for _, f := range c.Fields {
		if f.Field.Name == name {
			return f, f.Field.CheckDecoderScope(context, scope)
		}
	}
	return nil, false
}

func NewStruct(t reflect.Type) *Struct {
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
	s.ByCamelName = make(map[string]*StructField, len(s.Fields))
	for _, f := range s.Fields {
		s.ByCamelName[f.Field.Name] = f
	}
}

func (s *Struct) AddField(f *StructField) {
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
		Field:     &Field{Type: ft, Kind: ft.Kind(), TitleCase: f.Name},
		Idx:       f.Index,
	}
	anonymous := f.Anonymous
	res.Field.ParseTag(f.Tag)

	if res.Field.Name == "" {
		res.Field.Name = stringer.ToCamelCase(f.Name)
	} else {
		anonymous = false
	}

	if res.Field.Kind == reflect.Struct {
		if ft != c.Type {
			s := Cache.Get(ft)
			res.Field.Struct = s
			if anonymous {
				for _, f := range s.Fields {
					c.AddField(&StructField{Field: f.Field, Anonymous: f.Anonymous, Idx: append(res.Idx, f.Idx...)})
				}
				// Do not add the struct as a field if it's embedded
				return
			}
		} else {
			if anonymous {
				// Ignore self embedding
				return
			}
			res.Field.Struct = c
		}
	}

	res.Field.ObjectKey = []byte(`"` + res.Field.Name + `":`)
	res.Field.DBName = `"` + GetDBName(f, res) + `"`
	c.AddField(res)
}
