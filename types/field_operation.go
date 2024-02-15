package types

import "github.com/deveox/blaze/scopes"

type Operation int

func (s Operation) CanRead() bool {
	switch s {
	case OPERATION_READ, OPERATION_READ_CREATE, OPERATION_READ_UPDATE, OPERATION_ALL:
		return true
	}
	return false
}

func (s Operation) CanWrite(scope scopes.Decoding) bool {
	switch s {
	case OPERATION_ALL, OPERATION_WRITE:
		return true
	case OPERATION_CREATE, OPERATION_READ_CREATE:
		return scope == scopes.DECODE_CREATE || scope == scopes.DECODE_ANY
	case OPERATION_UPDATE, OPERATION_READ_UPDATE:
		return scope == scopes.DECODE_UPDATE || scope == scopes.DECODE_ANY
	default:
		return false
	}
}

func (s Operation) String() string {
	switch s {
	case OPERATION_ALL:
		return "all operations"
	case OPERATION_READ:
		return "read operation"
	case OPERATION_WRITE:
		return "write operation"
	case OPERATION_CREATE:
		return "create operation"
	case OPERATION_READ_CREATE:
		return "read and create operation"
	case OPERATION_UPDATE:
		return "update operation"
	case OPERATION_READ_UPDATE:
		return "read and update operation"
	case OPERATION_IGNORE:
		return "no operation"
	default:
		return "unknown scope"
	}
}

const (
	OPERATION_ALL Operation = iota
	OPERATION_READ
	OPERATION_WRITE
	OPERATION_CREATE
	OPERATION_READ_CREATE
	OPERATION_UPDATE
	OPERATION_READ_UPDATE
	OPERATION_IGNORE
)
