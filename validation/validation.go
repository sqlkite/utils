package validation

import "src.sqlkite.com/utils"

var (
	globalPool *Pool

	Required             = M(utils.VAL_REQUIRED, "required")
	InvalidStringType    = M(utils.VAL_STRING_TYPE, "must be a string")
	InvalidStringLength  = M(utils.VAL_STRING_LEN, "must be between %d and %d characters")
	InvalidStringPattern = M(utils.VAL_STRING_PATTERN, "is not valid")
	InvalidIntType       = M(utils.VAL_INT_TYPE, "must be a number")
	InvalidIntMin        = M(utils.VAL_INT_MIN, "must be greater or equal to %d")
	InvalidIntMax        = M(utils.VAL_INT_MAX, "must be less than or equal to %d")
	InvalidIntRange      = M(utils.VAL_INT_RANGE, "must be between %d and %d")
	InvalidBoolType      = M(utils.VAL_BOOL_TYPE, "must be true or false")
	InvalidUUIDType      = M(utils.VAL_UUID_TYPE, "must be a uuid")
)

func Checkout() *Result {
	return globalPool.Checkout()
}

type Meta struct {
	Code  uint32
	Error string
}

func M(code uint32, error string) Meta {
	return Meta{
		Code:  code,
		Error: error,
	}
}

type Invalid struct {
	Code  uint32 `json:"code"`
	Error string `json:"error"`
	Data  any    `json:"data,omitempty"`
}

type InvalidField struct {
	Invalid
	Field string `json:"field"`
}

type DataRange struct {
	Min any `json:"min"`
	Max any `json:"max"`
}

func Range(min any, max any) any {
	return DataRange{
		Min: min,
		Max: max,
	}
}

type DataMin struct {
	Min any `json:"min"`
}

func Min(min any) any {
	return DataMin{
		Min: min,
	}
}

type DataMax struct {
	Max any `json:"max"`
}

func Max(max any) any {
	return DataMax{
		Max: max,
	}
}
