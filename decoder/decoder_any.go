package decoder

import "reflect"

func decodeAny(d *Decoder, v reflect.Value) error {
	d.SkipWhitespace()
	c := d.char()
	switch c {
	case '"':
		str := reflect.ValueOf("")
		err := decodeString(d, str)
		if err != nil {
			return err
		}
		v.Set(str)
		return nil
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-':
		fl := reflect.ValueOf(float64(0))
		err := decodeFloat(d, fl)
		if err != nil {
			return err
		}
		v.Set(fl)
		return nil
	case 't', 'f':
		bl := reflect.ValueOf(false)
		err := decodeBool(d, bl)
		if err != nil {
			return err
		}
		v.Set(bl)
		return nil
	case 'n':
		err := d.ScanNull()
		if err != nil {
			return err
		}
		v.SetZero()
		return nil
	case '[':
		sl := reflect.ValueOf([]any{})
		err := d.decodeSlice(sl, decodeAny)
		if err != nil {
			return err
		}
		v.Set(sl)
		return nil
	case '{':
		mp := reflect.ValueOf(map[string]any{})
		err := d.decodeMap(mp, decodeString, decodeAny)
		if err != nil {
			return err
		}
		v.Set(mp)
		return nil
	case TERMINATION_CHAR:
		return d.Error("[Blaze decode()] unexpected end of input, expected beginning of value")
	default:
		return d.Error("[Blaze decode()] invalid char, expected beginning of value")
	}
}
