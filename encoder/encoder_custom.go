package encoder

import (
	"reflect"
)

func encodeAddressableCustom(e *Encoder, v reflect.Value) error {
	return encodeCustom(e, v.Addr())
}

func encodePtrCustom(e *Encoder, v reflect.Value) error {
	if v.Kind() == reflect.Pointer && v.IsNil() {
		e.WriteString("null")
		return nil
	}
	return encodeCustom(e, v)
}

func encodeCustom(e *Encoder, v reflect.Value) error {
	m := v.Interface().(Marshaler)
	return m.MarshalBlaze(e)
}
