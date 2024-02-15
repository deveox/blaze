package encoder

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type encoderFunc func(*Encoder, reflect.Value) error

// TODO consider using a bitmap for better performance
type encoderFns map[string]encoderFunc

func (ds encoderFns) Get(v reflect.Value) (encoderFunc, error) {
	t := v.Type().String()
	t = strings.TrimLeft(t, "*")
	fn, ok := ds[t]
	if ok {
		return fn, nil
	}
	return nativeEncoder, nil
}

func (ds encoderFns) GetType(t reflect.Type) (encoderFunc, error) {
	ts := t.String()
	ts = strings.TrimLeft(ts, "*")
	fn, ok := ds[ts]
	if ok {
		return fn, nil
	}
	return nativeEncoder, nil
}

func (ds encoderFns) add(t reflect.Type, fn encoderFunc) {
	ts := strings.TrimLeft(t.String(), "*")
	ds[ts] = fn
}

func marshaler(d *Encoder, v reflect.Value) error {
	if v.CanInterface() {
		val := v.Interface()
		if u, ok := val.(json.Marshaler); ok {
			b, err := u.MarshalJSON()
			if err != nil {
				return err
			}
			d.bytes = append(d.bytes, b...)
			return nil
		}
		return d.ErrorF("[Blaze marshaler()] invalid attempt to encode using json.Marshaler, type '%s' does not implement json.Marshaler", v.Type())
	}
	return d.ErrorF("[Blaze marshaler()] invalid attempt to encode using json.Marshaler, type '%s' is not addressable", v.Type())
}

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
