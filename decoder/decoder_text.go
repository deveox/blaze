package decoder

import (
	"encoding"
	"reflect"
)

func decodeAddressableText(d *Decoder, v reflect.Value) error {
	return decodeText(d, v.Addr())
}

func decodePtrText(d *Decoder, v reflect.Value) error {
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	return decodeText(d, v)
}

// TODO: unquote bytes
func decodeText(d *Decoder, v reflect.Value) error {
	err := d.Skip()
	if err != nil {
		return err
	}
	u := v.Interface().(encoding.TextUnmarshaler)
	return u.UnmarshalText(d.Buf[d.start:d.pos])
}
