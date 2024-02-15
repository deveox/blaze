package decoder

import (
	"sync"

	"github.com/deveox/blaze/scopes"
)

type Config struct {
	ContextScope scopes.Context
	decoderPool  sync.Pool
}

func (c *Config) Unmarshal(data []byte, v any) error {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	return t.Decode(v)
}

func (c *Config) UnmarshalScoped(data []byte, v any, operation scopes.Decoding) error {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.OperationScope = operation
	return t.Decode(v)
}

func (c *Config) UnmarshalScopedWithChanges(data []byte, v any, operation scopes.Decoding) ([]string, error) {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.OperationScope = operation
	t.Changes = make([]string, 0, 10)
	err := t.Decode(v)
	changes := t.Changes
	t.Changes = nil
	return changes, err
}

func (c *Config) UnmarshalWithChanges(data []byte, v any) ([]string, error) {
	t := c.NewDecoder(data)
	defer c.decoderPool.Put(t)
	t.Changes = make([]string, 0, 10)
	err := t.Decode(v)
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
	t := &Decoder{
		ContextScope: c.ContextScope,
	}
	t.init(data)
	return t
}
