package validation

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
)

type InputValidator interface {
	validate(input typed.Typed, res *Result)
	argsToTyped(args *fasthttp.Args, dest typed.Typed)
}

func Input() *input {
	return &input{}
}

type input struct {
	validators []InputValidator
}

func (i *input) Field(v InputValidator) *input {
	i.validators = append(i.validators, v)
	return i
}

func (i *input) Validate(input typed.Typed, res *Result) bool {
	len := res.Len()
	for _, validator := range i.validators {
		validator.validate(input, res)
	}
	return res.Len() == len
}

func (i *input) ValidateArgs(args *fasthttp.Args, res *Result) (typed.Typed, bool) {
	validators := i.validators
	input := make(typed.Typed, len(validators))
	for _, validator := range i.validators {
		validator.argsToTyped(args, input)
	}
	return input, i.Validate(input, res)
}

func inputError(field string, meta Meta, data any, args ...any) InvalidField {
	err := meta.Error
	if len(args) > 0 {
		err = fmt.Sprintf(err, args...)
	}
	return InvalidField{
		Field: field,
		Invalid: Invalid{
			Data:  data,
			Code:  meta.Code,
			Error: err,
		},
	}
}
