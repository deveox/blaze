package scopes

type Context int

const (
	CONTEXT_HTTP Context = iota
	CONTEXT_DB
)
