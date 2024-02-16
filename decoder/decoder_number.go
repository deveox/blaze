package decoder

import (
	"reflect"
	"strconv"
)

func (d *Decoder) SkipNumber(withDecimal, withExponent bool) error {
	d.pos++
	for {
		c := d.char()
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			d.pos++
		case 'e':
			if withExponent {
				return d.SkipExponent()
			}
			return d.Error("[Blaze SkipAnyNumber()] invalid char, expected integer")
		case '.':
			if withDecimal {
				return d.SkipFloat()
			}
			return d.Error("[Blaze SkipAnyNumber()] invalid char, expected integer")
		default:
			return nil
		}
	}
}

func (d *Decoder) SkipFloat() error {
	d.pos++
	c := d.char()
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return d.SkipNumber(false, true)
	default:
		return d.Error("[Blaze SkipFloat()] invalid char, expected integer")
	}
}

func (d *Decoder) SkipZero(isFloat bool) error {
	d.pos++
	c := d.char()
	switch c {
	case '.':
		if isFloat {
			return d.SkipFloat()
		}
		return d.Error("[Blaze SkipZero()] invalid char, expected 'e' (exponent) or end of number")
	case 'e':
		if isFloat {
			return d.SkipExponent()
		}
		return d.Error("[Blaze SkipZero()] invalid char, expected end of number")
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		if isFloat {
			return d.Error("[Blaze SkipZero()] invalid char, expected '.' or 'e' (exponent)")
		}
		return d.Error("[Blaze SkipZero()] invalid char, expected end of number")
	default:
		return nil
	}
}

func (d *Decoder) SkipExponent() error {
	d.pos++
	c := d.char()
	switch c {
	case '+', '-':
		d.pos++
		c = d.char()
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			d.SkipNumber(false, false)
			return nil
		default:
			return d.Error("[Blaze SkipExponent()] invalid char, expected integer")
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		d.SkipNumber(false, false)
		return nil
	default:
		return d.Error("[Blaze SkipExponent()] invalid char, expected '+', '-' or integer")
	}
}

func (d *Decoder) SkipMinus(isFloat bool) error {
	d.start = d.pos
	d.pos++
	c := d.char()
	switch c {
	case '0':
		return d.SkipZero(isFloat)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return d.SkipNumber(isFloat, isFloat)
	default:
		return d.Error("[Blaze Skip()] invalid char, expected number")
	}
}

func decodeInt(d *Decoder, v reflect.Value) error {
	d.SkipWhitespace()
	c := d.char()
	d.start = d.pos
	switch c {
	case '"':
		d.pos++
		err := decodeInt(d, v)
		if err != nil {
			return err
		}
		d.pos++
		return nil
	case 'n':
		err := d.ScanNull()
		if err != nil {
			return err
		}
		v.SetInt(0)
		return nil
	case '-':
		err := d.SkipMinus(false)
		if err != nil {
			return err
		}
	case '0':
		err := d.SkipZero(false)
		if err != nil {
			return err
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		err := d.SkipNumber(false, false)
		if err != nil {
			return err
		}
	default:
		return d.Error("[Blaze decodeInt()] invalid char, expected '-' or integer")
	}

	str := BytesToString(d.Buf[d.start:d.pos])
	n, err := strconv.ParseInt(str, 10, v.Type().Bits())
	if err != nil {
		return d.Error(err.Error())
	}
	v.SetInt(n)
	return nil
}

func decodeUint(d *Decoder, v reflect.Value) error {
	d.SkipWhitespace()
	c := d.char()
	d.start = d.pos
	switch c {
	case '"':
		d.pos++
		err := decodeUint(d, v)
		if err != nil {
			return err
		}
		d.pos++
		return nil
	case 'n':
		err := d.ScanNull()
		if err != nil {
			return err
		}
		v.SetUint(0)
		return nil
	case '0':
		err := d.SkipZero(false)
		if err != nil {
			return err
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		err := d.SkipNumber(false, false)
		if err != nil {
			return err
		}
	default:
		return d.Error("[Blaze decodeInt()] invalid char, expected '-' or integer")
	}

	str := BytesToString(d.Buf[d.start:d.pos])
	n, err := strconv.ParseUint(str, 10, v.Type().Bits())
	if err != nil {
		return d.Error(err.Error())
	}
	v.SetUint(n)
	return nil
}

func decodeFloat(d *Decoder, v reflect.Value) error {
	d.SkipWhitespace()
	c := d.char()
	d.start = d.pos
	switch c {
	case '"':
		d.pos++
		err := decodeFloat(d, v)
		if err != nil {
			return err
		}
		d.pos++
		return nil
	case 'n':
		err := d.ScanNull()
		if err != nil {
			return err
		}
		v.SetFloat(0)
		return nil
	case '-':
		err := d.SkipMinus(true)
		if err != nil {
			return err
		}
	case '0':
		err := d.SkipZero(true)
		if err != nil {
			return err
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		err := d.SkipNumber(true, true)
		if err != nil {
			return err
		}
	default:
		return d.Error("[Blaze decodeInt()] invalid char, expected '-' or integer")
	}

	str := BytesToString(d.Buf[d.start:d.pos])
	n, err := strconv.ParseFloat(str, v.Type().Bits())
	if err != nil {
		return d.Error(err.Error())
	}
	v.SetFloat(n)
	return nil
}
