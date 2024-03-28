package encoder

import (
	"reflect"

	"github.com/deveox/blaze/types"
)

func encodeStruct(e *Encoder, v reflect.Value, si *types.Struct) error {
	e.depth++
	defer func() {
		e.depth--
	}()
	anonymous := e.anonymous
	if !anonymous {
		e.WriteByte('{')
	}

	for _, fi := range si.Fields {
		ok := fi.CheckEncoderScope(e.config.Scope)
		if !ok {
			continue
		}
		f := v.Field(fi.Idx)

		// Handle zero values

		if f.IsZero() && !fi.KeepEmpty {
			continue
		}
		oldLen := len(e.bytes)
		switch f.Kind() {
		case reflect.Struct:
			if !fi.Anonymous {
				e.Write(fi.ObjectKey)
				oldLen = len(e.bytes)
			}
			e.anonymous = fi.Anonymous
			if err := e.encode(f); err != nil {
				return err
			}
		default:
			e.Write(fi.ObjectKey)
			oldLen = len(e.bytes)
			if err := e.encode(f); err != nil {
				return err
			}
		}
		if len(e.bytes) == oldLen {
			if fi.Anonymous && f.Kind() == reflect.Struct {
				continue
			}
			e.bytes = e.bytes[:len(e.bytes)-len(fi.ObjectKey)]
		} else {
			e.WriteByte(',')
		}
	}
	last := len(e.bytes) - 1
	if anonymous {
		if e.bytes[last] == ',' {
			e.bytes = e.bytes[:last]
		}
	} else {
		if e.bytes[last] == ',' {
			e.bytes[last] = '}'
		} else {
			e.bytes = append(e.bytes, '}')
		}
	}
	return nil
}

func newStructEncoder(t reflect.Type) EncoderFn {
	si := types.Cache.Get(t)
	return func(e *Encoder, v reflect.Value) error {
		return encodeStruct(e, v, si)
	}
}
