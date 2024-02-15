package scopes

type Context int

const (
	CONTEXT_ADMIN Context = iota
	CONTEXT_CLIENT
	CONTEXT_DB
)
