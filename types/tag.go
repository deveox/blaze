package types

const (
	TAG                         = "json"
	TAG_BLAZE                   = "blaze"
	TAG_SCOPE_CLIENT            = "client"
	TAG_SCOPE_ADMIN             = "admin"
	TAG_SCOPE_DB                = "db"
	TAG_SCOPE_HTTP              = "http"
	TV_SCOPE_READ               = "read"
	TV_SCOPE_WRITE              = "write"
	TV_SCOPE_CREATE             = "create"
	TV_SCOPE_UPDATE             = "update"
	TV_SCOPE_READ_CREATE        = "read.create"
	TV_SCOPE_READ_UPDATE        = "read.update"
	TV_SCOPE_READ_WRITE         = "read.write"
	TV_SCOPE_CREATE_UPDATE      = "create.update"
	TV_SCOPE_READ_CREATE_UPDATE = "read.create.update"
)

func tagPartToOperation(s string) Operation {
	switch s {
	case TV_SCOPE_CREATE:
		return OPERATION_CREATE
	case TV_SCOPE_READ:
		return OPERATION_READ
	case TV_SCOPE_WRITE:
		return OPERATION_WRITE
	case TV_SCOPE_UPDATE:
		return OPERATION_UPDATE
	case TV_SCOPE_READ_CREATE:
		return OPERATION_READ_CREATE
	case TV_SCOPE_READ_UPDATE:
		return OPERATION_READ_UPDATE
	case TV_SCOPE_READ_WRITE:
		return OPERATION_ALL
	case TV_SCOPE_CREATE_UPDATE:
		return OPERATION_WRITE
	case TV_SCOPE_READ_CREATE_UPDATE:
		return OPERATION_ALL
	default:
		return OPERATION_ALL
	}
}
