package encoder

import "reflect"

func encodeStringOrBytes[T string | []byte](e *Encoder, v T) error {
	e.WriteByte('"')
	for i := 0; i < len(v); i++ {
		switch v[i] {
		case '"':
			e.bytes = append(e.bytes, '\\', '"')
		case '\\':
			e.bytes = append(e.bytes, '\\', '\\')
		case '\n':
			e.bytes = append(e.bytes, '\\', 'n')
		case '\r':
			e.bytes = append(e.bytes, '\\', 'r')
		case '\t':
			e.bytes = append(e.bytes, '\\', 't')
		default:
			e.bytes = append(e.bytes, v[i])
		}
	}
	e.WriteByte('"')
	return nil
}

func encodeString(e *Encoder, v reflect.Value) error {
	return encodeStringOrBytes(e, v.String())
}
