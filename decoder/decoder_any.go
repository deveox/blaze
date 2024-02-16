package decoder

import "reflect"

// TODO remove allocations
// TODO add tests
func decodeAny(d *Decoder, v reflect.Value) error {
	d.SkipWhitespace()
	c := d.char()
	switch c {
	case '"':
		str, err := d.decodeToString()
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(str))
		return nil
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-':
		fl, err := d.decodeToFloat(64)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(fl))
		return nil
	case 't', 'f':
		bol, err := d.decodeToBool()
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(bol))
		return nil
	case 'n':
		err := d.ScanNull()
		if err != nil {
			return err
		}
		v.SetZero()
		return nil
	case '[':
		an := make([]any, 0)
		vv := reflect.ValueOf(&an)
		err := d.decodeSlice(vv.Elem(), decodeAny)
		if err != nil {
			return err
		}
		v.Set(vv.Elem())
		return nil
	case '{':
		mp := make(map[string]any)
		vv := reflect.ValueOf(&mp)
		err := d.decodeMap(vv.Elem(), decodeString, decodeAny)
		if err != nil {
			return err
		}
		v.Set(vv.Elem())
		return nil
	case TERMINATION_CHAR:
		return d.Error("[Blaze decode()] unexpected end of input, expected beginning of value")
	default:
		return d.Error("[Blaze decode()] invalid char, expected beginning of value")
	}
}
