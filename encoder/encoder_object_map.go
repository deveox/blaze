package encoder

import "reflect"

func (e *Encoder) EncodeMap(v reflect.Value, keyEnc, valueEnc EncoderFn) error {
	e.WriteByte('{')
	key := v.Type().Key()
	iter := v.MapRange()
	for {
		next := iter.Next()
		if next {
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
			if err := valueEnc(e, iter.Value()); err != nil {
				return err
			}
			e.WriteByte(',')
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
