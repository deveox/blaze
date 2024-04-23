package encoder

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/deveox/blaze/ctx"
	"github.com/deveox/blaze/scopes"
)

const MAX_DEPTH = 10000

type Encoder struct {
	*ctx.Ctx
	config    *Config
	bytes     []byte
	depth     int
	fields    *fields
	anonymous bool
}

// GetCurrentPath will return the path of the current field being encoded if encoder is created by MarshalPartial
// Otherwise it will return an empty string
func (e *Encoder) GetCurrentPath() string {
	return e.fields.currentPath
}

// GetFields will return the fields that should be encoded if encoder is created by MarshalPartial
// Otherwise it will return an empty slice
func (e *Encoder) GetFields() []string {
	return e.fields.fields
}

func (e *Encoder) Context() scopes.Context {
	return e.config.Scope
}

func (e *Encoder) Marshal(v any) ([]byte, error) {
	return e.config.MarshalCtx(v, e.Ctx)
}

func (e *Encoder) MarshalPartial(v any, fields []string, short bool) ([]byte, error) {
	return e.config.MarshalPartialCtx(v, fields, short, e.Ctx)
}

func (e *Encoder) Encode(v any) error {
	return e.encode(reflect.ValueOf(v))
}

// EncodePartial works like MarshalPartial but it will respect partial settings of the encoder
// If fields were set by MarshalPartial, it will preserve them.
// If you set short parameter to 1 it will use short encoding for the fields, 0 will not use short encoding, -1 will preserve the current setting
// If you pass fields, their names will be prepended with the current path of the encoder
func (e *Encoder) EncodePartial(v any, fields []string, short int) error {
	oldFields := *e.fields
	if fields != nil || short == 1 {
		e.fields.enabled = true
	}
	if e.fields.currentPath != "" {
		for _, f := range fields {
			e.fields.fields = append(e.fields.fields, e.fields.currentPath+"."+f)
		}
	} else {
		e.fields.fields = append(e.fields.fields, fields...)
	}
	if short == 1 {
		e.fields.short = true
	} else if short == 0 {
		e.fields.short = false
	}

	err := e.encode(reflect.ValueOf(v))
	if err != nil {
		return err
	}
	*e.fields = oldFields
	return nil
}

func (e *Encoder) marshal(v any) ([]byte, error) {
	err := e.encode(reflect.ValueOf(v))
	if err != nil {
		return nil, err
	}
	res := e.bytes[:len(e.bytes)]

	e.bytes = e.bytes[:0]
	return res, nil
}

func (e *Encoder) encode(v reflect.Value) error {
	if e.depth > MAX_DEPTH {
		return e.ErrorF("exceeded max depth of %d", MAX_DEPTH)
	}
	if !v.IsValid() {
		e.WriteString("null")
		return nil
	}
	err := getEncoderFn(v.Type())(e, v)
	e.anonymous = false
	return err
}

func (e *Encoder) Error(msg string) error {
	return errors.New(msg)
}

func (e *Encoder) ErrorF(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func AddIndent(b []byte) []byte {
	res := make([]byte, 0, len(b))
	depth := 0
	for i := 0; i < len(b); i++ {
		if b[i] == '{' || b[i] == '[' {
			depth++
			res = append(res, b[i], '\n')
			for j := 0; j < depth-1; j++ {
				res = append(res, '\t')
			}

		} else if b[i] == '}' || b[i] == ']' {
			depth--
			res = append(res, '\n')
			for j := 0; j < depth-1; j++ {
				res = append(res, '\t')
			}
			res = append(res, b[i])
		} else if b[i] == ',' {
			res = append(res, b[i], '\n')
			for j := 0; j < depth-1; j++ {
				res = append(res, '\t')
			}
		} else {
			res = append(res, b[i])
		}
	}
	return res
}

func encodeInterface(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		e.WriteString("null")
		return nil
	}
	return e.encode(v.Elem())
}

func encodePtr(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		e.WriteString("null")
		return nil
	}
	return e.encode(v.Elem())
}

func encodeUnsupported(e *Encoder, v reflect.Value) error {
	return e.ErrorF("[blaze encodeUnsupported()] unsupported type: %s", v.Type())
}
