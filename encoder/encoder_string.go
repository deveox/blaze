package encoder

import "reflect"

func encodeStringOrBytes[T string | []byte](e *Encoder, v T) error {
	e.Grow(len(v) + 2)
	b := e.AvailableBuffer()
	b = append(b, '"')
	for i := 0; i < len(v); i++ {
		switch v[i] {
		case '"':
			b = append(b, '\\', '"')
		case '\\':
			b = append(b, '\\', '\\')
		case '\n':
			b = append(b, '\\', 'n')
		case '\r':
			b = append(b, '\\', 'r')
		case '\t':
			b = append(b, '\\', 't')
		default:
			b = append(b, v[i])
		}
	}
	b = append(b, '"')
	e.Write(b)
	return nil
}

func encodeString(e *Encoder, v reflect.Value) error {
	return encodeStringOrBytes(e, v.String())
}
