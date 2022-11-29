package validation

import (
	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
	"src.sqlkite.com/utils/uuid"
)

func UUID() *UUIDValidator {
	return &UUIDValidator{
		errType:     invalid(InvalidUUIDType, nil),
		errRequired: invalid(Required, nil),
	}
}

type UUIDValidator struct {
	dflt        string
	required    bool
	errType     Invalid
	errRequired Invalid
}

func (i *UUIDValidator) argsToTyped(field string, args *fasthttp.Args, t typed.Typed) {
	if value := args.Peek(field); value != nil {
		t[field] = string(value)
	}
}

func (i *UUIDValidator) validate(field string, input typed.Typed, res *Result) {
	value, exists := input.StringIf(field)

	if !exists {
		if _, exists = input[field]; !exists && i.required {
			res.add(InvalidField{i.errRequired, field})
		} else if exists {
			res.add(InvalidField{i.errType, field})
		}
		if dflt := i.dflt; dflt != "" {
			input[field] = dflt
		}
		return
	}

	if !uuid.IsValid(value) {
		res.add(InvalidField{i.errType, field})
	}
}

func (i *UUIDValidator) Required() *UUIDValidator {
	i.required = true
	return i
}

func (i *UUIDValidator) Default(value string) *UUIDValidator {
	i.dflt = value
	return i
}
