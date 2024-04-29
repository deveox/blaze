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
	keep := e.keep
	for _, fi := range si.Fields {
		ok := fi.Field.CheckEncoderScope(e.config.Scope)
		if !ok {
			continue
		}

		f := fi.Value(v)

		// Handle zero values
		if f.IsZero() {
			if fi.Field.KeepEmpty {
				e.keep = true
			} else {
				continue
			}
		}
		e.fields.currentPath = path
		err := encodeStructField(e, f, fi, f.Kind())
		if err != nil {
			return err
		}
		e.keep = false
	}
	e.fields.currentPath = path
	last := len(e.bytes) - 1
	if anonymous {
		if e.bytes[last] == ',' {
			e.bytes = e.bytes[:last]
		}
	} else {
		switch e.bytes[last] {
		case '{':
			if keep || e.depth == 1 {
				e.WriteByte('}')
			} else {
				e.bytes = e.bytes[:last]
			}
		case ',':
			e.bytes[last] = '}'
		default:
			e.WriteByte('}')
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
		case reflect.Array, reflect.Slice, reflect.Map:
			if e.fields.Has(fi.Field.Name, fi.Field.Short) {
				if !fi.Field.Short {
					e.fields.enabled = false
				}
			}
		case reflect.Struct:
			// Encode full struct if its field name specified
			if e.fields.Has(fi.Field.Name, fi.Field.Short) {
				if !fi.Field.Short && fi.Field.Struct != nil {
					e.fields.enabled = false
				}
			} else {
				if fi.Field.Struct == nil {
					return nil
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
	if fi.Field.StringEncoding {
		if err := encodeString(e, v); err != nil {
			return err
		}
	} else {
		if err := e.encode(v); err != nil {
			return err
		}
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
