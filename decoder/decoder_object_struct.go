package decoder

import (
	"fmt"
	"reflect"

	"github.com/deveox/blaze/types"
)

func (d *Decoder) decodeStruct(v reflect.Value, si *types.Struct) error {
	d.depth++
	if d.depth > MAX_DEPTH {
		return d.Error("[Blaze decodeStruct()] max depth reached")
	}
	d.SkipWhitespace()
	c := d.char()
	prefix := d.ChangesPrefix
	switch c {
	case '{':
		d.pos++
	case 'n':
		err := d.ScanNull()
		if err != nil {
			return err
		}
		for _, fi := range si.Fields {
			ok := fi.Field.CheckDecoderScope(d.config.Scope, d.operation)
			if ok {
				f := fi.Value(v)
				if f.IsZero() {
					continue
				}
				if d.Changes != nil {
					if prefix == "" {
						d.Changes = append(d.Changes, fi.Field.Name)
					} else {
						d.Changes = append(d.Changes, fmt.Sprintf("%s.%s", prefix, fi.Field.Name))
					}
				}
				f.SetZero()
			}
		}
		return nil
	default:
		return d.Error("[Blaze decodeStruct()] expected '{' or 'null'")
	}

	for {
		d.SkipWhitespace()
		c := d.char()
		switch c {
		case '}':
			d.pos++
			d.depth--
			return nil
		case '"':
		case TERMINATION_CHAR:
			return d.Error("[Blaze decodeStruct()] unexpected end of input, expected object key or '}'")
		default:
			return d.Error("[Blaze decodeStruct()] expected object key or '}'")
		}
		start := d.pos
		err := d.SkipString()
		if err != nil {
			return err
		}
		fName := BytesToString(d.Buf[start+1 : d.pos-1])
		d.SkipWhitespace()
		c = d.char()
		if c != ':' {
			return d.Error("[Blaze decodeStruct()] expected ':'")
		}
		d.pos++
		d.SkipWhitespace()
		field, ok := si.GetDecoderField(fName, d.config.Scope, d.operation)
		// fmt.Printf("\nfield %v %s %#v\n\n", ok, v.Type(), field)
		if ok {
			fv := field.Value(v)
			if d.Changes != nil {
				if prefix == "" {
					d.ChangesPrefix = field.Field.Name
				} else {
					d.ChangesPrefix = fmt.Sprintf("%s.%s", prefix, field.Field.Name)
				}
				d.Changes = append(d.Changes, d.ChangesPrefix)
			}
			oldLen := len(d.Changes)
			if field.Field.StringDecoding && d.char() == '"' {
				s, err := d.DecodeString()
				if err != nil {
					return err
				}
				nd := d.Decoder([]byte(s))
				if err := nd.decode(fv); err != nil {
					nd.Release()
					return err
				}
				nd.Release()
			} else {
				if err := d.decode(fv); err != nil {
					return err
				}
			}

			if field.Field.Struct != nil {
				if len(d.Changes) == oldLen && len(d.Changes) > 0 {
					d.Changes = d.Changes[:len(d.Changes)-1]
				}
			}

		} else {
			err := d.Skip()
			if err != nil {
				return err
			}
		}
		// fmt.Println(1, string(d.Buf[d.pos:]))
		d.SkipWhitespace()
		c = d.char()
		switch c {
		case '}':
			d.pos++
			d.depth--
			return nil
		case ',':
			d.pos++
		case TERMINATION_CHAR:
			return d.Error("[Blaze decodeStruct()] unexpected end of input, expected ',' or '}'")
		default:
			return d.Error("[Blaze decodeStruct()] expected ',' or '}'")
		}
	}
}

func newStructDecoder(t reflect.Type) DecoderFn {
	si := types.Cache.Get(t)
	return func(d *Decoder, v reflect.Value) error {
		return d.decodeStruct(v, si)
	}
}
