package encoder

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/deveox/blaze/scopes"
	"github.com/deveox/blaze/types"
)

func Marshal(v any) ([]byte, error) {
	e := NewEncoder()
	defer encodersPool.Put(e)
	return e.Encode(v)
}
func MarshalScoped(v any, context scopes.Context) ([]byte, error) {
	e := NewEncoder()
	defer encodersPool.Put(e)
	e.ContextScope = context
	return e.Encode(v)
}

func NewEncoder() *Encoder {
	if v := encodersPool.Get(); v != nil {
		e := v.(*Encoder)
		e.bytes = e.bytes[:0]
		e.ContextScope = 0
		return e
	}
	return &Encoder{bytes: make([]byte, 0, 2048)}
}

var encodersPool sync.Pool

const MAX_DEPTH = 1000

type Encoder struct {
	bytes.Buffer
	ContextScope scopes.Context
	depth        int
	anonymous    bool
}

func (e *Encoder) Encode(v any) ([]byte, error) {
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
	enc, err := EncoderFns.Get(v)
	if err != nil {
		return err
	}
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.Type().Name() != "" {
			break
		}
		if v.IsNil() {
			e.bytes = append(e.bytes, "null"...)
			return nil
		}
		v = v.Elem()

	}
	err = enc(e, v)
	if err != nil {
		return err
	}
	e.anonymous = false
	return nil
}

func (e *Encoder) nativeEncoder(v reflect.Value) error {
	switch v.Kind() {

	case reflect.Struct:

		return e.EncodeStruct(v)
	case reflect.Slice, reflect.Array:
		if v.IsNil() {
			e.bytes = append(e.bytes, "null"...)
			return nil
		}
		if v.IsZero() {
			e.bytes = append(e.bytes, '[', ']')
			return nil
		}
		return e.EncodeSlice(v)
	case reflect.Interface, reflect.Ptr:
		for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
			if v.IsNil() {
				e.bytes = append(e.bytes, "null"...)
				return nil
			} else {
				v = v.Elem()
			}
		}
		return e.nativeEncoder(v)
	case reflect.Map:
		if v.IsNil() {
			e.bytes = append(e.bytes, "null"...)
			return nil
		}
		if v.IsZero() {
			e.bytes = append(e.bytes, '{', '}')
			return nil
		}
		return e.EncodeMap(v)
	// case reflect.Uintptr:

	// case reflect.UnsafePointer:

	// case reflect.Complex64, reflect.Complex128:

	// case reflect.Invalid:

	default:
		return fmt.Errorf("unsupported type: %v", v.Kind())
	}
	return nil
}

func (e *Encoder) EncodeStruct(v reflect.Value) error {
	e.depth++
	defer func() {
		e.depth--
	}()

	if !e.anonymous {
		e.bytes = append(e.bytes, '{')
	}

	t := v.Type()
	si, _ := types.Cache.Get(t)

	for _, fi := range si.Fields {
		ok := fi.CheckEncoderScope(e.ContextScope)
		if !ok {
			continue
		}
		f := v.Field(fi.Idx)

		// Handle zero values

		if f.IsZero() && !fi.KeepEmpty {
			continue
		}

		switch f.Kind() {
		case reflect.Struct:
			if !fi.Anonymous {
				e.bytes = append(e.bytes, fi.ObjectKey...)
			}
			e.anonymous = fi.Anonymous
			if err := e.encode(f); err != nil {
				return err
			}
		default:
			e.bytes = append(e.bytes, fi.ObjectKey...)
			if err := e.encode(f); err != nil {
				return err
			}
		}

		e.bytes = append(e.bytes, ',')
	}
	if e.anonymous {
		if e.bytes[len(e.bytes)-1] == ',' {
			e.bytes = e.bytes[:len(e.bytes)-1]
		}
	} else {
		if e.bytes[len(e.bytes)-1] == ',' {
			e.bytes[len(e.bytes)-1] = '}'
		} else {
			e.bytes = append(e.bytes, '}')
		}
	}
	return nil
}

func (e *Encoder) EncodeSlice(v reflect.Value) error {
	e.bytes = append(e.bytes, '[')
	n := v.Len()
	elem := v.Type().Elem()
	enc, err := EncoderFns.GetType(elem)
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		f := v.Index(i)
		for f.Kind() == reflect.Interface || f.Kind() == reflect.Ptr {
			if f.Type().Name() != "" {
				break
			}
			if f.IsNil() {
				e.bytes = append(e.bytes, "null"...)
				return nil
			}
			f = f.Elem()

		}
		err = enc(e, f)
		if err != nil {
			return err
		}
		e.bytes = append(e.bytes, ',')
	}
	if e.bytes[len(e.bytes)-1] == ',' {
		e.bytes[len(e.bytes)-1] = ']'
	} else {
		e.bytes = append(e.bytes, ']')
	}

	return nil
}

func (e *Encoder) EncodeMap(v reflect.Value) error {
	e.bytes = append(e.bytes, '{')
	key := v.Type().Key()
	iter := v.MapRange()
	for {
		next := iter.Next()
		if next {
			switch key.Kind() {
			case reflect.String:
				if err := e.encode(iter.Key()); err != nil {
					return err
				}
			default:
				e.bytes = append(e.bytes, '"')
				if err := e.encode(iter.Key()); err != nil {
					return err
				}
				e.bytes = append(e.bytes, '"')
			}
			e.bytes = append(e.bytes, ':')
			if err := e.encode(iter.Value()); err != nil {
				return err
			}
			e.bytes = append(e.bytes, ',')
		} else {
			if e.bytes[len(e.bytes)-1] == ',' {
				e.bytes[len(e.bytes)-1] = '}'
			} else {
				e.bytes = append(e.bytes, '}')
			}
			break
		}
	}
	return nil
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
