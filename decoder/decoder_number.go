package decoder

import (
	"reflect"
	"strconv"
)

func (d *Decoder) SkipAnyNumber() error {
	for {
		c := d.char(d.ptr, d.pos)
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			d.pos++
		case 'e':
			return d.SkipExponent()
		case '.':
			return d.SkipFloat()
		default:
			return nil
		}
	}
}

func (d *Decoder) SkipFloat() error {
	d.pos++
	c := d.char(d.ptr, d.pos)
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		d.pos++
		return d.SkipInteger()
	default:
		return d.Error("[Blaze SkipFloat()] invalid char, expected integer")
	}
}

func (d *Decoder) SkipZero() error {
	d.pos++
	c := d.char(d.ptr, d.pos)
	switch c {
	case '.':
		return d.SkipFloat()
	case 'e':
		return d.SkipExponent()
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		return d.Error("[Blaze SkipZero()] invalid char, expected '.' or 'e' (exponent)")
	default:
		return nil
	}
}

func (d *Decoder) SkipExponent() error {
	d.pos++
	c := d.char(d.ptr, d.pos)
	switch c {
	case '+', '-':
		d.pos++
		c = d.char(d.ptr, d.pos)
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			d.pos++
			d.SkipZeroOrMoreDigit()
			return nil
		default:
			return d.Error("[Blaze SkipExponent()] invalid char, expected integer")
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		d.pos++
		d.SkipZeroOrMoreDigit()
		return nil
	default:
		return d.Error("[Blaze SkipExponent()] invalid char, expected '+', '-' or integer")
	}
}

func (d *Decoder) SkipInteger() error {
	for {
		c := d.char(d.ptr, d.pos)
		switch c {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			d.pos++
		case 'e':
			return d.SkipExponent()
		default:
			return nil
		}
	}
}

func (d *Decoder) SkipZeroOrMoreDigit() {
	for {
		c := d.char(d.ptr, d.pos)
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			d.pos++
		default:
			return
		}
	}
}

func (d *Decoder) DecodeFloat(bitSize int) (float64, error) {
	str := BytesToString(d.Buf[d.start:d.pos])
	return strconv.ParseFloat(str, bitSize)
}

func (d *Decoder) decodeFloat(v reflect.Value) error {
	n, err := d.DecodeFloat(v.Type().Bits())
	if err != nil {
		return d.Error(err.Error())
	}

	v.SetFloat(n)
	return nil
}

func (d *Decoder) DecodeInt(bitSize int) (int64, error) {
	str := BytesToString(d.Buf[d.start:d.pos])
	return strconv.ParseInt(str, 10, bitSize)
}

func (d *Decoder) decodeInt(v reflect.Value) error {
	n, err := d.DecodeInt(v.Type().Bits())
	if err != nil {
		return d.Error(err.Error())
	}

	v.SetInt(n)
	return nil
}

func (d *Decoder) DecodeUint(bitSize int) (uint64, error) {
	str := BytesToString(d.Buf[d.start:d.pos])
	return strconv.ParseUint(str, 10, bitSize)
}

func (d *Decoder) decodeUint(v reflect.Value) error {
	n, err := d.DecodeUint(v.Type().Bits())
	if err != nil {
		return d.Error(err.Error())
	}

	v.SetUint(n)
	return nil
}

func (t *Decoder) decodeNumber(v reflect.Value) error {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return t.decodeFloat(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return t.decodeInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return t.decodeUint(v)
	case reflect.Interface:
		if v.NumMethod() > 0 {
			return t.ErrorF("[Blaze decodeNumber()] unsupported custom interface %s", v.Kind().String())
		}

		n, err := t.DecodeFloat(64)
		if err != nil {
			return t.Error(err.Error())
		}
		v.Set(reflect.ValueOf(n))
	default:
		return t.ErrorF("[Blaze decodeNumber()] can't decode number into Go type '%s'", v.Type())
	}
	return nil
}
