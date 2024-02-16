package decoder

import "reflect"

func decodeAddressableCustom(d *Decoder, v reflect.Value) error {
	return decodeCustom(d, v.Addr())
}

func decodePtrCustom(d *Decoder, v reflect.Value) error {
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	return decodeCustom(d, v)
}

func decodeCustom(d *Decoder, v reflect.Value) error {
	err := d.Skip()
	if err != nil {
		return err
	}
	u := v.Interface().(Unmarshaler)
	return u.UnmarshalBlaze(d, d.Buf[d.start:d.pos])
}
