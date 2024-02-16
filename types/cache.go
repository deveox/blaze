package types

import (
	"reflect"
	"sync"

	"github.com/deveox/gu/mirror"
)

type cache struct {
	m sync.Map
}

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
	v := NewStruct(t)
	c.m.Store(t, v)
	v.init()
	return v
}

var Cache = &cache{}
