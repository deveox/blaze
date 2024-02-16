package decoder

import (
	"reflect"
)

func (d *Decoder) SkipArray() error {
	level := 1
	d.pos++
	for {
		c := d.char()
		switch c {
		case '[':
			level++
		case ']':
			level--
			if level == 0 {
				d.pos++
				return nil
			}
		case '"':
			err := d.SkipString()
			if err != nil {
				return err
			}
			continue
		case TERMINATION_CHAR:
			return d.Error("[Blaze SkipArray()] unexpected end of input, expected ']'")
		}
		d.pos++
	}
}

func (d *Decoder) ScanArray() (int, error) {
	d.SkipWhitespace()

	start := d.pos
	size := 0
	d.pos++
	for {
		c := d.char()
		switch c {
		case ',':
			size++
			d.pos++
			err := d.Skip()
			if err != nil {
				return 0, err
			}
			continue
		case ']':
			d.pos++

			d.start = start
			return size, nil
		case TERMINATION_CHAR:
			return 0, d.Error("[Blaze ScanArray()] unexpected end of input, expected ']'")
		default:
			size++
			err := d.Skip()
			if err != nil {
				return 0, err
			}
			continue
		}
	}
}

func (d *Decoder) decodeArray(v reflect.Value, elemDecoder DecoderFn) error {

	d.SkipWhitespace()
	c := d.char()
	switch c {
	case '[':
		d.pos++
	case 'n':
		err := d.ScanNull()
		if err != nil {
			return err
		}
		v.SetZero()
		return nil
	default:
		return d.Error("[Blaze decodeArray()] invalid char, expected '[' or 'null'")
	}
	d.depth++
	if d.depth > MAX_DEPTH {
		return d.Error("[Blaze decodeArray()] max depth reached")
	}
	i := -1
	for {
		c := d.char()
		d.SkipWhitespace()
		switch c {
		case ',':
			i++
			d.pos++
			if i < v.Len() {
				d.SkipWhitespace()
				err := elemDecoder(d, v.Index(i))
				if err != nil {
					return err
				}
			} else {
				err := d.Skip()
				if err != nil {
					return err
				}
			}
		case ']':
			d.pos++
			d.depth--
			return nil
		case TERMINATION_CHAR:
			return d.Error("[Blaze decodeArray()] unexpected end of input, expected ']'")
		default:
			i++
			if i < v.Len() {
				err := elemDecoder(d, v.Index(i))
				if err != nil {
					return err
				}
			} else {
				err := d.Skip()
				if err != nil {
					return err
				}
			}
		}
	}
}

func newArrayDecoder(d reflect.Type) DecoderFn {
	elemDecoder := newDecoderFn(d.Elem(), true)
	return func(d *Decoder, v reflect.Value) error {
		return d.decodeArray(v, elemDecoder)
	}
}
