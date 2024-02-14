package encoder

import (
	"sync"

	"github.com/deveox/blaze/scopes"
)

type Config struct {
	ContextScope scopes.Context
	UserScope    scopes.User
	Operation    scopes.Operation
	pool         sync.Pool
}

func (c *Config) NewEncoder() *Encoder {
	if v := c.pool.Get(); v != nil {
		e := v.(*Encoder)
		e.bytes = e.bytes[:0]
		return e
	}
	e := &Encoder{bytes: make([]byte, 0, 2048)}
	e.ContextScope = c.ContextScope
	e.UserScope = c.UserScope
	return e
}

func (c *Config) Marshal(v any) ([]byte, error) {
	e := c.NewEncoder()
	defer c.pool.Put(e)
	return e.Encode(v)
}

func (c *Config) MarshalOperation(v any, operation scopes.Operation) ([]byte, error) {
	t := c.NewEncoder()
	defer c.pool.Put(t)
	t.OperationScope = operation
	return t.Encode(v)
}
