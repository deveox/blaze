package encoder

import (
	"encoding/json"
	"reflect"
)

func encodeAddressableStd(e *Encoder, v reflect.Value) error {
	va := v.Addr()
	// if va.IsNil() {
	// 	e.WriteString("null")
	// 	return nil
	// }
	return encodeStd(e, va)
}

func encodePtrStd(e *Encoder, v reflect.Value) error {
	if v.Kind() == reflect.Pointer && v.IsNil() {
		e.WriteString("null")
		return nil
	}
	return encodeStd(e, v)
}

func encodeStd(e *Encoder, v reflect.Value) error {
	m := v.Interface().(json.Marshaler)
	b, err := m.MarshalJSON()
	if err != nil {
		return err
	}
	e.Write(b)
	return nil
}
