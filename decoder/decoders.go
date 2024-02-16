package decoder

import (
	"encoding"
	"encoding/json"
	"reflect"
	"sync"

	"github.com/deveox/gu/async"
)

type DecoderFn func(*Decoder, reflect.Value) error

// TODO consider using a bitmap for better performance

var decoders = &async.Map[reflect.Type, DecoderFn]{}

func getDecoderFn(t reflect.Type) DecoderFn {
	if fi, ok := decoders.Load(t); ok {
		return fi
	}

	// To deal with recursive types, populate the map with an
	// indirect func before we build it. This type waits on the
	// real func (f) to be ready and then calls it. This indirect
	// func is only used for recursive types.
	var (
		wg sync.WaitGroup
		f  DecoderFn
	)
	wg.Add(1)
	fi, loaded := decoders.LoadOrStore(t, DecoderFn(func(e *Decoder, v reflect.Value) error {
		wg.Wait()
		return f(e, v)

	}))
	if loaded {
		return fi
	}

	// Compute the real encoder and replace the indirect func with it.
	f = newDecoderFn(t, true)
	wg.Done()
	decoders.Store(t, f)
	return f
}

type Unmarshaler interface {
	UnmarshalBlaze(e *Decoder, data []byte) error
}

var (
	stdUnmarshaler  = reflect.TypeFor[json.Unmarshaler]()
	textUnmarshaler = reflect.TypeFor[encoding.TextMarshaler]()
	unmarshaler     = reflect.TypeFor[Unmarshaler]()
)

func newDecoderFn(t reflect.Type, allowAddr bool) DecoderFn {
	// If we have a non-pointer value whose type implements
	// Marshaler with a value receiver, then we're better off taking
	// the address of the value - otherwise we end up with an
	// allocation as we cast the value to an interface.
	if t.Kind() != reflect.Pointer && allowAddr {
		ptr := reflect.PointerTo(t)
		if ptr.Implements(unmarshaler) {
			return newIfAddressable(decodeAddressableCustom, newDecoderFn(t, false))
		}
		if t.Implements(stdUnmarshaler) {
			return newIfAddressable(decodeAddressableStd, newDecoderFn(t, false))
		}
		if t.Implements(textUnmarshaler) {
			return newIfAddressable(decodeAddressableText, newDecoderFn(t, false))
		}

	}

	if t.Implements(unmarshaler) {
		return decodePtrCustom
	}

	if t.Implements(stdUnmarshaler) {
		return decodePtrStd
	}

	if t.Implements(textUnmarshaler) {
		return decodePtrText
	}

	switch t.Kind() {
	case reflect.Bool:
		return decodeBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return decodeInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return decodeUint
	case reflect.Float32, reflect.Float64:
		return decodeFloat
	case reflect.String:
		return decodeString
	case reflect.Interface:
		if t.NumMethod() == 0 {
			return decodeAny
		}
		return decodeInterface
	case reflect.Struct:
		return newStructDecoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Slice:
		return newSliceDecoder(t)
	case reflect.Array:
		return newArrayDecoder(t)
	case reflect.Pointer:
		return decodePtr
	default:
		return decodeInvalid
	}
}

func RegisterDecoder[T any](fn DecoderFn) {
	decoders.Store(reflect.TypeFor[T](), fn)
}
