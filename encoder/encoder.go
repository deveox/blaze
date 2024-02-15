package encoder

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"sync"

	"github.com/deveox/blaze/scopes"
	"github.com/deveox/blaze/types"
)

func Marshal(v any) ([]byte, error) {
	e := NewEncoder()
	defer encoders.Put(e)
	return e.Encode(v)
}
func MarshalScoped(v any, context scopes.Context) ([]byte, error) {
	e := NewEncoder()
	defer encoders.Put(e)
	e.ContextScope = context
	return e.Encode(v)
}

func NewEncoder() *Encoder {
	if v := encoders.Get(); v != nil {
		e := v.(*Encoder)
		e.bytes = e.bytes[:0]
		e.ContextScope = 0
		return e
	}
	return &Encoder{bytes: make([]byte, 0, 2048)}
}

var encoders sync.Pool

const MAX_DEPTH = 1000

type Encoder struct {
	bytes []byte

	ContextScope scopes.Context
	depth        int
	anonymous    bool
}

func (e *Encoder) Encode(v any) ([]byte, error) {
	err := e.encode(reflect.ValueOf(v))
	if err != nil {
		return nil, err
	}
	res := e.bytes[:len(e.bytes)]

	e.bytes = e.bytes[:0]
	return res, nil
}

func (e *Encoder) encode(v reflect.Value) error {
	if e.depth > MAX_DEPTH {
		return e.ErrorF("exceeded max depth of %d", MAX_DEPTH)
	}
	enc, err := EncoderFns.Get(v)
	if err != nil {
		return err
	}
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.Type().Name() != "" {
			break
		}
		if v.IsNil() {
			e.bytes = append(e.bytes, "null"...)
			return nil
		}
		v = v.Elem()

	}
	err = enc(e, v)
	if err != nil {
		return err
	}
	e.anonymous = false
	return nil
}

func (e *Encoder) nativeEncoder(v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		e.bytes = append(e.bytes, '"')
		str := v.String()
		for i := 0; i < len(str); i++ {
			switch str[i] {
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
				e.bytes = append(e.bytes, str[i])
			}
		}
		e.bytes = append(e.bytes, '"')
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := v.Int()
		if v == 0 {
			e.bytes = append(e.bytes, '0')
			return nil
		}
		e.bytes = strconv.AppendInt(e.bytes, v, 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := v.Uint()
		if v == 0 {
			e.bytes = append(e.bytes, '0')
			return nil
		}
		e.bytes = strconv.AppendUint(e.bytes, v, 10)
	case reflect.Float32:
		if v.IsZero() {
			e.bytes = append(e.bytes, '0')
			return nil
		}
		return e.EncodeFloat(v, 32)
	case reflect.Float64:
		if v.IsZero() {
			e.bytes = append(e.bytes, '0')
			return nil
		}
		return e.EncodeFloat(v, 64)
	case reflect.Bool:
		b := v.Bool()
		if b {
			e.bytes = append(e.bytes, "true"...)
		} else {
			e.bytes = append(e.bytes, "false"...)
		}
	case reflect.Struct:

		return e.EncodeStruct(v)
	case reflect.Slice, reflect.Array:
		if v.IsNil() {
			e.bytes = append(e.bytes, "null"...)
			return nil
		}
		if v.IsZero() {
			e.bytes = append(e.bytes, '[', ']')
			return nil
		}
		return e.EncodeSlice(v)
	case reflect.Interface, reflect.Ptr:
		for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
			if v.IsNil() {
				e.bytes = append(e.bytes, "null"...)
				return nil
			} else {
				v = v.Elem()
			}
		}
		return e.nativeEncoder(v)
	case reflect.Map:
		if v.IsNil() {
			e.bytes = append(e.bytes, "null"...)
			return nil
		}
		if v.IsZero() {
			e.bytes = append(e.bytes, '{', '}')
			return nil
		}
		return e.EncodeMap(v)
	// case reflect.Uintptr:

	// case reflect.UnsafePointer:

	// case reflect.Complex64, reflect.Complex128:

	// case reflect.Invalid:

	default:
		return fmt.Errorf("unsupported type: %v", v.Kind())
	}
	return nil
}

