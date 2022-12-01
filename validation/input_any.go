package validation

import (
	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils"
	"src.sqlkite.com/utils/typed"
)

type AnyRule interface {
	fields(fields []string) AnyRule
	Validate(fields []string, value any, object typed.Typed, input typed.Typed, res *Result) any
}

func Any() *AnyValidator {
	return &AnyValidator{}
}

type AnyValidator struct {
	field       string
	fields      []string
	dflt        any
	required    bool
	rules       []AnyRule
	errRequired InvalidField
}

func (v *AnyValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	field := v.field
	if value := args.Peek(field); value != nil {
		t[field] = utils.B2S(value)
	}
}

func (v *AnyValidator) validate(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
	value, exists := object[field]

	if !exists {
		if v.required {
			res.addField(v.errRequired)
		} else if dflt := v.dflt; dflt != 0 {
			object[field] = dflt
		}
		return
	}

	fields := v.fields
	for _, rule := range v.rules {
		value = rule.Validate(fields, value, object, input, res)
	}
	object[field] = value
}

func (v *AnyValidator) addField(field string) InputValidator {
	field, fields := expandFields(field, v.fields, v.field)

	rules := make([]AnyRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.fields(fields)
	}

	return &AnyValidator{
		field:       field,
		fields:      fields,
		dflt:        v.dflt,
		required:    v.required,
		rules:       rules,
		errRequired: invalidField(fields, Required, nil),
	}
}

func (v *AnyValidator) Required() *AnyValidator {
	v.required = true
	return v
}

func (v *AnyValidator) Default(value any) *AnyValidator {
	v.dflt = value
	return v
}

func (v *AnyValidator) Func(fn func(fields []string, value any, object typed.Typed, input typed.Typed, res *Result) any) *AnyValidator {
	v.rules = append(v.rules, AnyFunc{fn: fn})
	return v
}

type AnyFunc struct {
	fn func([]string, any, typed.Typed, typed.Typed, *Result) any
}

func (v AnyFunc) Validate(fields []string, value any, object typed.Typed, input typed.Typed, res *Result) any {
	return v.fn(fields, value, object, input, res)
}

func (r AnyFunc) fields(fields []string) AnyRule {
	return r
}
