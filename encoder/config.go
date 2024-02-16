package encoder

import (
	"sync"

	"github.com/deveox/blaze/scopes"
)

type Config struct {
	Scope scopes.Context
	pool  sync.Pool
}

func (c *Config) NewEncoder() *Encoder {
	if v := c.pool.Get(); v != nil {
		e := v.(*Encoder)
		e.bytes = e.bytes[:0]
		return e
	}
	e := &Encoder{bytes: make([]byte, 0, 2048), config: c}
	return e
}

func (c *Config) Marshal(v any) ([]byte, error) {
	e := c.NewEncoder()
	defer c.pool.Put(e)
	return e.marshal(v)
}
