package decoder

import (
	"sync"

	"github.com/deveox/blaze/ctx"
	"github.com/deveox/blaze/scopes"
)

type Config struct {
	Scope       scopes.Context
	decoderPool sync.Pool
}

func (c *Config) Unmarshal(data []byte, v any) error {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.Ctx.Clear()
	return t.unmarshal(v)
}

func (c *Config) UnmarshalCtx(data []byte, v any, ctx *ctx.Ctx) error {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.Ctx = ctx
	return t.unmarshal(v)
}

func (c *Config) UnmarshalScoped(data []byte, v any, operation scopes.Decoding) error {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.operation = operation
	t.Ctx.Clear()
	return t.unmarshal(v)
}

func (c *Config) UnmarshalScopedCtx(data []byte, v any, operation scopes.Decoding, ctx *ctx.Ctx) error {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.operation = operation
	t.Ctx = ctx
	return t.unmarshal(v)
}

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
