package validation

import (
	"regexp"

	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
)

type StringFn interface {
	Validate(field string, value string, rest typed.Typed, res *Result) string
}

type StringConverter func(field string, value string, input typed.Typed, res *Result) any
type StringFuncValidator func(field string, value string, input typed.Typed, res *Result) string

func String() *StringValidator {
	return &StringValidator{
		errType:     invalid(InvalidStringType, nil),
		errRequired: invalid(Required, nil),
	}
}

type StringValidator struct {
	dflt        string
	required    bool
	converter   StringConverter
	validators  []StringFn
	errType     Invalid
	errRequired Invalid
}

func (i *StringValidator) argsToTyped(field string, args *fasthttp.Args, t typed.Typed) {
	if value := args.Peek(field); value != nil {
		t[field] = string(value)
	}
}

func (i *StringValidator) validate(field string, input typed.Typed, res *Result) {
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

	for _, validator := range i.validators {
		value = validator.Validate(field, value, input, res)
	}

	if converter := i.converter; converter != nil {
		input[field] = converter(field, value, input, res)
	} else {
		input[field] = value
	}
}

func (i *StringValidator) Required() *StringValidator {
	i.required = true
	return i
}

func (i *StringValidator) Default(value string) *StringValidator {
	i.dflt = value
	return i
}

func (i *StringValidator) Length(min int, max int) *StringValidator {
	i.validators = append(i.validators, StringLen{
		min: min,
		max: max,
		err: invalid(InvalidStringLength, Range(min, max), min, max),
	})
	return i
}

func (i *StringValidator) Pattern(pattern string) *StringValidator {
	re := regexp.MustCompile(pattern)
	i.validators = append(i.validators, StringPattern{
		pattern: re,
		err:     invalid(InvalidStringPattern, nil),
	})
	return i
}

func (i *StringValidator) Func(fn StringFuncValidator) *StringValidator {
	i.validators = append(i.validators, StringFunc{fn})
	return i
}

func (i *StringValidator) Convert(fn StringConverter) *StringValidator {
	i.converter = fn
	return i
}

type StringLen struct {
	min int
	max int
	err Invalid
}

func (v StringLen) Validate(field string, value string, rest typed.Typed, res *Result) string {
	if min := v.min; min > 0 && len(value) < min {
		res.add(InvalidField{v.err, field})
	}
	if max := v.max; max > 0 && len(value) > max {
		res.add(InvalidField{v.err, field})
	}
	return value
}

type StringPattern struct {
	pattern *regexp.Regexp
	err     Invalid
}

func (v StringPattern) Validate(field string, value string, rest typed.Typed, res *Result) string {
	if !v.pattern.MatchString(value) {
		res.add(InvalidField{v.err, field})
	}
	return value
}

type StringFunc struct {
	fn func(string, string, typed.Typed, *Result) string
}

func (v StringFunc) Validate(field string, value string, rest typed.Typed, res *Result) string {
	return v.fn(field, value, rest, res)
}
