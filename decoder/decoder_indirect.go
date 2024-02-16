package decoder

import "reflect"

func newIfAddressable(then, otherwise DecoderFn) DecoderFn {
	return func(d *Decoder, v reflect.Value) error {
		if v.CanAddr() {
			return then(d, v)
		}
		return otherwise(d, v)
	}
}

func decodeInterface(d *Decoder, v reflect.Value) error {
	if v.IsNil() {
		return d.ErrorF("[Blaze decodeInterface()] cannot decode into nil interface '%s'", v.Type())
	}
	return d.decode(v.Elem())
}

func decodePtr(d *Decoder, v reflect.Value) error {
	c := d.char()
	if c == 'n' {
		err := d.ScanNull()
		if err != nil {
			return err
		}
		if v.CanSet() {
			v.SetZero()
			return nil
		}
		d.pos -= 4
	}
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	return d.decode(v.Elem())
}

func decodeInvalid(d *Decoder, v reflect.Value) error {
	return d.ErrorF("[Blaze decodeInvalid()] cannot decode to unsupported type '%s'", v.Type())
}