func (e *Encoder) EncodeFloat(v reflect.Value, bits int) error {
	f := v.Float()
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return fmt.Errorf("unsupported value: %v", f)
	}

	// Convert as if by ES6 number to string conversion.
	// This matches most other JSON generators.
	// See golang.org/issue/6384 and golang.org/issue/14135.
	// Like fmt %g, but the exponent cutoffs are different
	// and exponents themselves are not padded to two digits.
	abs := math.Abs(f)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if bits == 64 && (abs < 1e-6 || abs >= 1e21) || bits == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) {
			fmt = 'e'
		}
	}
	sl := len(e.bytes)
	e.bytes = strconv.AppendFloat(e.bytes, f, fmt, -1, int(bits))
	if fmt == 'e' {
		n := len(e.bytes) - sl
		// clean up e-09 to e-9
		if n >= 4 && e.bytes[n-4] == 'e' && e.bytes[n-3] == '-' && e.bytes[n-2] == '0' {
			e.bytes[n-2] = e.bytes[n-1]
			e.bytes = e.bytes[:n-1]
		}
	}
	return nil
}

func (e *Encoder) EncodeStruct(v reflect.Value) error {
	e.depth++
	defer func() {
		e.depth--
	}()

	if !e.anonymous {
		e.bytes = append(e.bytes, '{')
	}

	t := v.Type()
	si, _ := types.Cache.Get(t)

	for _, fi := range si.Fields {
		ok := fi.CheckEncoderScope(e.ContextScope)
		if !ok {
			continue
		}
		f := v.Field(fi.Idx)

		// Handle zero values

		if f.IsZero() && !fi.KeepEmpty {
			continue
		}

		switch f.Kind() {
		case reflect.Struct:
			if !fi.Anonymous {
				e.bytes = append(e.bytes, fi.ObjectKey...)
			}
			e.anonymous = fi.Anonymous
			if err := e.encode(f); err != nil {
				return err
			}
		default:
			e.bytes = append(e.bytes, fi.ObjectKey...)
			if err := e.encode(f); err != nil {
				return err
			}
		}

		e.bytes = append(e.bytes, ',')
	}
	if e.anonymous {
		if e.bytes[len(e.bytes)-1] == ',' {
			e.bytes = e.bytes[:len(e.bytes)-1]
		}
	} else {
		if e.bytes[len(e.bytes)-1] == ',' {
			e.bytes[len(e.bytes)-1] = '}'
		} else {
			e.bytes = append(e.bytes, '}')
		}
	}
	return nil
}

func (e *Encoder) EncodeSlice(v reflect.Value) error {
	e.bytes = append(e.bytes, '[')
	n := v.Len()
	elem := v.Type().Elem()
	enc, err := EncoderFns.GetType(elem)
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		f := v.Index(i)
		for f.Kind() == reflect.Interface || f.Kind() == reflect.Ptr {
			if f.Type().Name() != "" {
				break
			}
			if f.IsNil() {
				e.bytes = append(e.bytes, "null"...)
				return nil
			}
			f = f.Elem()

		}
		err = enc(e, f)
		if err != nil {
			return err
		}
		e.bytes = append(e.bytes, ',')
	}
	if e.bytes[len(e.bytes)-1] == ',' {
		e.bytes[len(e.bytes)-1] = ']'
	} else {
		e.bytes = append(e.bytes, ']')
	}

	return nil
}

func (e *Encoder) EncodeMap(v reflect.Value) error {
	e.bytes = append(e.bytes, '{')
	key := v.Type().Key()
	iter := v.MapRange()
	for {
		next := iter.Next()
		if next {
			switch key.Kind() {
			case reflect.String:
				if err := e.encode(iter.Key()); err != nil {
					return err
				}
			default:
				e.bytes = append(e.bytes, '"')
				if err := e.encode(iter.Key()); err != nil {
					return err
				}
				e.bytes = append(e.bytes, '"')
			}
			e.bytes = append(e.bytes, ':')
			if err := e.encode(iter.Value()); err != nil {
				return err
			}
			e.bytes = append(e.bytes, ',')
		} else {
			if e.bytes[len(e.bytes)-1] == ',' {
				e.bytes[len(e.bytes)-1] = '}'
			} else {
				e.bytes = append(e.bytes, '}')
			}
			break
		}
	}
	return nil
}

func (e *Encoder) Error(msg string) error {
	return errors.New(msg)
}

func (e *Encoder) ErrorF(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func AddIndent(b []byte) []byte {
	res := make([]byte, 0, len(b))
	depth := 0
	for i := 0; i < len(b); i++ {
		if b[i] == '{' || b[i] == '[' {
			depth++
			res = append(res, b[i], '\n')

		} else if b[i] == '}' || b[i] == ']' {
			res = append(res, '\n')

			res = append(res, b[i])
		} else {
			res = append(res, b[i])
		}
	}
	return res
}
