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
	e.bytes = strconv.AppendInt(e.bytes, va, 10)
	return nil
}

func encodeUint(e *Encoder, v reflect.Value) error {
	va := v.Uint()
	if va == 0 {
		e.WriteByte('0')
		return nil
	}
	e.bytes = strconv.AppendUint(e.bytes, va, 10)
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

func encodeFloat32(e *Encoder, v reflect.Value) error {
	return e.EncodeFloat(v, 32)
}

func encodeFloat64(e *Encoder, v reflect.Value) error {
	return e.EncodeFloat(v, 64)
}
