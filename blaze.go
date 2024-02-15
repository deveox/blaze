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
