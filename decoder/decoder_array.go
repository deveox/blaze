package decoder

import (
	"reflect"
)

func (t *Decoder) SkipArray() error {
	level := 1
	t.pos++
	for {
		c := t.char(t.ptr, t.pos)
		switch c {
		case '[':
			t.depth++
			level++
		case ']':
			level--
			t.depth--
			if level == 0 {
				t.pos++
				return nil
			}
		case '"':
			err := t.SkipString()
			if err != nil {
				return err
			}
			continue
		case TERMINATION_CHAR:
			return t.Error("[Blaze SkipArray()] unexpected end of input, expected ']'")
		}
		t.pos++
	}
}

func (t *Decoder) ScanArray() (int, error) {
	start := t.pos
	size := 0
	t.pos++
	for {
		c := t.char(t.ptr, t.pos)
		switch c {
		case ',':
			size++
			t.pos++
			err := t.Skip()
			if err != nil {
				return 0, err
			}
			continue
		case ']':
			t.pos++

			t.start = start
			return size, nil
		case TERMINATION_CHAR:
			return 0, t.Error("[Blaze ScanArray()] unexpected end of input, expected ']'")
		default:
			size++
			err := t.Skip()
			if err != nil {
				return 0, err
			}
			continue
		}
	}
}

func (t *Decoder) decodeArrayOrSlice(v reflect.Value) error {
	t.start = t.pos
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			et := v.Type().Elem()
			v.Set(reflect.New(et))
			// if et.Kind() == reflect.Slice {
			// 	v.Set(reflect.MakeSlice(et))
			// } else {
			// }
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Slice:
		return t.decodeSlice(v)
	case reflect.Array:
		return t.decodeArray(v)
	case reflect.Interface:
		if v.NumMethod() > 0 {
			return t.ErrorF("[Blaze decodeArrayOrSlice()] unsupported custom interface %s", v.Type())
		}
		slice, err := t.DecodeArray()
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(slice))
		return nil
	default:
		return t.ErrorF("[Blaze decodeArrayOrSlice()] can't decode array into Go type '%s'", v.Type())
	}
}

func (t *Decoder) DecodeArray() ([]any, error) {
	var sv []any
	slice := reflect.ValueOf(sv)
	err := t.decodeSlice(slice)
	return slice.Interface().([]any), err
}

func (t *Decoder) decodeSlice(v reflect.Value) error {
	size, err := t.ScanArray()
	if err != nil {
		return err
	}
	if size == 0 {
		return nil

	}
	cap := v.Cap()
	if cap < size {
		v.Grow(size - cap)
	}
	v.SetLen(size)
	i := -1
	t.pos = t.start + 1
	for {
		c := t.char(t.ptr, t.pos)
		switch c {
		case ' ', '\t', '\n', '\r':
			t.pos++
		case ',':
			i++
			t.pos++
			t.SkipWhitespace()
			err := t.decode(v.Index(i))
			if err != nil {
				return err
			}
		case ']':
			t.pos++
			return nil
		case TERMINATION_CHAR:
			return t.Error("[Blaze decodeSlice()] unexpected end of input, expected ']'")
		default:
			i++
			err := t.decode(v.Index(i))
			if err != nil {
				return err
			}
		}
	}
}

func (t *Decoder) decodeArray(v reflect.Value) error {
	t.pos = t.start + 1
	i := -1
	for {
		c := t.char(t.ptr, t.pos)
		switch c {
		case ' ', '\t', '\n', '\r':
			t.pos++
		case ',':
			i++
			t.pos++
			if i < v.Len() {
				t.SkipWhitespace()
				err := t.decode(v.Index(i))
				if err != nil {
					return err
				}
			} else {
				err := t.Skip()
				if err != nil {
					return err
				}
			}
		case ']':
			t.pos++
			return nil
		case TERMINATION_CHAR:
			return t.Error("[Blaze decodeArray()] unexpected end of input, expected ']'")
		default:
			i++
			if i < v.Len() {
				err := t.decode(v.Index(i))
				if err != nil {
					return err
				}
			} else {
				err := t.Skip()
				if err != nil {
					return err
				}
			}
		}
	}
}
