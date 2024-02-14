package decoder

import (
	"reflect"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

func (t *Decoder) SkipString() error {
	t.pos++
	for {
		c := t.char(t.ptr, t.pos)
		switch c {
		case '"':
			t.pos++
			return nil
		case '\\':
			t.pos += 2
		case TERMINATION_CHAR:
			return t.Error("[Blaze SkipString()] unexpected end of input, expected '\"'")
		default:
			t.pos++
		}
	}
}

func (t *Decoder) DecodeString() (string, error) {
	t.start = t.pos
	t.pos = t.start + 1
	for {
		c := t.char(t.ptr, t.pos)
		switch c {
		case '"':
			str := string(t.Buf[t.start+1 : t.pos])
			t.pos++
			return str, nil
		case '\\':
			return t.decodeStringWithEscapes()
		case TERMINATION_CHAR:
			return "", t.Error("[Blaze DecodeString()] unexpected end of input, expected '\"'")
		}
		t.pos++
	}
}

func (t *Decoder) decodeStringWithEscapes() (string, error) {
	str := t.Buf[t.start+1 : t.pos]
loop:
	for {
		c := t.char(t.ptr, t.pos)
		switch c {
		case '"':
			return string(str), nil
		case '\\':
			t.pos++
			switch t.Buf[t.pos] {
			case 'u':
				t.pos++
				// Read 4 bytes assuming it's a UTF-8 code unit.
				r, err := t.getU4()
				if err != nil {
					return "", err
				}
				// Check if rune can be a surrogate.
				if utf16.IsSurrogate(r) {
					r2, err := t.getU4()
					if err != nil {
						return "", err
					}
					// If it's a valid surrogate pair, decode it.
					if dec := utf16.DecodeRune(r, r2); dec != unicode.ReplacementChar {
						str = utf8.AppendRune(str, dec)
						continue loop
					}
					// Invalid surrogate; fall back to replacement rune.
					r = unicode.ReplacementChar
				}
				str = utf8.AppendRune(str, r)
				continue loop
			case '"':
				str = append(str, '"')
			case '\\':
				str = append(str, '\\')
			case '/':
				str = append(str, '/')
			case 'b':
				str = append(str, '\b')
			case 'f':
				str = append(str, '\f')
			case 'n':
				str = append(str, '\n')
			case 'r':
				str = append(str, '\r')
			case 't':
				str = append(str, '\t')
			default:
				return "", t.Error("[Blaze decodeStringWithEscapes()] invalid escape character")
			}
		case TERMINATION_CHAR:
			t.Error("[Blaze decodeStringWithEscapes()] unexpected end of input, expected '\"'")
		default:
			str = append(str, t.Buf[t.pos])

		}
		t.pos++
	}
}

func (t *Decoder) getU4() (rune, error) {
	var r rune
	for i := 0; i < 4; i++ {
		c := t.char(t.ptr, t.pos)
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			c = c - '0'
		case 'a', 'b', 'c', 'd', 'e', 'f':
			c = c - 'a' + 10
		case 'A', 'B', 'C', 'D', 'E', 'F':
			c = c - 'A' + 10
		case TERMINATION_CHAR:
			return -1, t.Error("[Blaze getU4()] unexpected end of input, expected UTF-8 code unit (hexadecimal digit)")
		default:
			return -1, t.Error("[Blaze getU4()] invalid character, expected UTF-8 code unit (hexadecimal digit)")
		}
		r = r*16 + rune(c)
		t.pos++
	}
	return r, nil
}

func (t *Decoder) decodeString(v reflect.Value) error {
	str, err := t.DecodeString()
	if err != nil {
		return err
	}
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String:
		v.SetString(str)
	case reflect.Interface:
		if v.NumMethod() > 0 {
			return t.ErrorF("[Blaze decodeString()] unsupported custom interface %s", v.Kind().String())
		}
		v.Set(reflect.ValueOf(str))
	default:
		return t.ErrorF("[Blaze decodeString()] can't decode string into Go type  '%s'", v.Type())
	}
	return nil
}
