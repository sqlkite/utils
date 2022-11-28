package validation

import (
	"strconv"

	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils"
	"src.sqlkite.com/utils/typed"
)

type IntValidator interface {
	Validate(value int, rest typed.Typed, res *Result) int
}

func Int(field string) *InputInt {
	return &InputInt{
		field:       field,
		errType:     inputError(field, InvalidIntType, nil),
		errRequired: inputError(field, Required, nil),
	}
}

type InputInt struct {
	dflt        int
	field       string
	required    bool
	validators  []IntValidator
	errType     InvalidField
	errRequired InvalidField
}

func (i *InputInt) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	field := i.field
	if value := args.Peek(field); value != nil {
		if n, err := strconv.ParseInt(utils.B2S(value), 10, 0); err == nil {
			t[field] = n
		} else {
			t[field] = value
		}
	}
}

func (i *InputInt) validate(input typed.Typed, res *Result) {
	field := i.field
	value, exists := input.IntIf(field)

	if !exists {
		if _, exists := input[field]; !exists {
			if i.required {
				res.add(i.errRequired)
			} else if dflt := i.dflt; dflt != 0 {
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

func (i *InputInt) Required() *InputInt {
	i.required = true
	return i
}

func (i *InputInt) Default(value int) *InputInt {
	i.dflt = value
	return i
}

func (i *InputInt) Min(min int) *InputInt {
	i.validators = append(i.validators, IntMin{
		min: min,
		err: inputError(i.field, InvalidIntMin, Min(min), min),
	})
	return i
}

func (i *InputInt) Max(max int) *InputInt {
	i.validators = append(i.validators, IntMax{
		max: max,
		err: inputError(i.field, InvalidIntMax, Max(max), max),
	})
	return i
}
func (i *InputInt) Range(min int, max int) *InputInt {
	i.validators = append(i.validators, IntRange{
		min: min,
		max: max,
		err: inputError(i.field, InvalidIntRange, Range(min, max), min, max),
	})
	return i
}

func (i *InputInt) Func(fn func(field string, value int, input typed.Typed, res *Result) int) *InputInt {
	i.validators = append(i.validators, IntFunc{
		fn:    fn,
		field: i.field,
	})
	return i
}

type IntMin struct {
	min int
	err InvalidField
}

func (v IntMin) Validate(value int, rest typed.Typed, res *Result) int {
	if value < v.min {
		res.add(v.err)
	}
	return value
}

type IntMax struct {
	max int
	err InvalidField
}

func (v IntMax) Validate(value int, rest typed.Typed, res *Result) int {
	if value > v.max {
		res.add(v.err)
	}
	return value
}

type IntRange struct {
	min int
	max int
	err InvalidField
}

func (v IntRange) Validate(value int, rest typed.Typed, res *Result) int {
	if value < v.min || value > v.max {
		res.add(v.err)
	}
	return value
}

type IntFunc struct {
	field string
	fn    func(string, int, typed.Typed, *Result) int
}

func (v IntFunc) Validate(value int, rest typed.Typed, res *Result) int {
	return v.fn(v.field, value, rest, res)
}
