package encoder

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func encodeInt(e *Encoder, v reflect.Value) error {
	va := v.Int()
	if va == 0 {
		e.WriteByte('0')
		return nil
	}
	b := e.AvailableBuffer()
	b = strconv.AppendInt(b, va, 10)
	e.Write(b)
	return nil
}

func encodeUint(e *Encoder, v reflect.Value) error {
	va := v.Uint()
	if va == 0 {
		e.WriteByte('0')
		return nil
	}
	b := e.AvailableBuffer()
	b = strconv.AppendUint(b, va, 10)
	e.Write(b)
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
	b := e.AvailableBuffer()
	sl := e.Len()
	b = strconv.AppendFloat(b, f, fmt, -1, int(bits))
	if fmt == 'e' {
		n := len(b) - sl
		// clean up e-09 to e-9
		if n >= 4 && b[n-4] == 'e' && b[n-3] == '-' && b[n-2] == '0' {
			b[n-2] = b[n-1]
			b = b[:n-1]
		}
	}
	return nil
}

func encodeFloat32(e *Encoder, v reflect.Value) error {
	return e.EncodeFloat(v, 32)
}

func encodeFloat64(e *Encoder, v reflect.Value) error {
	return e.EncodeFloat(v, 64)
}
