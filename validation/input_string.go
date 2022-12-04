package validation

import (
	"regexp"

	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
)

type StringRule interface {
	clone() StringRule
	Validate(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string
}

type StringConverter func(field Field, value string, object typed.Typed, input typed.Typed, res *Result) any
type StringFuncValidator func(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string

func String() *StringValidator {
	return &StringValidator{
		errReq:  Required(),
		errType: InvalidStringType(),
	}
}

type StringValidator struct {
	field     Field
	dflt      string
	required  bool
	converter StringConverter
	rules     []StringRule
	errReq    Invalid
	errType   Invalid
}

func (v *StringValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	fieldName := v.field.Name
	if value := args.Peek(fieldName); value != nil {
		t[fieldName] = string(value)
	}
}

func (v *StringValidator) validate(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
	fieldName := field.Name

	value, exists := object.StringIf(fieldName)
	if !exists {
		if _, exists = object[fieldName]; !exists && v.required {
			res.AddInvalidField(field, v.errReq)
		} else if exists {
			res.AddInvalidField(field, v.errType)
		}
		if dflt := v.dflt; dflt != "" {
			object[fieldName] = dflt
		}
		return
	}

	for _, rule := range v.rules {
		value = rule.Validate(field, value, object, input, res)
	}

	if converter := v.converter; converter != nil {
		object[fieldName] = converter(field, value, object, input, res)
	} else {
		object[fieldName] = value
	}
}

func (v *StringValidator) addField(fieldName string) InputValidator {
	field := v.field.add(fieldName)

	rules := make([]StringRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.clone()
	}

	return &StringValidator{
		field:     field,
		dflt:      v.dflt,
		required:  v.required,
		converter: v.converter,
		rules:     rules,
		errReq:    v.errReq,
		errType:   v.errType,
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
	v.rules = append(v.rules, StringChoice{
		valid: valid,
		err:   InvalidStringChoice(valid),
	})
	return v
}

func (v *StringValidator) Length(min int, max int) *StringValidator {
	v.rules = append(v.rules, StringLen{
		min: min,
		max: max,
		err: InvalidStringLength(min, max),
	})
	return v
}

func (v *StringValidator) Pattern(pattern string, errorMessage ...string) *StringValidator {
	v.rules = append(v.rules, StringPattern{
		pattern: regexp.MustCompile(pattern),
		err:     InvalidStringPattern(errorMessage...),
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
	err Invalid
}

func (r StringLen) Validate(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string {
	if min := r.min; min > 0 && len(value) < min {
		res.AddInvalidField(field, r.err)
	}
	if max := r.max; max > 0 && len(value) > max {
		res.AddInvalidField(field, r.err)
	}
	return value
}

func (r StringLen) clone() StringRule {
	return StringLen{
		min: r.min,
		max: r.max,
		err: r.err,
	}
}

type StringPattern struct {
	pattern *regexp.Regexp
	err     Invalid
}

func (r StringPattern) Validate(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string {
	if !r.pattern.MatchString(value) {
		res.AddInvalidField(field, r.err)
	}
	return value
}

func (r StringPattern) clone() StringRule {
	return StringPattern{
		pattern: r.pattern,
		err:     r.err,
	}
}

type StringChoice struct {
	valid []string
	err   Invalid
}

func (r StringChoice) Validate(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string {
	for _, valid := range r.valid {
		if value == valid {
			return value
		}
	}
	res.AddInvalidField(field, r.err)
	return value
}

func (r StringChoice) clone() StringRule {
	return StringChoice{
		valid: r.valid,
		err:   r.err,
	}
}

type StringFunc struct {
	fn StringFuncValidator
}

func (r StringFunc) Validate(field Field, value string, object typed.Typed, input typed.Typed, res *Result) string {
	return r.fn(field, value, object, input, res)
}

func (r StringFunc) clone() StringRule {
	return r
}
