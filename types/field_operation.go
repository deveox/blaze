package types

type Operation int

func (s Operation) String() string {
	switch s {
	case OPERATION_ALL:
		return "all scopes"
	case OPERATION_READ:
		return "read scope"
	case OPERATION_WRITE:
		return "write scope"
	case OPERATION_CREATE:
		return "create scope"
	case OPERATION_READ_CREATE:
		return "read and create scope"
	case OPERATION_UPDATE:
		return "update scope"
	case OPERATION_READ_UPDATE:
		return "read and update scope"
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
)
