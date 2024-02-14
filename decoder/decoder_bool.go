package decoder

import "reflect"

func (t *Decoder) SkipTrue() error {
	t.start = t.pos - 1

	c := t.char(t.ptr, t.pos)
	if c != 'r' {
		return t.ErrorF("[Blaze SkipTrue()] invalid char, expected 'true'")
	}
	t.pos++
	c = t.char(t.ptr, t.pos)
	if c != 'u' {
		return t.ErrorF("[Blaze SkipTrue()] invalid char, expected 'true'")
	}
	t.pos++
	c = t.char(t.ptr, t.pos)
	if c != 'e' {
		return t.ErrorF("[Blaze SkipTrue()] invalid char, expected 'true'")
	}
	t.pos++
	return nil
}

func (t *Decoder) SkipFalse() error {
	t.start = t.pos - 1

	c := t.char(t.ptr, t.pos)
	if c != 'a' {
		return t.ErrorF("[Blaze SkipFalse()] invalid char, expected 'false'")
	}
	t.pos++
	c = t.char(t.ptr, t.pos)
	if c != 'l' {
		return t.ErrorF("[Blaze SkipFalse()] invalid char, expected 'false'")
	}
	t.pos++
	c = t.char(t.ptr, t.pos)
	if c != 's' {
		return t.ErrorF("[Blaze SkipFalse()] invalid char, expected 'false'")
	}
	t.pos++
	c = t.char(t.ptr, t.pos)
	if c != 'e' {
		return t.ErrorF("[Blaze SkipFalse()] invalid char, expected 'false'")
	}
	t.pos++
	return nil
}

func (t *Decoder) DecodeBool() bool {
	c := t.char(t.ptr, t.start)
	return c == 't'
}

func (t *Decoder) decodeBool(v reflect.Value) error {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Bool:
		v.SetBool(t.DecodeBool())
	case reflect.Interface:
		if v.NumMethod() > 0 {
			return t.ErrorF("[Blaze decodeBool()] unsupported custom interface %s", v.Kind().String())
		}
		v.Set(reflect.ValueOf(t.DecodeBool()))
	default:
		return t.ErrorF("[Blaze decodeBool()] can't decode bool into Go type '%s'", v.Type())
	}
	return nil
}
