package decoder

import (
	"encoding/json"
	"reflect"
)

func decodeAddressableStd(d *Decoder, v reflect.Value) error {
	return decodeStd(d, v.Addr())
}

func decodePtrStd(d *Decoder, v reflect.Value) error {
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	return decodeStd(d, v)
}

func decodeStd(d *Decoder, v reflect.Value) error {
	err := d.Skip()
	if err != nil {
		return err
	}
	u := v.Interface().(json.Unmarshaler)
	return u.UnmarshalJSON(d.Buf[d.start:d.pos])
}
