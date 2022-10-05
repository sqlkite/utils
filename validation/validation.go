package validation

import "src.goblgobl.com/utils"

var (
	globalPool *Pool

	Required             = M(utils.VAL_REQUIRED, "required")
	InvalidStringType    = M(utils.VAL_STRING_TYPE, "must be a string")
	InvalidStringLength  = M(utils.VAL_STRING_LEN, "must be between %d and %d characters")
	InvalidStringPattern = M(utils.VAL_STRING_PATTERN, "is not valid")
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
