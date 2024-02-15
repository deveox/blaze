package encoder

import "reflect"

func encodeBool(e *Encoder, v reflect.Value) error {
	b := v.Bool()
	if b {
		e.WriteString("true")
	} else {
		e.WriteString("false")
	}
	return nil
}
