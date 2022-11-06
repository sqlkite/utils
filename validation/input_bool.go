package validation

import (
	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/typed"
)

type BoolValidator interface {
	Validate(value bool, rest typed.Typed, res *Result) bool
}

func Bool(field string) *InputBool {
	return &InputBool{
		field:       field,
		errType:     inputError(field, InvalidBoolType, nil),
		errRequired: inputError(field, Required, nil),
	}
}

type InputBool struct {
	dflt        bool
	field       string
	coerce      bool
	required    bool
	validators  []BoolValidator
	errType     InvalidField
	errRequired InvalidField
}

func (i *InputBool) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	field := i.field
	if value := args.Peek(field); value != nil {
		switch utils.B2S(value) {
		case "true", "TRUE", "True":
			t[field] = true
		case "false", "FALSE", "False":
			t[field] = false
		default:
			t[field] = value
		}
	}
}

func (i *InputBool) Required() *InputBool {
	i.required = true
	return i
}

func (i *InputBool) Coerce() *InputBool {
	i.coerce = true
	return i
}

func (i *InputBool) validate(input typed.Typed, res *Result) {
	field := i.field
	value, exists := input.BoolIf(field)

	if !exists {
		if _, exists := input[field]; !exists {
			if i.required {
				res.add(i.errRequired)
			} else if dflt := i.dflt; dflt != false {
				input[field] = dflt
			}
			return
		}
		res.add(i.errType)
		return
	}

	for _, validator := range i.validators {
		value = validator.Validate(value, input, res)
	}
	input[field] = value
}

func (i *InputBool) Default(value bool) *InputBool {
	i.dflt = value
	return i
}

func (i *InputBool) Func(fn func(field string, value bool, input typed.Typed, res *Result) bool) *InputBool {
	i.validators = append(i.validators, BoolFunc{
		fn:    fn,
		field: i.field,
	})
	return i
}

type BoolFunc struct {
	field string
	fn    func(string, bool, typed.Typed, *Result) bool
}

func (v BoolFunc) Validate(value bool, rest typed.Typed, res *Result) bool {
	return v.fn(v.field, value, rest, res)
}
