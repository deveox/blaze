package decoder

import "reflect"

func (d *Decoder) SkipTrue() error {
	d.start = d.pos - 1

	c := d.char()
	if c != 'r' {
		return d.ErrorF("[Blaze SkipTrue()] invalid char, expected 'true'")
	}
	d.pos++
	c = d.char()
	if c != 'u' {
		return d.ErrorF("[Blaze SkipTrue()] invalid char, expected 'true'")
	}
	d.pos++
	c = d.char()
	if c != 'e' {
		return d.ErrorF("[Blaze SkipTrue()] invalid char, expected 'true'")
	}
	d.pos++
	return nil
}

func (d *Decoder) SkipFalse() error {
	d.start = d.pos - 1

	c := d.char()
	if c != 'a' {
		return d.ErrorF("[Blaze SkipFalse()] invalid char, expected 'false'")
	}
	d.pos++
	c = d.char()
	if c != 'l' {
		return d.ErrorF("[Blaze SkipFalse()] invalid char, expected 'false'")
	}
	d.pos++
	c = d.char()
	if c != 's' {
		return d.ErrorF("[Blaze SkipFalse()] invalid char, expected 'false'")
	}
	d.pos++
	c = d.char()
	if c != 'e' {
		return d.ErrorF("[Blaze SkipFalse()] invalid char, expected 'false'")
	}
	d.pos++
	return nil
}

func decodeBool(d *Decoder, v reflect.Value) error {
	d.SkipWhitespace()
	c := d.char()
	switch c {
	case 'n':
		err := d.ScanNull()
		if err != nil {
			return err
		}
		v.SetBool(false)
		return nil
	case 't':
		d.pos++
		err := d.SkipTrue()
		if err != nil {
			return err
		}
		v.SetBool(true)
		return nil
	case 'f':
		d.pos++
		err := d.SkipFalse()
		if err != nil {
			return err
		}
		v.SetBool(false)
		return nil
	default:
		return d.ErrorF("[Blaze decodeBool()] invalid char, expected 't' or 'f'")
	}
}
