package validation

import (
	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils"
	"src.sqlkite.com/utils/typed"
)

type AnyRule interface {
	clone() AnyRule
	Validate(field Field, value any, object typed.Typed, input typed.Typed, res *Result) any
}

func Any() *AnyValidator {
	return &AnyValidator{
		errReq: Required(),
	}
}

type AnyValidator struct {
	field    Field
	dflt     any
	required bool
	rules    []AnyRule
	errReq   Invalid
}

func (v *AnyValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	fieldName := v.field.Name
	if value := args.Peek(fieldName); value != nil {
		t[fieldName] = utils.B2S(value)
	}
}

func (v *AnyValidator) validate(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
	fieldName := field.Name

	value, exists := object[fieldName]
	if !exists {
		if v.required {
			res.AddInvalidField(field, v.errReq)
		} else if dflt := v.dflt; dflt != 0 {
			object[fieldName] = dflt
		}
		return
	}

	for _, rule := range v.rules {
		value = rule.Validate(field, value, object, input, res)
	}
	object[fieldName] = value
}

func (v *AnyValidator) addField(fieldName string) InputValidator {
	field := v.field.add(fieldName)

	rules := make([]AnyRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.clone()
	}

	return &AnyValidator{
		field:    field,
		dflt:     v.dflt,
		required: v.required,
		rules:    rules,
		errReq:   v.errReq,
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

func (v *AnyValidator) Func(fn func(field Field, value any, object typed.Typed, input typed.Typed, res *Result) any) *AnyValidator {
	v.rules = append(v.rules, AnyFunc{fn: fn})
	return v
}

type AnyFunc struct {
	fn func(Field, any, typed.Typed, typed.Typed, *Result) any
}

func (v AnyFunc) Validate(field Field, value any, object typed.Typed, input typed.Typed, res *Result) any {
	return v.fn(field, value, object, input, res)
}

func (r AnyFunc) clone() AnyRule {
	return r
}
