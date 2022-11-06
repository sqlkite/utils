package validation

import (
	"regexp"

	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils/typed"
)

type StringValidator interface {
	Clone(field string) StringValidator
	Validate(value string, rest typed.Typed, res *Result) string
}

type StringConverter func(field string, value string, input typed.Typed, res *Result) any
type StringFuncValidator func(field string, value string, input typed.Typed, res *Result) string

func String(field string) *InputString {
	return &InputString{
		field:       field,
		errType:     inputError(field, InvalidStringType, nil),
		errRequired: inputError(field, Required, nil),
	}
}

type InputString struct {
	field       string
	dflt        string
	required    bool
	converter   StringConverter
	validators  []StringValidator
	errType     InvalidField
	errRequired InvalidField
}

func (i *InputString) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	field := i.field
	if value := args.Peek(field); value != nil {
		t[field] = string(value)
	}
}

func (i *InputString) Clone(field string) *InputString {
	validators := make([]StringValidator, len(i.validators))
	for i, validator := range i.validators {
		validators[i] = validator.Clone(field)
	}

	return &InputString{
		field:       field,
		dflt:        i.dflt,
		required:    i.required,
		converter:   i.converter,
		validators:  validators,
		errType:     inputError(field, InvalidStringType, nil),
		errRequired: inputError(field, Required, nil),
	}
}

func (i *InputString) validate(input typed.Typed, res *Result) {
	field := i.field
	value, exists := input.StringIf(field)

	if !exists {
		if _, exists = input[field]; !exists && i.required {
			res.add(i.errRequired)
		} else if exists {
			res.add(i.errType)
		}
		if dflt := i.dflt; dflt != "" {
			input[field] = dflt
		}
		return
	}

	for _, validator := range i.validators {
		value = validator.Validate(value, input, res)
	}

	if converter := i.converter; converter != nil {
		input[field] = converter(field, value, input, res)
	} else {
		input[field] = value
	}
}

func (i *InputString) Required() *InputString {
	i.required = true
	return i
}

func (i *InputString) Default(value string) *InputString {
	i.dflt = value
	return i
}

func (i *InputString) Length(min int, max int) *InputString {
	i.validators = append(i.validators, newStringLen(i.field, min, max))
	return i
}

func (i *InputString) Pattern(pattern string) *InputString {
	re := regexp.MustCompile(pattern)
	i.validators = append(i.validators, newStringPattern(i.field, re))
	return i
}

func (i *InputString) Func(fn StringFuncValidator) *InputString {
	i.validators = append(i.validators, newStringFunc(i.field, fn))
	return i
}

func (i *InputString) Convert(fn StringConverter) *InputString {
	i.converter = fn
	return i
}

type StringLen struct {
	min int
	max int
	err InvalidField
}

func newStringLen(field string, min int, max int) StringLen {
	return StringLen{
		min: min,
		max: max,
		err: inputError(field, InvalidStringLength, Range(min, max), min, max),
	}
}

func (v StringLen) Clone(field string) StringValidator {
	return newStringLen(field, v.min, v.max)
}

func (v StringLen) Validate(value string, rest typed.Typed, res *Result) string {
	if min := v.min; min > 0 && len(value) < min {
		res.add(v.err)
	}
	if max := v.max; max > 0 && len(value) > max {
		res.add(v.err)
	}
	return value
}

type StringPattern struct {
	pattern *regexp.Regexp
	err     InvalidField
}

func newStringPattern(field string, pattern *regexp.Regexp) StringPattern {
	return StringPattern{
		pattern: pattern,
		err:     inputError(field, InvalidStringPattern, nil),
	}
}

func (v StringPattern) Clone(field string) StringValidator {
	return newStringPattern(field, v.pattern)
}

func (v StringPattern) Validate(value string, rest typed.Typed, res *Result) string {
	if !v.pattern.MatchString(value) {
		res.add(v.err)
	}
	return value
}

type StringFunc struct {
	field string
	fn    func(string, string, typed.Typed, *Result) string
}

func newStringFunc(field string, fn StringFuncValidator) StringFunc {
	return StringFunc{
		fn:    fn,
		field: field,
	}
}
func (v StringFunc) Clone(field string) StringValidator {
	return newStringFunc(field, v.fn)
}

func (v StringFunc) Validate(value string, rest typed.Typed, res *Result) string {
	return v.fn(v.field, value, rest, res)
}
