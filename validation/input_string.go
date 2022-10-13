package validation

import (
	"regexp"

	"src.goblgobl.com/utils/typed"
)

type StringValidator interface {
	Validate(value string, rest typed.Typed, res *Result) string
}

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
	validators  []StringValidator
	errType     InvalidField
	errRequired InvalidField
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
	input[field] = value
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
	i.validators = append(i.validators, StringLen{
		min: min,
		max: max,
		err: inputError(i.field, InvalidStringLength, Range(min, max), min, max),
	})
	return i
}

func (i *InputString) Pattern(pattern string) *InputString {
	i.validators = append(i.validators, StringPattern{
		pattern: regexp.MustCompile(pattern),
		err:     inputError(i.field, InvalidStringPattern, nil),
	})
	return i
}

func (i *InputString) Func(fn func(field string, value string, input typed.Typed, res *Result) string) *InputString {
	i.validators = append(i.validators, StringFunc{
		fn:    fn,
		field: i.field,
	})
	return i
}

type StringLen struct {
	min int
	max int
	err InvalidField
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

func (v StringFunc) Validate(value string, rest typed.Typed, res *Result) string {
	return v.fn(v.field, value, rest, res)
}
