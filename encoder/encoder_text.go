package encoder

import (
	"encoding"
	"reflect"
)

func encodePtrText(e *Encoder, v reflect.Value) error {
	if v.Kind() == reflect.Pointer && v.IsNil() {
		e.WriteString("null")
		return nil
	}
	return encodeText(e, v)
}

func encodeAddressableText(e *Encoder, v reflect.Value) error {
	va := v.Addr()
	// if va.IsNil() {
	// 	e.WriteString("null")
	// 	return nil
	// }
	return encodeText(e, va)
}

func encodeText(e *Encoder, v reflect.Value) error {
	m := v.Interface().(encoding.TextMarshaler)
	b, err := m.MarshalText()
	if err != nil {
		return err
	}
	return encodeString(e, b)
}
