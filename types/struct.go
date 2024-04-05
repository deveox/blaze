package types

import (
	"fmt"
	"reflect"

	"github.com/deveox/blaze/scopes"

	"github.com/deveox/gu/mirror"
	"github.com/deveox/gu/stringer"
)

type Struct struct {
	Type   reflect.Type
	Fields []*StructField
}

func (c *Struct) GetField(name string) (*StructField, bool) {
	for _, f := range c.Fields {
		if f.Field.Name == name {
			return f, true
		}
	}
	return nil, false
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
		Field: &Field{Type: f.Type},
		Idx:   f.Index,
	}
	anonymous := f.Anonymous
	res.Field.ParseTag(f.Tag)

	if res.Field.Name == "" {
		res.Field.Name = stringer.ToCamelCase(f.Name)
	} else {
		anonymous = false
	}

	if ft.Kind() == reflect.Struct {
		if ft != c.Type {
			s := Cache.Get(f.Type)
			res.Field.Struct = s
			if anonymous {
				for _, f := range s.Fields {
					c.AddField(&StructField{Field: f.Field, Idx: append(res.Idx, f.Idx...)})
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
	res.Field.ObjectKey = []byte(fmt.Sprintf("\"%s\":", res.Field.Name))
	c.AddField(res)
}
