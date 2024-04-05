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
		err := encodeStructField(e, f, fi, f.Kind())
		if err != nil {
			return err
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

func encodeStructField(e *Encoder, v reflect.Value, fi *types.StructField, kind reflect.Kind) error {
	enabled := e.fields.enabled
	if enabled {
		switch kind {
		case reflect.Ptr, reflect.Interface:
			f := v
			for f.Kind() == reflect.Ptr || f.Kind() == reflect.Interface {
				if f.IsNil() {
					f.Set(reflect.New(f.Type().Elem()))
				}
				f = f.Elem()
			}
			return encodeStructField(e, v, fi, f.Kind())
		case reflect.Struct:
			// Encode full struct if its field name specified
			if e.fields.Has(fi.Field.Name, fi.Field.Short) {
				if !fi.Field.Short {
					e.fields.enabled = false
				}
			}
			// Otherwise, continue to check its fields
		default:
			// Skip if field name is not in the partial list
			if !e.fields.Has(fi.Field.Name, fi.Field.Short) {
				return nil
			}
		}
	}
	e.Write(fi.Field.ObjectKey)
	oldLen := len(e.bytes)
	if err := e.encode(v); err != nil {
		return err
	}
	if len(e.bytes) == oldLen {
		e.bytes = e.bytes[:len(e.bytes)-len(fi.Field.ObjectKey)]
	} else {
		e.WriteByte(',')
	}
	if enabled != e.fields.enabled {
		e.fields.enabled = enabled
	}
	return nil
}

func newStructEncoder(t reflect.Type) EncoderFn {
	si := types.Cache.Get(t)
	return func(e *Encoder, v reflect.Value) error {
		return encodeStruct(e, v, si)
	}
}
