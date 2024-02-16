package encoder

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/deveox/blaze/scopes"
)

const MAX_DEPTH = 10000

type Encoder struct {
	config    *Config
	bytes     []byte
	depth     int
	anonymous bool
}

func (e *Encoder) Context() scopes.Context {
	return e.config.Scope
}

func (e *Encoder) Marshal(v any) ([]byte, error) {
	return e.config.Marshal(v)
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

		} else if b[i] == '}' || b[i] == ']' {
			res = append(res, '\n')

			res = append(res, b[i])
		} else if b[i] == ',' {
			res = append(res, b[i], '\n')
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
