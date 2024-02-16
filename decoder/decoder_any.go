package decoder

import "reflect"

// TODO remove allocations
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
		var fl *float64
		vv := reflect.ValueOf(fl)
		err := decodeFloat(d, vv)
		if err != nil {
			return err
		}
		v.Set(vv)
		return nil
	case 't', 'f':
		var bol *bool
		vv := reflect.ValueOf(bol)
		err := decodeBool(d, vv)
		if err != nil {
			return err
		}
		v.Set(vv)
		return nil
	case 'n':
		err := d.ScanNull()
		if err != nil {
			return err
		}
		v.SetZero()
		return nil
	case '[':
		var an *[]any
		vv := reflect.ValueOf(an)
		err := d.decodeSlice(vv, decodeAny)
		if err != nil {
			return err
		}
		v.Set(vv)
		return nil
	case '{':
		var mp *map[string]any
		vv := reflect.ValueOf(mp)
		err := d.decodeMap(vv, decodeString, decodeAny)
		if err != nil {
			return err
		}
		v.Set(vv)
		return nil
	case TERMINATION_CHAR:
		return d.Error("[Blaze decode()] unexpected end of input, expected beginning of value")
	default:
		return d.Error("[Blaze decode()] invalid char, expected beginning of value")
	}
}
