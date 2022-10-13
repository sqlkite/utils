package validation

import (
	"src.goblgobl.com/utils/typed"
)

type BoolValidator interface {
	Validate(value bool, rest typed.Typed, res *Result) bool
}

func Bool(field string, required bool) *InputBool {
	return &InputBool{
		field:       field,
		required:    required,
		errType:     inputError(field, InvalidBoolType, nil),
		errRequired: inputError(field, Required, nil),
	}
}

type InputBool struct {
	dflt        bool
	field       string
	required    bool
	validators  []BoolValidator
	errType     InvalidField
	errRequired InvalidField
}

func (i *InputBool) validate(input typed.Typed, res *Result) {
	field := i.field
	value, exists := input.BoolIf(field)

	if !exists {
		if _, exists = input[field]; !exists && i.required {
			res.add(i.errRequired)
		} else if exists {
			res.add(i.errType)
		}
		if dflt := i.dflt; dflt != false {
			input[field] = dflt
		}
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
