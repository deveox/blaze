package blaze

import (
	"github.com/deveox/blaze/decoder"
	"github.com/deveox/blaze/encoder"
	"github.com/deveox/blaze/scopes"
)

func Unmarshal(data []byte, v interface{}) error {
	return decoder.Unmarshal(data, v)
}

func UnmarshalScoped(data []byte, v any, context scopes.Context, scope scopes.Decoding) error {
	return decoder.UnmarshalScoped(data, v, context, scope)
}

func UnmarshalScopedWithChanges(data []byte, v any, context scopes.Context, scope scopes.Decoding) ([]string, error) {
	return decoder.UnmarshalScopedWithChanges(data, v, context, scope)
}

func RegisterDecoder[T any](value T, fn decoder.CustomFn[T]) {
	decoder.RegisterDecoder(value, fn)
}

func RegisterUnmarshaler(value any) {
	decoder.RegisterUnmarshaler(value)
}

var DBDecoder = &decoder.Config{
	ContextScope: scopes.CONTEXT_DB,
}

var ClientDecoder = &decoder.Config{
	ContextScope: scopes.CONTEXT_CLIENT,
}

var AdminDecoder = &decoder.Config{
	ContextScope: scopes.CONTEXT_ADMIN,
}

func Marshal(v any) ([]byte, error) {
	return encoder.Marshal(v)
}

func MarshalScoped(v any, context scopes.Context) ([]byte, error) {
	return encoder.MarshalScoped(v, context)
}

var DBEncoder = &encoder.Config{
	ContextScope: scopes.CONTEXT_DB,
}

var ClientEncoder = &encoder.Config{
	ContextScope: scopes.CONTEXT_CLIENT,
}

var AdminEncoder = &encoder.Config{
	ContextScope: scopes.CONTEXT_ADMIN,
}

func RegisterEncoder[T any](value T, fn encoder.CustomFn[T]) {
	encoder.RegisterEncoder(value, fn)
}

func RegisterMarshaler(value any) {
	encoder.RegisterMarshaler(value)
}
