package validation

import (
	"strconv"

	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils"
	"src.sqlkite.com/utils/typed"
)

type IntFn interface {
	Validate(field string, value int, rest typed.Typed, res *Result) int
}

func Int() *IntValidator {
	return &IntValidator{
		errType:     invalid(InvalidIntType, nil),
		errRequired: invalid(Required, nil),
	}
}

type IntValidator struct {
	dflt        int
	required    bool
	validators  []IntFn
	errType     Invalid
	errRequired Invalid
}

func (i *IntValidator) argsToTyped(field string, args *fasthttp.Args, t typed.Typed) {
	if value := args.Peek(field); value != nil {
		if n, err := strconv.ParseInt(utils.B2S(value), 10, 0); err == nil {
			t[field] = n
		} else {
			t[field] = value
		}
	}
}

func (i *IntValidator) validate(field string, input typed.Typed, res *Result) {
	value, exists := input.IntIf(field)

	if !exists {
		if _, exists := input[field]; !exists {
			if i.required {
				res.add(InvalidField{i.errRequired, field})
			} else if dflt := i.dflt; dflt != 0 {
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

func (i *IntValidator) Required() *IntValidator {
	i.required = true
	return i
}

func (i *IntValidator) Default(value int) *IntValidator {
	i.dflt = value
	return i
}

func (i *IntValidator) Min(min int) *IntValidator {
	i.validators = append(i.validators, IntMin{
		min: min,
		err: invalid(InvalidIntMin, Min(min), min),
	})
	return i
}

func (i *IntValidator) Max(max int) *IntValidator {
	i.validators = append(i.validators, IntMax{
		max: max,
		err: invalid(InvalidIntMax, Max(max), max),
	})
	return i
}
func (i *IntValidator) Range(min int, max int) *IntValidator {
	i.validators = append(i.validators, IntRange{
		min: min,
		max: max,
		err: invalid(InvalidIntRange, Range(min, max), min, max),
	})
	return i
}

func (i *IntValidator) Func(fn func(field string, value int, input typed.Typed, res *Result) int) *IntValidator {
	i.validators = append(i.validators, IntFunc{fn})
	return i
}

type IntMin struct {
	min int
	err Invalid
}

func (v IntMin) Validate(field string, value int, rest typed.Typed, res *Result) int {
	if value < v.min {
		res.add(InvalidField{v.err, field})
	}
	return value
}

type IntMax struct {
	max int
	err Invalid
}

func (v IntMax) Validate(field string, value int, rest typed.Typed, res *Result) int {
	if value > v.max {
		res.add(InvalidField{v.err, field})
	}
	return value
}

type IntRange struct {
	min int
	max int
	err Invalid
}

func (v IntRange) Validate(field string, value int, rest typed.Typed, res *Result) int {
	if value < v.min || value > v.max {
		res.add(InvalidField{v.err, field})
	}
	return value
}

type IntFunc struct {
	fn func(string, int, typed.Typed, *Result) int
}

func (v IntFunc) Validate(field string, value int, rest typed.Typed, res *Result) int {
	return v.fn(field, value, rest, res)
}
