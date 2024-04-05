package blaze

import (
	"github.com/deveox/blaze/ctx"
	"github.com/deveox/blaze/decoder"
	"github.com/deveox/blaze/encoder"
	"github.com/deveox/blaze/scopes"
)

var AdminDecoder = &decoder.Config{Scope: scopes.CONTEXT_ADMIN}

func Unmarshal(data []byte, v interface{}) error {
	return AdminDecoder.Unmarshal(data, v)
}

func UnmarshalCtx(data []byte, v interface{}, ctx *ctx.Ctx) error {
	return AdminDecoder.UnmarshalCtx(data, v, ctx)
}

func UnmarshalScoped(data []byte, v any, scope scopes.Decoding) error {
	return AdminDecoder.UnmarshalScoped(data, v, scope)
}

func UnmarshalScopedCtx(data []byte, v any, scope scopes.Decoding, ctx *ctx.Ctx) error {
	return AdminDecoder.UnmarshalScopedCtx(data, v, scope, ctx)
}

func UnmarshalScopedWithChanges(data []byte, v any, scope scopes.Decoding) ([]string, error) {
	return AdminDecoder.UnmarshalScopedWithChanges(data, v, scope)
}

func UnmarshalScopedWithChangesCtx(data []byte, v any, scope scopes.Decoding, ctx *ctx.Ctx) ([]string, error) {
	return AdminDecoder.UnmarshalScopedWithChangesCtx(data, v, scope, ctx)
}

func RegisterDecoder[T any](fn decoder.DecoderFn) {
	decoder.RegisterDecoder[T](fn)
}

var AdminEncoder = &encoder.Config{Scope: scopes.CONTEXT_ADMIN}

func Marshal(v any) ([]byte, error) {
	return AdminEncoder.Marshal(v)
}

func MarshalCtx(v any, ctx *ctx.Ctx) ([]byte, error) {
	return AdminEncoder.MarshalCtx(v, ctx)
}

func MarshalPartial(v any, fields []string, short bool) ([]byte, error) {
	return AdminEncoder.MarshalPartial(v, fields, short)
}

func MarshalPartialCtx(v any, fields []string, short bool, ctx *ctx.Ctx) ([]byte, error) {
	return AdminEncoder.MarshalPartialCtx(v, fields, short, ctx)
}

func RegisterEncoder[T any](fn encoder.EncoderFn) {
	encoder.RegisterEncoder[T](fn)
}

func DecCtx[T any](d *decoder.Decoder, key string) (res T, ok bool) {
	ok, v := d.Get(key)
	if !ok {
		return res, false
	}
	res, ok = v.(T)
	return res, ok
}

func EncCtx[T any](e *encoder.Encoder, key string) (res T, ok bool) {
	ok, v := e.Get(key)
	if !ok {
		return res, false
	}
	res, ok = v.(T)
	return res, ok
}
