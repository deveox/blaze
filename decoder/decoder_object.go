package decoder

import (
	"fmt"
	"reflect"

	"github.com/deveox/blaze/types"
)

func (t *Decoder) ScanObject() (int, error) {
	start := t.pos
	size := 0
	t.pos++
	for {

		c := t.char(t.ptr, t.pos)
		switch c {
		case ',':
			t.pos++
		case ':':
			t.pos++
			size++
		case '}':
			t.pos++
			t.start = start
			return size, nil
		case TERMINATION_CHAR:
			return 0, t.Error("[Blaze ScanObject()] unexpected end of input, expected '}'")
		default:
			err := t.Skip()
			if err != nil {
				return 0, err
			}
		}
	}
}

func (t *Decoder) SkipObject() error {
	level := 1
	for {
		t.pos++
		c := t.char(t.ptr, t.pos)
		switch c {
		case '{':
			level++
			t.depth++
			if t.depth > MAX_DEPTH {
				return t.Error("[Blaze SkipObject()] maximum depth reached")
			}
		case '}':
			level--
			if level == 0 {
				t.pos++
				return nil
			}
		case '"':
			t.SkipString()
			t.pos--
		case TERMINATION_CHAR:
			return t.Error("[Blaze SkipObject()] unexpected end of input, expected '}'")
		}
	}
}

func (t *Decoder) decodeObject(v reflect.Value) error {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Map:
		return t.decodeMap(v)
	case reflect.Struct:
		return t.decodeStruct(v)
	default:
		return t.ErrorF("[Blaze decodeObject()] can't decode JSON object into Go Type '%s'", v.Type().String())
	}
}

func (t *Decoder) decodeMap(v reflect.Value) error {
	size, err := t.ScanObject()
	if err != nil {
		return err
	}
	if v.IsNil() {
		v.Set(reflect.MakeMapWithSize(v.Type(), size))
	}
	t.pos = t.start + 1
	for {
		t.SkipWhitespace()
		c := t.char(t.ptr, t.pos)
		switch c {
		case '}':
			t.pos++
			return nil
		case '"':
		case TERMINATION_CHAR:
			return t.Error("[Blaze Decode Map] unexpected end of input, expected object key")
		default:
			return t.Error("[Blaze Decode Map] expected object key")
		}

		key := reflect.New(v.Type().Key()).Elem()
		switch key.Kind() {
		case reflect.String:
			err := t.decodeString(key)
			if err != nil {
				return err
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			start := t.pos + 1
			err := t.SkipString()
			if err != nil {
				return err
			}
			t.start = start
			t.pos--
			switch key.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				err := t.decodeInt(key)
				if err != nil {
					return err
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				err := t.decodeUint(key)
				if err != nil {
					return err
				}
			case reflect.Float32, reflect.Float64:
				err := t.decodeFloat(key)
				if err != nil {
					return err
				}
			}
			t.pos++
		default:
			return t.ErrorF("[Blaze decodeMap()] can't decode JSON object key into Go Type '%s'", key.Type().String())
		}
		t.SkipWhitespace()
		c = t.char(t.ptr, t.pos)
		if c != ':' {
			return t.Error("[Blaze decodeMap()] expected ':'")
		}
		t.pos++
		t.SkipWhitespace()
		value := reflect.New(v.Type().Elem()).Elem()
		if err := t.decode(value); err != nil {
			return err
		}
		v.SetMapIndex(key, value)
		t.SkipWhitespace()
		c = t.char(t.ptr, t.pos)
		switch c {
		case '}':
			t.pos++
			return nil
		case ',':
			t.pos++
		case TERMINATION_CHAR:
			return t.Error("[Blaze decodeMap()] unexpected end of input, expected ',' or '}'")
		default:
			return t.Error("[Blaze decodeMap()] expected ',' or '}'")
		}
	}
}

func (d *Decoder) decodeStruct(v reflect.Value) error {
	si, err := types.Cache.Get(v.Type())
	if err != nil {
		return err
	}
	d.pos = d.pos + 1

	prefix := d.ChangesPrefix
	for {
		d.SkipWhitespace()
		c := d.char(d.ptr, d.pos)
		switch c {
		case '}':
			d.pos++
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
		c = d.char(d.ptr, d.pos)
		if c != ':' {
			return d.Error("[Blaze decodeStruct()] expected ':'")
		}
		d.pos++
		d.SkipWhitespace()
		field, embedded, ok := si.GetDecoderField(fName, d.ContextScope, d.OperationScope)
		// fmt.Printf("\nfield %v %s %#v\n\n", ok, v.Type(), field)
		if ok {
			var fv reflect.Value
			if d.Changes != nil {
				if prefix == "" {
					d.ChangesPrefix = fName
				} else {
					d.ChangesPrefix = fmt.Sprintf("%s.%s", prefix, fName)
				}
				d.Changes = append(d.Changes, d.ChangesPrefix)
			}
			if embedded != nil {
				fv = embedded.Value(v)
			} else {
				fv = v.Field(field.Idx)
			}
			if err := d.decode(fv); err != nil {
				return err
			}
		} else {

			err := d.Skip()
			if err != nil {
				return err
			}
		}
		// fmt.Println(1, string(d.Buf[d.pos:]))
		d.SkipWhitespace()
		c = d.char(d.ptr, d.pos)
		switch c {
		case '}':
			d.pos++
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
