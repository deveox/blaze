package decoder

import (
	"reflect"
)

func (d *Decoder) ScanObject() (int, error) {
	start := d.pos
	size := 0
	d.pos++
	for {

		c := d.char()
		switch c {
		case ',':
			d.pos++
		case ':':
			d.pos++
			size++
		case '}':
			d.pos++
			d.start = start
			return size, nil
		case TERMINATION_CHAR:
			return 0, d.Error("[Blaze ScanObject()] unexpected end of input, expected '}'")
		default:
			err := d.Skip()
			if err != nil {
				return 0, err
			}
		}
	}
}

func (d *Decoder) SkipObject() error {
	level := 1
	for {
		d.pos++
		c := d.char()
		switch c {
		case '{':
			level++

		case '}':
			level--
			if level == 0 {
				d.pos++
				return nil
			}
		case '"':
			d.SkipString()
			d.pos--
		case TERMINATION_CHAR:
			return d.Error("[Blaze SkipObject()] unexpected end of input, expected '}'")
		}
	}
}

func (d *Decoder) decodeMap(v reflect.Value, keyDec, elemDec DecoderFn) error {
	d.SkipWhitespace()
	c := d.char()
	switch c {
	case 'n':
		err := d.ScanNull()
		if err != nil {
			return err
		}
		v.SetZero()
		return nil
	case '{':
	default:
		return d.Error("[Blaze decodeMap()] expected '{' or 'null'")
	}
	d.depth++
	if d.depth > MAX_DEPTH {
		return d.Error("[Blaze decodeMap()] max depth reached")
	}
	size, err := d.ScanObject()
	if err != nil {
		return err
	}
	if v.IsNil() {
		v.Set(reflect.MakeMapWithSize(v.Type(), size))
	}
	d.pos = d.start + 1
	for {
		d.SkipWhitespace()
		c := d.char()
		switch c {
		case '}':
			d.pos++
			d.depth--
			return nil
		case '"':
		case TERMINATION_CHAR:
			return d.Error("[Blaze Decode Map] unexpected end of input, expected object key")
		default:
			return d.Error("[Blaze Decode Map] expected object key")
		}

		key := reflect.New(v.Type().Key()).Elem()
		err := keyDec(d, key)
		if err != nil {
			return err
		}
		d.SkipWhitespace()

		c = d.char()
		if c != ':' {
			return d.Error("[Blaze decodeMap()] expected ':'")
		}
		d.pos++
		d.SkipWhitespace()
		value := reflect.New(v.Type().Elem()).Elem()
		if err := elemDec(d, value); err != nil {
			return err
		}
		v.SetMapIndex(key, value)
		d.SkipWhitespace()
		c = d.char()
		switch c {
		case '}':
			d.pos++
			d.depth--
			return nil
		case ',':
			d.pos++
		case TERMINATION_CHAR:
			return d.Error("[Blaze decodeMap()] unexpected end of input, expected ',' or '}'")
		default:
			return d.Error("[Blaze decodeMap()] expected ',' or '}'")
		}
	}
}

func newMapEncoder(t reflect.Type) DecoderFn {
	keyDec := newDecoderFn(t.Key(), true)
	elemDec := newDecoderFn(t.Elem(), true)
	return func(d *Decoder, v reflect.Value) error {
		return d.decodeMap(v, keyDec, elemDec)
	}
}
