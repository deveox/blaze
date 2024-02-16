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

	if !e.anonymous {
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

		switch f.Kind() {
		case reflect.Struct:
			if !fi.Anonymous {
				e.Write(fi.ObjectKey)
			}
			e.anonymous = fi.Anonymous
			if err := e.encode(f); err != nil {
				return err
			}
		default:
			e.Write(fi.ObjectKey)
			if err := e.encode(f); err != nil {
				return err
			}
		}
		e.WriteByte(',')
	}
	last := len(e.bytes) - 1
	if e.anonymous {
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
