package decoder

import (
	"reflect"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

func (d *Decoder) SkipString() error {
	d.pos++
	for {
		c := d.char()
		switch c {
		case '"':
			d.pos++
			return nil
		case '\\':
			d.pos += 2
		case TERMINATION_CHAR:
			return d.Error("[Blaze SkipString()] unexpected end of input, expected '\"'")
		default:
			d.pos++
		}
	}
}

func (d *Decoder) DecodeString() (string, error) {
	d.start = d.pos
	d.pos = d.start + 1
	for {
		c := d.char()
		switch c {
		case '"':
			str := string(d.Buf[d.start+1 : d.pos])
			d.pos++
			return str, nil
		case '\\':
			b, err := d.unquoteString()
			if err != nil {
				return "", err
			}
			return string(b), nil
		case TERMINATION_CHAR:
			return "", d.Error("[Blaze DecodeString()] unexpected end of input, expected '\"'")
		}
		d.pos++
	}
}

func (d *Decoder) unquoteString() ([]byte, error) {
	str := d.Buf[d.start+1 : d.pos]
loop:
	for {
		c := d.char()
		switch c {
		case '"':
			return str, nil
		case '\\':
			d.pos++
			switch d.Buf[d.pos] {
			case 'u':
				d.pos++
				// Read 4 bytes assuming it's a UTF-8 code unit.
				r, err := d.getU4()
				if err != nil {
					return str, err
				}
				// Check if rune can be a surrogate.
				if utf16.IsSurrogate(r) {
					r2, err := d.getU4()
					if err != nil {
						return str, err
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
				return str, d.Error("[Blaze decodeStringWithEscapes()] invalid escape character")
			}
		case TERMINATION_CHAR:
			d.Error("[Blaze decodeStringWithEscapes()] unexpected end of input, expected '\"'")
		default:
			str = append(str, d.Buf[d.pos])

		}
		d.pos++
	}
}

func (t *Decoder) getU4() (rune, error) {
	var r rune
	for i := 0; i < 4; i++ {
		c := t.char()
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

func decodeString(d *Decoder, v reflect.Value) error {
	d.SkipWhitespace()
	c := d.char()
	switch c {
	case '"':
	case 'n':
		err := d.ScanNull()
		if err != nil {
			return err
		}
		v.SetZero()
		return nil
	default:
		return d.Error("[Blaze decodeString()] invalid char, expected '\"'")
	}
	str, err := d.DecodeString()
	if err != nil {
		return err
	}
	v.SetString(str)
	return nil
}
