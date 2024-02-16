package blaze

import (
	"github.com/deveox/blaze/decoder"
	"github.com/deveox/blaze/encoder"
	"github.com/deveox/blaze/scopes"
)

func Unmarshal(data []byte, v interface{}) error {
	return decoder.Unmarshal(data, v)
}

func UnmarshalScoped(data []byte, v any, scope scopes.Decoding) error {
	return decoder.UnmarshalScoped(data, v, scope)
}

func UnmarshalScopedWithChanges(data []byte, v any, scope scopes.Decoding) ([]string, error) {
	return decoder.UnmarshalScopedWithChanges(data, v, scope)
}

func RegisterDecoder[T any](fn decoder.DecoderFn) {
	decoder.RegisterDecoder[T](fn)
}

func Marshal(v any) ([]byte, error) {
	return encoder.Marshal(v)
}

func RegisterEncoder[T any](fn encoder.EncoderFn) {
	encoder.RegisterEncoder[T](fn)
}
