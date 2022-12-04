package validation

import (
	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
)

type BoolRule interface {
	clone() BoolRule
	Validate(field Field, value bool, object typed.Typed, input typed.Typed, res *Result) bool
}

func Bool() *BoolValidator {
	return &BoolValidator{
		errReq:  Required(),
		errType: InvalidBoolType(),
	}
}

type BoolValidator struct {
	field    Field
	dflt     bool
	required bool
	rules    []BoolRule
	errReq   Invalid
	errType  Invalid
}

func (v *BoolValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	fieldName := v.field.Name
	if value := args.Peek(fieldName); value != nil {
		// switch string([]byte) is optimized by Go
		switch string(value) {
		case "true", "TRUE", "True":
			t[fieldName] = true
		case "false", "FALSE", "False":
			t[fieldName] = false
		default:
			t[fieldName] = value
		}
	}
}

func (v *BoolValidator) validate(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
	fieldName := field.Name

	value, exists := object.BoolIf(fieldName)
	if !exists {
		if _, exists := object[fieldName]; !exists {
			if v.required {
				res.AddInvalidField(field, v.errReq)
			} else if dflt := v.dflt; dflt != false {
				object[fieldName] = dflt
			}
			return
		}
		res.AddInvalidField(field, v.errType)
		return
	}

	for _, rule := range v.rules {
		value = rule.Validate(field, value, object, input, res)
	}
	object[fieldName] = value
}

func (v *BoolValidator) addField(fieldName string) InputValidator {
	field := v.field.add(fieldName)

	rules := make([]BoolRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.clone()
	}

	return &BoolValidator{
		field:    field,
		dflt:     v.dflt,
		required: v.required,
		rules:    rules,
		errReq:   v.errReq,
		errType:  v.errType,
	}
}

func (v *BoolValidator) Required() *BoolValidator {
	v.required = true
	return v
}

func (v *BoolValidator) Default(value bool) *BoolValidator {
	v.dflt = value
	return v
}

func (v *BoolValidator) Func(fn func(field Field, value bool, object typed.Typed, input typed.Typed, res *Result) bool) *BoolValidator {
	v.rules = append(v.rules, BoolFunc{fn})
	return v
}

type BoolFunc struct {
	fn func(Field, bool, typed.Typed, typed.Typed, *Result) bool
}

func (v BoolFunc) Validate(field Field, value bool, object typed.Typed, input typed.Typed, res *Result) bool {
	return v.fn(field, value, object, input, res)
}

func (r BoolFunc) clone() BoolRule {
	return r
}
