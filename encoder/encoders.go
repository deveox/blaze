package encoder

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/deveox/gu/async"
)

type encoderFunc func(*Encoder, reflect.Value) error

// TODO consider using a bitmap for better performance
var encoders = &async.Map[reflect.Type, encoderFunc]{}

func getEncoderFn(t reflect.Type) encoderFunc {
	if fi, ok := encoders.Load(t); ok {
		return fi
	}

	// To deal with recursive types, populate the map with an
	// indirect func before we build it. This type waits on the
	// real func (f) to be ready and then calls it. This indirect
	// func is only used for recursive types.
	var (
		wg sync.WaitGroup
		f  encoderFunc
	)
	wg.Add(1)
	fi, loaded := encoders.LoadOrStore(t, encoderFunc(func(e *Encoder, v reflect.Value) error {
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
	MarshalBlaze(e *Encoder) ([]byte, error)
}

var (
	stdMarshaler  = reflect.TypeFor[json.Marshaler]()
	textMarshaler = reflect.TypeFor[encoding.TextMarshaler]()
	marshaler     = reflect.TypeFor[Marshaler]()
)

func newEncoderFn(t reflect.Type, allowAddr bool) encoderFunc {
	// If we have a non-pointer value whose type implements
	// Marshaler with a value receiver, then we're better off taking
	// the address of the value - otherwise we end up with an
	// allocation as we cast the value to an interface.
	if t.Kind() != reflect.Pointer && allowAddr && reflect.PointerTo(t).Implements(stdMarshaler) {
		return newCondAddrEncoder(encodeAddressableStd, newEncoderFn(t, false))
	}
	if t.Implements(stdMarshaler) {
		return encodePtrStd
	}
	if t.Kind() != reflect.Pointer && allowAddr && reflect.PointerTo(t).Implements(textMarshaler) {
		return newCondAddrEncoder(encodeAddressableText, newEncoderFn(t, false))
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
		return interfaceEncoder
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Pointer:
		return newPtrEncoder(t)
	default:
		return unsupportedTypeEncoder
	}
}

func Get(v reflect.Value) (encoderFunc, error) {
	for v.Kind() == reflect.Interface {
		v = v.Elem()

	}
	t := v.Type()
	fn, ok := encoders.Load(t)
	if ok {
		return fn, nil
	}
	return nativeEncoder, nil
}

func newEncoder(v reflect.Value) (encoderFunc, error) {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:

	}
}

// func marshaler(d *Encoder, v reflect.Value) error {
// 	if v.CanInterface() {
// 		val := v.Interface()
// 		if u, ok := val.(json.Marshaler); ok {
// 			b, err := u.MarshalJSON()
// 			if err != nil {
// 				return err
// 			}
// 			d.bytes = append(d.bytes, b...)
// 			return nil
// 		}
// 		return d.ErrorF("[Blaze marshaler()] invalid attempt to encode using json.Marshaler, type '%s' does not implement json.Marshaler", v.Type())
// 	}
// 	return d.ErrorF("[Blaze marshaler()] invalid attempt to encode using json.Marshaler, type '%s' is not addressable", v.Type())
// }

func nativeEncoder(d *Encoder, v reflect.Value) error {
	return d.nativeEncoder(v)
}

type CustomFn[T any] func(d *Encoder, v T) ([]byte, error)

func customEncoder[T any](fn CustomFn[T]) encoderFunc {
	return func(e *Encoder, v reflect.Value) error {
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
		if v.CanInterface() {
			b, err := fn(e, v.Interface().(T))
			if err != nil {
				return err
			}
			e.bytes = append(e.bytes, b...)
			return nil
		}
		return e.ErrorF("[Blaze customEncoder()] can't encode value of type '%s' because it's not addressable", v.Type())
	}
}

var EncoderFns = encoderFns{}

func RegisterMarshaler(value any) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr:
		if v.Type().NumMethod() > 0 {
			val := v.Interface()
			if _, ok := val.(json.Marshaler); ok {
				EncoderFns.add(v.Elem().Type(), marshaler)
				return
			}
			panic(fmt.Sprintf("blaze: can't register unmarshaler for type '%s' because it doesn't implement json.Unmarshaler", v.Type()))
		}
		panic(fmt.Sprintf("blaze: can't register unmarshaler for type '%s' because it doesn't have any methods", v.Type()))

	default:
		panic(fmt.Sprintf("blaze: can't register unmarshaler for type '%s', you should pass a pointer", v.Type()))
	}
}

func RegisterEncoder[T any](value T, fn CustomFn[T]) {
	EncoderFns.add(reflect.TypeOf(value), customEncoder(fn))
}

func init() {
	RegisterMarshaler(&time.Time{})
}
