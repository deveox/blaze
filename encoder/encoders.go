package encoder

import (
	"encoding"
	"encoding/json"
	"reflect"
	"sync"

	"github.com/deveox/gu/async"
)

type EncoderFn func(*Encoder, reflect.Value) error

// TODO consider using a bitmap for better performance
var encoders = &async.Map[reflect.Type, EncoderFn]{}

func getEncoderFn(t reflect.Type) EncoderFn {
	if fi, ok := encoders.Load(t); ok {
		return fi
	}

	// To deal with recursive types, populate the map with an
	// indirect func before we build it. This type waits on the
	// real func (f) to be ready and then calls it. This indirect
	// func is only used for recursive types.
	var (
		wg sync.WaitGroup
		f  EncoderFn
	)
	wg.Add(1)
	fi, loaded := encoders.LoadOrStore(t, EncoderFn(func(e *Encoder, v reflect.Value) error {
		wg.Wait()
		return f(e, v)

	}))
	if loaded {
		return fi
	}

	// Compute the real encoder and replace the indirect func with it.
	f = newEncoderFn(t, true)
	wg.Done()
	encoders.Store(t, f)
	return f
}

type Marshaler interface {
	MarshalBlaze(e *Encoder) error
}

var (
	stdMarshaler  = reflect.TypeFor[json.Marshaler]()
	textMarshaler = reflect.TypeFor[encoding.TextMarshaler]()
	marshaler     = reflect.TypeFor[Marshaler]()
)

func newEncoderFn(t reflect.Type, allowAddr bool) EncoderFn {
	// If we have a non-pointer value whose type implements
	// Marshaler with a value receiver, then we're better off taking
	// the address of the value - otherwise we end up with an
	// allocation as we cast the value to an interface.
	if t.Kind() != reflect.Pointer && allowAddr {
		ptr := reflect.PointerTo(t)
		if ptr.Implements(marshaler) {
			return newCondAddrEncoder(encodeAddressableCustom, newEncoderFn(t, false))
		}
		if ptr.Implements(stdMarshaler) {
			return newCondAddrEncoder(encodeAddressableStd, newEncoderFn(t, false))
		}
		if ptr.Implements(textMarshaler) {
			return newCondAddrEncoder(encodeAddressableText, newEncoderFn(t, false))
		}
	}

	if t.Implements(marshaler) {
		return encodePtrCustom
	}

	if t.Implements(stdMarshaler) {
		return encodePtrStd
	}

	if t.Implements(textMarshaler) {
		return encodePtrText
	}

	switch t.Kind() {
	case reflect.Bool:
		return encodeBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encodeInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return encodeUint
	case reflect.Float32:
		return encodeFloat32
	case reflect.Float64:
		return encodeFloat64
	case reflect.String:
		return encodeString
	case reflect.Interface:
		return encodeInterface
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Pointer:
		return encodePtr
	default:
		return encodeUnsupported
	}
}

func RegisterEncoder[T any](fn EncoderFn) {
	encoders.Store(reflect.TypeFor[T](), fn)
}
