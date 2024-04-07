package decoder

import (
	"sync"

	"github.com/deveox/blaze/ctx"
	"github.com/deveox/blaze/scopes"
)

// Config is a configuration for the decoder.
// It can be used to define the scope of the decoding ones and reuse it multiple times.
type Config struct {
	Scope       scopes.Context
	decoderPool sync.Pool
}

// Unmarshal decodes the data into the given value.
func (c *Config) Unmarshal(data []byte, v any) error {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.Ctx.Clear()
	return t.unmarshal(v)
}

// UnmarshalCtx sets the [*ctx.Ctx] and decodes the data into the given value.
func (c *Config) UnmarshalCtx(data []byte, v any, ctx *ctx.Ctx) error {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.Ctx = ctx
	return t.unmarshal(v)
}

// UnmarshalScoped decodes the data into the given value with the given scope.
func (c *Config) UnmarshalScoped(data []byte, v any, operation scopes.Decoding) error {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.operation = operation
	t.Ctx.Clear()
	return t.unmarshal(v)
}

// UnmarshalScopedCtx sets the [*ctx.Ctx] and decodes the data into the given value with the given scope.
func (c *Config) UnmarshalScopedCtx(data []byte, v any, operation scopes.Decoding, ctx *ctx.Ctx) error {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.operation = operation
	t.Ctx = ctx
	return t.unmarshal(v)
}

// UnmarshalScopedWithChanges decodes the data into the given value with the given scope and returns the changes.
// The changes are the paths of the fields that have been changed.
func (c *Config) UnmarshalScopedWithChanges(data []byte, v any, operation scopes.Decoding) ([]string, error) {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.operation = operation
	t.Changes = make([]string, 0, 10)
	t.Ctx.Clear()
	err := t.unmarshal(v)
	changes := t.Changes
	t.Changes = nil
	return changes, err
}

// UnmarshalScopedWithChangesCtx sets the [*ctx.Ctx] and decodes the data into the given value with the given scope and returns the changes.
// The changes are the paths of the fields that have been changed.
func (c *Config) UnmarshalScopedWithChangesCtx(data []byte, v any, operation scopes.Decoding, ctx *ctx.Ctx) ([]string, error) {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.operation = operation
	t.Ctx = ctx
	t.Changes = make([]string, 0, 10)
	err := t.unmarshal(v)
	changes := t.Changes
	t.Changes = nil
	return changes, err

}

// UnmarshalWithChanges decodes the data into the given value and returns the changes.
// The changes are the paths of the fields that have been changed.
func (c *Config) UnmarshalWithChanges(data []byte, v any) ([]string, error) {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.Changes = make([]string, 0, 10)
	t.Ctx.Clear()
	err := t.unmarshal(v)
	changes := t.Changes
	t.Changes = nil
	return changes, err
}

// UnmarshalWithChangesCtx sets the [*ctx.Ctx] and decodes the data into the given value and returns the changes.
// The changes are the paths of the fields that have been changed.
func (c *Config) UnmarshalWithChangesCtx(data []byte, v any, ctx *ctx.Ctx) ([]string, error) {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.Ctx = ctx
	t.Changes = make([]string, 0, 10)
	err := t.unmarshal(v)
	changes := t.Changes
	t.Changes = nil
	return changes, err
}

// NewDecoder creates a new decoder with the given data.
func (c *Config) NewDecoder(data []byte) *Decoder {
	if v := c.decoderPool.Get(); v != nil {
		d := v.(*Decoder)
		d.init(data)
		return d
	}
	t := &Decoder{config: c, Ctx: &ctx.Ctx{}}
	t.init(data)
	return t
}
