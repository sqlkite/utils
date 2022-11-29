package validation

import (
	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
)

type BoolFn interface {
	Validate(field string, value bool, rest typed.Typed, res *Result) bool
}

func Bool() *BoolValidator {
	return &BoolValidator{
		errType:     invalid(InvalidBoolType, nil),
		errRequired: invalid(Required, nil),
	}
}

type BoolValidator struct {
	dflt        bool
	required    bool
	validators  []BoolFn
	errType     Invalid
	errRequired Invalid
}

func (i *BoolValidator) argsToTyped(field string, args *fasthttp.Args, t typed.Typed) {
	if value := args.Peek(field); value != nil {
		// switch string([]byte) is optimized by Go
		switch string(value) {
		case "true", "TRUE", "True":
			t[field] = true
		case "false", "FALSE", "False":
			t[field] = false
		default:
			t[field] = value
		}
	}
}

func (i *BoolValidator) Required() *BoolValidator {
	i.required = true
	return i
}

func (i *BoolValidator) validate(field string, input typed.Typed, res *Result) {
	value, exists := input.BoolIf(field)

	if !exists {
		if _, exists := input[field]; !exists {
			if i.required {
				res.add(InvalidField{i.errRequired, field})
			} else if dflt := i.dflt; dflt != false {
				input[field] = dflt
			}
			return
		}
		res.add(InvalidField{i.errType, field})
		return
	}

	for _, validator := range i.validators {
		value = validator.Validate(field, value, input, res)
	}
	input[field] = value
}

func (i *BoolValidator) Default(value bool) *BoolValidator {
	i.dflt = value
	return i
}

func (i *BoolValidator) Func(fn func(field string, value bool, input typed.Typed, res *Result) bool) *BoolValidator {
	i.validators = append(i.validators, BoolFunc{fn})
	return i
}

type BoolFunc struct {
	fn func(string, bool, typed.Typed, *Result) bool
}

func (v BoolFunc) Validate(field string, value bool, rest typed.Typed, res *Result) bool {
	return v.fn(field, value, rest, res)
}
