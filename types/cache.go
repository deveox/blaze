package types

import (
	"reflect"
	"sync"

	"github.com/deveox/gu/mirror"
)

type cache struct {
	m sync.Map
}

func (c *cache) Get(t reflect.Type) (*Struct, error) {
	v, ok := c.load(t)
	if ok {
		return v, nil
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

func (c *cache) store(t reflect.Type) (*Struct, error) {
	t = mirror.DerefType(t)
	v := NewStruct(t)
	c.m.Store(t, v)
	err := v.init()
	return v, err
}

var Cache = &cache{}
