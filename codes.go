package utils

const (
	VAL_REQUIRED       = 1001
	VAL_STRING_TYPE    = 1002
	VAL_STRING_LEN     = 1003
	VAL_STRING_PATTERN = 1004
	VAL_INT_TYPE       = 1005
	VAL_INT_MIN        = 1006
	VAL_INT_MAX        = 1007
	VAL_INT_RANGE      = 1008
	VAL_BOOL_TYPE      = 1009
	VAL_UUID_TYPE      = 1010
	VAL_ARRAY_TYPE     = 1011

	RES_SERVER_ERROR         = 2001
	RES_SERIALIZATION_ERROR  = 2002
	RES_INVALID_JSON_PAYLOAD = 2003
	RES_VALIDATION           = 2004

	ERR_INVALID_LOG_LEVEL  = 3001
	ERR_INVALID_LOG_FORMAT = 3002
	ERR_PG_INIT            = 3003
	ERR_SQLITE_INIT        = 3004
)
