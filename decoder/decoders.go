package decoder

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type decoderFunc func(*Decoder, reflect.Value) error

// TODO consider using a bitmap for better performance
type decoders map[string]decoderFunc

func (ds decoders) Get(v reflect.Value) (decoderFunc, error) {
	t := v.Type().String()
	t = strings.TrimLeft(t, "*")
	fn, ok := ds[t]
	if ok {
		return fn, nil
	}
	return nativeDecoder, nil
}

func (ds decoders) add(t reflect.Type, fn decoderFunc) {
	ts := strings.TrimLeft(t.String(), "*")
	ds[ts] = fn
}

func unmarshaler(d *Decoder, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Ptr:
		for v.Kind() == reflect.Ptr {
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			if v.Type().NumMethod() > 0 {
				val := v.Interface()
				if u, ok := val.(json.Unmarshaler); ok {
					err := d.Skip()
					if err != nil {
						return err
					}
					return u.UnmarshalJSON(d.Buf[d.start:d.pos])
				}
			}
			v = v.Elem()
		}
		return d.Error("[Blaze unmarshaler()] invalid attempt to decode using json.Unmarshaler, passed pointer(or pointer chain) does not implement json.Unmarshaler")
	default:
		if v.CanAddr() {
			v = v.Addr()
			if v.Type().NumMethod() > 0 {
				val := v.Interface()
				if u, ok := val.(json.Unmarshaler); ok {
					err := d.Skip()
					if err != nil {
						return err
					}
					return u.UnmarshalJSON(d.Buf[d.start:d.pos])
				}
			}
			return d.Error("[Blaze unmarshaler()] invalid attempt to decode using json.Unmarshaler, indirect value does not implement json.Unmarshaler")
		}
		return d.Error("[Blaze unmarshaler()] invalid attempt to decode using json.Unmarshaler, indirect value is not addressable")
	}
}

func nativeDecoder(d *Decoder, v reflect.Value) error {
	return d.nativeDecoder(v)
}

type CustomFn[T any] func(d *Decoder, bytes []byte) (T, error)

func customDecoder[T any](fn CustomFn[T]) decoderFunc {
	return func(d *Decoder, v reflect.Value) error {
		err := d.Skip()
		if err != nil {
			return err
		}
		val, err := fn(d, d.Buf[d.start:d.pos])
		if err != nil {
			return err
		}
		for v.Kind() == reflect.Ptr {
			if v.Type().Name() == "" {
				break
			}
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		}
		v.Set(reflect.ValueOf(val))
		return nil
	}
}

var Decoders = decoders{}

func RegisterUnmarshaler(value any) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr:
		if v.Type().NumMethod() > 0 {
			val := v.Interface()
			if _, ok := val.(json.Unmarshaler); ok {
				Decoders.add(v.Elem().Type(), unmarshaler)
				return
			}
			panic(fmt.Sprintf("blaze: can't register unmarshaler for type '%s' because it doesn't implement json.Unmarshaler", v.Type()))
		}
		panic(fmt.Sprintf("blaze: can't register unmarshaler for type '%s' because it doesn't have any methods", v.Type()))

	default:
		panic(fmt.Sprintf("blaze: can't register unmarshaler for type '%s', you should pass a pointer", v.Type()))
	}
}

func RegisterDecoder[T any](value T, fn CustomFn[T]) {
	Decoders.add(reflect.TypeOf(value), customDecoder(fn))
}

func init() {
	RegisterUnmarshaler(&time.Time{})
}
