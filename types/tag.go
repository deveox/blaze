package types

import (
	"strings"
)

const (
	TAG_NAME_JSON    = "json"
	TAG_NAME_BLAZE   = "blaze"
	TAG_SCOPE_CLIENT = "client"
	TAG_SCOPE_ADMIN  = "admin"
	TAG_KEEP         = "keep"
	TAG_NO_DB        = "no-db"
	TAG_NO_HTTP      = "no-http"
	TAG_SV_IGNORE    = "-"
	TAG_SV_READ      = "read"
	TAG_SV_WRITE     = "write"
	TAG_SV_CREATE    = "create"
	TAG_SV_UPDATE    = "update"
	TAG_SV_ALL       = "all"
)

func tagPartToOperation(s string) Operation {
	read := false
	update := false
	create := false

	var v string
	for {
		v, s, _ = strings.Cut(s, ".")
		switch v {
		case TAG_SV_READ:
			read = true
		case TAG_SV_WRITE:
			update = true
			create = true
		case TAG_SV_CREATE:
			create = true
		case TAG_SV_UPDATE:
			update = true
		case TAG_SV_ALL:
			return OPERATION_ALL
		case TAG_SV_IGNORE:
			return OPERATION_IGNORE
		}
		if s == "" {
			break
		}
	}
	if read && update && create {
		return OPERATION_ALL
	}
	if read && update {
		return OPERATION_READ_UPDATE
	}
	if read && create {
		return OPERATION_READ_CREATE
	}
	if read {
		return OPERATION_READ
	}
	if update && create {
		return OPERATION_WRITE
	}
	if update {
		return OPERATION_UPDATE
	}
	if create {
		return OPERATION_CREATE
	}
	return OPERATION_IGNORE
}
