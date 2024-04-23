package encoder

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
)

func encodeByteSlice(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		e.WriteString("null")
		return nil
	}

	s := v.Bytes()
	e.bytes = append(e.bytes, '"')
	e.bytes = base64.StdEncoding.AppendEncode(e.bytes, s)
	e.bytes = append(e.bytes, '"')
	return nil
}

func encodeRawMessage(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		e.WriteString("null")
		return nil
	}
	e.bytes = append(e.bytes, v.Bytes()...)
	return nil
}

func newSliceEncoder(t reflect.Type) EncoderFn {
	// Byte slices get special treatment; arrays don't.
	if t.Elem() == reflect.TypeFor[byte]() {
		if t == reflect.TypeFor[json.RawMessage]() {
			return encodeRawMessage
		}
		return encodeByteSlice
	}
	return newArrayEncoder(t)
}

func newArrayEncoder(t reflect.Type) EncoderFn {
	vEnc := newEncoderFn(t.Elem(), true)
	return func(e *Encoder, v reflect.Value) error {
		if v.Kind() == reflect.Slice && v.IsNil() {
			e.WriteString("null")
			return nil
		}
		if v.IsZero() {
			e.bytes = append(e.bytes, '[', ']')
			return nil
		}
		return e.EncodeSlice(v, vEnc)
	}
}

func (e *Encoder) EncodeSlice(v reflect.Value, valueEnc EncoderFn) (err error) {
	e.depth++
	defer func() {
		e.depth--
	}()
	e.bytes = append(e.bytes, '[')
	n := v.Len()

	for i := 0; i < n; i++ {
		f := v.Index(i)
		oldLen := len(e.bytes)
		err = valueEnc(e, f)
		if err != nil {
			return err
		}
		if len(e.bytes) != oldLen {
			e.bytes = append(e.bytes, ',')
		}
	}

	last := len(e.bytes) - 1
	switch e.bytes[last] {
	case '[':
		if e.keep || e.depth == 1 {
			e.WriteByte(']')
		} else {
			e.bytes = e.bytes[:last]
		}
	case ',':
		e.bytes[last] = ']'
	default:
		e.WriteByte(']')
	}

	return nil
}
