package decoder

import "reflect"

func (d *Decoder) SkipTrue() error {
	d.pos++

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
	d.pos++
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

func (d *Decoder) decodeToBool() (bool, error) {
	d.SkipWhitespace()
	c := d.char()
	switch c {
	case 'n':
		err := d.ScanNull()
		return false, err
	case 't':
		err := d.SkipTrue()
		return true, err
	case 'f':
		err := d.SkipFalse()
		return false, err
	default:
		return false, d.ErrorF("[Blaze decodeBool()] invalid char %s, expected 't' or 'f'", string(c))
	}
}

func decodeBool(d *Decoder, v reflect.Value) error {
	b, err := d.decodeToBool()
	if err != nil {
		return err
	}
	v.SetBool(b)
	return nil
}
