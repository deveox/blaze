package decoder

import (
	"reflect"
)

func (t *Decoder) SkipNull() error {
	c := t.char(t.ptr, t.pos)
	if c != 'n' {
		return t.ErrorF("[Blaze SkipNull()] invalid char, expected 'null'")
	}
	t.pos++
	c = t.char(t.ptr, t.pos)
	if c != 'u' {
		return t.ErrorF("[Blaze SkipNull()] invalid char, expected 'null'")
	}
	t.pos++
	c = t.char(t.ptr, t.pos)
	if c != 'l' {
		return t.ErrorF("[Blaze SkipNull()] invalid char, expected 'null'")
	}
	t.pos++
	c = t.char(t.ptr, t.pos)
	if c != 'l' {
		return t.ErrorF("[Blaze SkipNull()] invalid char, expected 'null'")
	}
	return nil
}

func (t *Decoder) decodeNull(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		if v.CanAddr() {
			v.SetZero()
		} else {
			for !v.CanAddr() {
				v = v.Elem()
			}
			v.SetZero()
		}
		return nil
	default:
		if !v.IsZero() {
			v.SetZero()
		}
	}
	return nil
}
