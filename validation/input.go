package validation

import (
	"fmt"

	"src.goblgobl.com/utils/typed"
)

type InputValidator interface {
	validate(input typed.Typed, res *Result)
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
