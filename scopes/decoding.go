package scopes

type Decoding int

const (
	DECODE_ANY Decoding = iota
	DECODE_CREATE
	DECODE_UPDATE
)
