package decoder

import (
	"encoding/base64"
	"reflect"
)

func (d *Decoder) decodeSlice(v reflect.Value, elemDecoder DecoderFn) error {
	d.SkipWhitespace()
	c := d.char()
	switch c {
	case '[':
	case 'n':
		err := d.ScanNull()
		if err != nil {
			return err
		}
		v.SetZero()
		return nil
	default:
		return d.Error("[Blaze decodeSlice()] invalid char, expected '[' or 'null'")
	}
	d.depth++
	if d.depth > MAX_DEPTH {
		return d.Error("[Blaze decodeSlice()] maximum depth reached")
	}
	size, err := d.ScanArray()
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
	d.pos = d.start + 1
	for {
		d.SkipWhitespace()
		c := d.char()
		switch c {
		case ',':
			i++
			d.pos++
			d.SkipWhitespace()
			err := elemDecoder(d, v.Index(i))
			if err != nil {
				return err
			}
		case ']':
			d.pos++
			d.depth--
			return nil
		case TERMINATION_CHAR:
			return d.Error("[Blaze decodeSlice()] unexpected end of input, expected ']'")
		default:
			i++
			err := elemDecoder(d, v.Index(i))
			if err != nil {
				return err
			}
		}
	}
}

func newSliceDecoder(t reflect.Type) DecoderFn {
	if t.Elem() == reflect.TypeFor[byte]() {
		return decodeBytes
	}
	elemDecoder := newDecoderFn(t.Elem(), true)
	return func(d *Decoder, v reflect.Value) error {
		return d.decodeSlice(v, elemDecoder)
	}
}

func decodeBytes(d *Decoder, v reflect.Value) error {
	d.SkipWhitespace()
	d.start = d.pos
	if d.char() != '"' {
		return d.Error("[Blaze decodeBytes()] expected '\"'")
	}
	d.pos++
	b, err := d.unquoteString()
	if err != nil {
		return err
	}
	res := make([]byte, base64.StdEncoding.DecodedLen(len(b)))
	n, err := base64.StdEncoding.Decode(res, b)
	if err != nil {
		return err
	}
	v.SetBytes(res[:n])
	return nil
}
