package types

import (
	"fmt"
	"reflect"

	"github.com/deveox/blaze/scopes"

	"github.com/deveox/gu/mirror"
	"github.com/deveox/gu/stringer"
)

type EmbeddedField struct {
	Field *Field
	Idx   []int
}

func (e *EmbeddedField) Value(v reflect.Value) reflect.Value {
	for _, i := range e.Idx {
		v = v.Field(i)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		}
	}
	return v
}

type Struct struct {
	Type           reflect.Type
	Fields         []*Field
	EmbeddedFields []*EmbeddedField
}

func (c *Struct) GetField(name string) (*Field, *EmbeddedField, bool) {
	for _, f := range c.Fields {
		if f.Name == name {
			return f, nil, true
		}
	}
	for _, f := range c.EmbeddedFields {
		if f.Field.Name == name {
			return f.Field, f, true
		}
	}
	return nil, nil, false
}

func (c *Struct) GetDecoderField(name string, context scopes.Context, scope scopes.Decoding) (*Field, *EmbeddedField, bool) {
	for _, f := range c.Fields {
		if f.Name == name {

			return f, nil, f.CheckDecoderScope(context, scope)
		}
	}

	for _, f := range c.EmbeddedFields {
		if f.Field.Name == name {
			return f.Field, f, f.Field.CheckDecoderScope(context, scope)
		}
	}
	return nil, nil, false
}

func NewStruct(t reflect.Type) *Struct {
	return &Struct{
		Type: t,
	}
}

func (s *Struct) init() error {
	n := s.Type.NumField()
	s.Fields = make([]*Field, 0, n)
	for i := 0; i < n; i++ {
		f := s.Type.Field(i)
		err := s.initField(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Struct) initField(f reflect.StructField) error {

	// Ignore unexported fields
	if !f.IsExported() {
		return nil
	}
	ft := mirror.DerefType(f.Type)

	res := &Field{
		Type:      f.Type,
		Idx:       f.Index[0],
		Anonymous: f.Anonymous,
	}
	res.ParseTag(f.Tag)

	if res.Name == "" {
		res.Name = stringer.ToCamelCase(f.Name)
	}

	if ft.Kind() == reflect.Struct {
		if ft != c.Type {
			s, err := Cache.Get(f.Type)
			if err != nil {
				return err
			}
			res.Struct = s
			if res.Anonymous {
				for _, f := range s.Fields {
					c.EmbeddedFields = append(c.EmbeddedFields, &EmbeddedField{Field: f, Idx: []int{res.Idx, f.Idx}})
				}
				for _, f := range s.EmbeddedFields {
					if len(f.Idx) >= 10 {
						return fmt.Errorf("embedded field %s depth is too deep: %d, maximum is 10", res.Name, len(f.Idx))
					}
					c.EmbeddedFields = append(c.EmbeddedFields, &EmbeddedField{Field: f.Field, Idx: append([]int{res.Idx}, f.Idx...)})
				}
			}
		} else {
			res.Struct = c
		}

	}
	res.ObjectKey = []byte(fmt.Sprintf("\"%s\":", res.Name))

	c.Fields = append(c.Fields, res)
	return nil
}
