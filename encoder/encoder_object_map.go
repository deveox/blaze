package encoder

import "reflect"

func (e *Encoder) EncodeMap(v reflect.Value, keyEnc, valueEnc EncoderFn) error {
	e.depth++
	defer func() {
		e.depth--
	}()

	e.WriteByte('{')
	key := v.Type().Key()
	iter := v.MapRange()
	for {
		next := iter.Next()
		if next {
			oldLen := len(e.bytes)
			switch key.Kind() {
			case reflect.String:
				if err := keyEnc(e, iter.Key()); err != nil {
					return err
				}
			default:
				e.WriteByte('"')
				if err := keyEnc(e, iter.Key()); err != nil {
					return err
				}
				e.WriteByte('"')
			}
			e.WriteByte(':')
			keyLen := len(e.bytes) - oldLen
			oldLen = len(e.bytes)
			if err := valueEnc(e, iter.Value()); err != nil {
				return err
			}
			if len(e.bytes) == oldLen {
				e.bytes = e.bytes[:len(e.bytes)-keyLen]
			} else {
				e.WriteByte(',')
			}
		} else {
			last := len(e.bytes) - 1
			switch e.bytes[last] {
			case '{':
				if e.keep || e.depth == 1 {
					e.WriteByte('}')
				} else {
					e.bytes = e.bytes[:last]
				}
			case ',':
				e.bytes[last] = '}'
			default:
				e.WriteByte('}')
			}
			break
		}
	}
	return nil
}

func newMapEncoder(t reflect.Type) EncoderFn {
	keyEnc := newEncoderFn(t.Key(), false)
	valueEnc := newEncoderFn(t.Elem(), false)
	return func(e *Encoder, v reflect.Value) error {
		if v.IsNil() {
			e.WriteString("null")
			return nil
		}
		if v.IsZero() {
			e.bytes = append(e.bytes, '{', '}')
			return nil
		}
		return e.EncodeMap(v, keyEnc, valueEnc)
	}
}
