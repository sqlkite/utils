package validation

import (
	"regexp"

	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
)

type StringRule interface {
	fields(fields []string) StringRule
	Validate(fields []string, value string, object typed.Typed, input typed.Typed, res *Result) string
}

type StringConverter func(fields []string, value string, object typed.Typed, input typed.Typed, res *Result) any
type StringFuncValidator func(fields []string, value string, object typed.Typed, input typed.Typed, res *Result) string

func String() *StringValidator {
	return &StringValidator{}
}

type StringValidator struct {
	field       string
	fields      []string
	dflt        string
	required    bool
	converter   StringConverter
	rules       []StringRule
	errType     InvalidField
	errRequired InvalidField
}

func (v *StringValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	field := v.field
	if value := args.Peek(field); value != nil {
		t[field] = string(value)
	}
}

func (v *StringValidator) validate(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
	value, exists := object.StringIf(field)

	if !exists {
		if _, exists = object[field]; !exists && v.required {
			res.addField(v.errRequired)
		} else if exists {
			res.addField(v.errType)
		}
		if dflt := v.dflt; dflt != "" {
			object[field] = dflt
		}
		return
	}

	fields := v.fields
	for _, rule := range v.rules {
		value = rule.Validate(fields, value, object, input, res)
	}

	if converter := v.converter; converter != nil {
		object[field] = converter(fields, value, object, input, res)
	} else {
		object[field] = value
	}
}

func (v *StringValidator) addField(field string) InputValidator {
	field, fields := expandFields(field, v.fields, v.field)

	rules := make([]StringRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.fields(fields)
	}

	return &StringValidator{
		field:       field,
		fields:      fields,
		dflt:        v.dflt,
		required:    v.required,
		converter:   v.converter,
		rules:       rules,
		errType:     invalidField(fields, InvalidStringType, nil),
		errRequired: invalidField(fields, Required, nil),
	}
}

func (v *StringValidator) Required() *StringValidator {
	v.required = true
	return v
}

func (v *StringValidator) Default(value string) *StringValidator {
	v.dflt = value
	return v
}

func (v *StringValidator) Choice(valid ...string) *StringValidator {
	v.rules = append(v.rules, StringChoice{valid: valid})
	return v
}

func (v *StringValidator) Length(min int, max int) *StringValidator {
	v.rules = append(v.rules, StringLen{
		min: min,
		max: max,
	})
	return v
}

func (v *StringValidator) Pattern(pattern string, error ...string) *StringValidator {
	errorMessage := ""
	if error != nil {
		errorMessage = error[0]
	}
	v.rules = append(v.rules, StringPattern{
		errorMessage: errorMessage,
		pattern:      regexp.MustCompile(pattern),
	})
	return v
}

func (v *StringValidator) Func(fn StringFuncValidator) *StringValidator {
	v.rules = append(v.rules, StringFunc{fn})
	return v
}

func (v *StringValidator) Convert(fn StringConverter) *StringValidator {
	v.converter = fn
	return v
}

type StringLen struct {
	min int
	max int
	err InvalidField
}

func (r StringLen) Validate(fields []string, value string, object typed.Typed, input typed.Typed, res *Result) string {
	if min := r.min; min > 0 && len(value) < min {
		res.addField(r.err)
	}
	if max := r.max; max > 0 && len(value) > max {
		res.addField(r.err)
	}
	return value
}

func (r StringLen) fields(fields []string) StringRule {
	min := r.min
	max := r.max
	return StringLen{
		min: min,
		max: max,
		err: invalidField(fields, InvalidStringLength, Range(min, max), min, max),
	}
}

type StringPattern struct {
	pattern      *regexp.Regexp
	err          InvalidField
	errorMessage string
}

func (r StringPattern) Validate(fields []string, value string, object typed.Typed, input typed.Typed, res *Result) string {
	if !r.pattern.MatchString(value) {
		res.addField(r.err)
	}
	return value
}

func (r StringPattern) fields(fields []string) StringRule {
	pattern := r.pattern
	errorMessage := r.errorMessage
	sp := StringPattern{
		pattern:      pattern,
		errorMessage: errorMessage,
		err:          invalidField(fields, InvalidStringPattern, nil),
	}
	if errorMessage != "" {
		sp.err.Error = errorMessage
	}
	return sp
}

type StringChoice struct {
	valid []string
	err   InvalidField
}

func (r StringChoice) Validate(fields []string, value string, object typed.Typed, input typed.Typed, res *Result) string {
	for _, valid := range r.valid {
		if value == valid {
			return value
		}
	}
	res.addField(r.err)
	return value
}

func (r StringChoice) fields(fields []string) StringRule {
	valid := r.valid
	return StringChoice{
		valid: valid,
		err:   invalidField(fields, InvalidStringChoice, Choice(valid)),
	}
}

type StringFunc struct {
	fn StringFuncValidator
}

func (r StringFunc) Validate(fields []string, value string, object typed.Typed, input typed.Typed, res *Result) string {
	return r.fn(fields, value, object, input, res)
}

func (r StringFunc) fields(fields []string) StringRule {
	return r
}
