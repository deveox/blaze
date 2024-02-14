package blaze

import (
	"github.com/deveox/blaze/decoder"
	"github.com/deveox/blaze/encoder"
	"github.com/deveox/blaze/scopes"
)

func Unmarshal(data []byte, v interface{}) error {
	return decoder.Unmarshal(data, v)
}

func UnmarshalScoped(data []byte, v any, scopes scopes.Scopes) error {
	return decoder.UnmarshalScoped(data, v, scopes)
}

func UnmarshalOperation(data []byte, v any, operation scopes.Operation) error {
	return decoder.UnmarshalOperation(data, v, operation)
}

var DBDecoder = &decoder.Config{
	ContextScope: scopes.CONTEXT_DB,
}

var ClientDecoder = &decoder.Config{
	UserScope: scopes.USER_CLIENT,
}

func Marshal(v any) ([]byte, error) {
	return encoder.Marshal(v)
}

func MarshalScoped(v any, scopes scopes.Scopes) ([]byte, error) {
	return encoder.MarshalScoped(v, scopes)
}

var DBEncoder = &encoder.Config{
	ContextScope: scopes.CONTEXT_DB,
}

var ClientEncoder = &encoder.Config{
	UserScope: scopes.USER_CLIENT,
}
