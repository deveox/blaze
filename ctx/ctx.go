package ctx

// Ctx is a context that can store key-value pairs.
type Ctx struct {
	*Ctx
	key string
	val any
}

// Set sets a key-value pair in the context.
func (c *Ctx) Set(key string, val any) {
	if key == "" {
		c.key = key
		c.val = val
	} else if c.Ctx == nil {
		c.Ctx = &Ctx{key: key, val: val}
	} else {
		c.Ctx.Set(key, val)
	}
}

func (c *Ctx) unset(key string) {
	if c.Ctx != nil {
		if c.Ctx.key == key {
			c.Ctx = c.Ctx.Ctx
			return
		}
		c.Ctx.Unset(key)
	}
}

// Unset removes a key-value pair from the context.
func (c *Ctx) Unset(key string) {
	if c.key == key {
		c.key = ""
	} else if c.Ctx != nil {
		c.unset(key)
	}
}

// Get returns a value by its key.
func (c *Ctx) Get(key string) (bool, any) {
	if c.key == key {
		return true, c.val
	}
	if c.Ctx == nil {
		return false, nil
	}
	return c.Ctx.Get(key)
}

// Clear removes all key-value pairs from the context.
func (c *Ctx) Clear() {
	c.Ctx = nil
}
