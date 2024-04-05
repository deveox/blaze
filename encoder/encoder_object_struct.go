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
	path := e.fields.currentPath

	for _, fi := range si.Fields {
		ok := fi.Field.CheckEncoderScope(e.config.Scope)
		if !ok {
			continue
		}

		f := fi.Value(v)

		// Handle zero values
		if f.IsZero() && !fi.Field.KeepEmpty {
			continue
		}
		e.fields.currentPath = path
		if !e.fields.Has(fi.Field.Name, fi.Field.Short) {
			continue
		}
		e.Write(fi.Field.ObjectKey)
		oldLen := len(e.bytes)
		if err := e.encode(f); err != nil {
			return err
		}
		if len(e.bytes) == oldLen {
			e.bytes = e.bytes[:len(e.bytes)-len(fi.Field.ObjectKey)]
		} else {
			e.WriteByte(',')
		}
	}
	e.fields.currentPath = path
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
