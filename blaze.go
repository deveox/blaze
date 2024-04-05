package blaze

import (
	"github.com/deveox/blaze/decoder"
	"github.com/deveox/blaze/encoder"
	"github.com/deveox/blaze/scopes"
)

var AdminDecoder = &decoder.Config{Scope: scopes.CONTEXT_ADMIN}

func Unmarshal(data []byte, v interface{}) error {
	return AdminDecoder.Unmarshal(data, v)
}

func UnmarshalScoped(data []byte, v any, scope scopes.Decoding) error {
	return AdminDecoder.UnmarshalScoped(data, v, scope)
}

func UnmarshalScopedWithChanges(data []byte, v any, scope scopes.Decoding) ([]string, error) {
	return AdminDecoder.UnmarshalScopedWithChanges(data, v, scope)
}

func RegisterDecoder[T any](fn decoder.DecoderFn) {
	decoder.RegisterDecoder[T](fn)
}

var AdminEncoder = &encoder.Config{Scope: scopes.CONTEXT_ADMIN}

func Marshal(v any) ([]byte, error) {
	return AdminEncoder.Marshal(v)
}

func MarshalPartial(v any, fields []string, short bool) ([]byte, error) {
	return AdminEncoder.MarshalPartial(v, fields, short)
}

func RegisterEncoder[T any](fn encoder.EncoderFn) {
	encoder.RegisterEncoder[T](fn)
}
