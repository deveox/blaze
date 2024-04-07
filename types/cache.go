package types

import (
	"reflect"
	"sync"

	"github.com/deveox/gu/mirror"
)

type cache struct {
	m sync.Map
}

// Get returns a struct meta information by its type.
// If the struct is not found in the cache, it will be created and stored.
func (c *cache) Get(t reflect.Type) *Struct {
	v, ok := c.load(t)
	if ok {
		return v
	}
	return c.store(t)
}

func (c *cache) load(t reflect.Type) (value *Struct, ok bool) {
	t = mirror.DerefType(t)
	v, ok := c.m.Load(t)
	if ok {
		value = v.(*Struct)
		return
	}
	return
}

func (c *cache) store(t reflect.Type) *Struct {
	t = mirror.DerefType(t)
	v := newStruct(t)
	c.m.Store(t, v)
	v.init()
	return v
}

// Cache is a global cache of struct meta information.
var Cache = &cache{}
