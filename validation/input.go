package validation

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
)

type InputValidator interface {
	validate(field string, input typed.Typed, res *Result)
	argsToTyped(field string, args *fasthttp.Args, dest typed.Typed)
}

type field struct {
	name      string
	validator InputValidator
}

func Input() *input {
	return &input{}
}

type input struct {
	fields []field
}

func (i *input) Field(name string, validator InputValidator) *input {
	i.fields = append(i.fields, field{name, validator})
	return i
}

func (i *input) Validate(input typed.Typed, res *Result) bool {
	len := res.Len()
	for _, field := range i.fields {
		field.validator.validate(field.name, input, res)
	}
	return res.Len() == len
}

func (i *input) ValidateArgs(args *fasthttp.Args, res *Result) (typed.Typed, bool) {
	fields := i.fields
	input := make(typed.Typed, len(fields))
	for _, field := range fields {
		field.validator.argsToTyped(field.name, args, input)
	}
	return input, i.Validate(input, res)
}

func invalid(meta Meta, data any, args ...any) Invalid {
	err := meta.Error
	if len(args) > 0 {
		err = fmt.Sprintf(err, args...)
	}
	return Invalid{
		Data:  data,
		Code:  meta.Code,
		Error: err,
	}
}
