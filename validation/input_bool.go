package validation

import (
	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
)

type BoolRule interface {
	fields(fields []string) BoolRule
	Validate(fields []string, value bool, object typed.Typed, input typed.Typed, res *Result) bool
}

func Bool() *BoolValidator {
	return &BoolValidator{}
}

type BoolValidator struct {
	fields      []string
	field       string
	dflt        bool
	required    bool
	rules       []BoolRule
	errType     InvalidField
	errRequired InvalidField
}

func (v *BoolValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	field := v.field
	if value := args.Peek(field); value != nil {
		// switch string([]byte) is optimized by Go
		switch string(value) {
		case "true", "TRUE", "True":
			t[field] = true
		case "false", "FALSE", "False":
			t[field] = false
		default:
			t[field] = value
		}
	}
}

func (v *BoolValidator) validate(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
	value, exists := object.BoolIf(field)

	if !exists {
		if _, exists := object[field]; !exists {
			if v.required {
				res.addField(v.errRequired)
			} else if dflt := v.dflt; dflt != false {
				object[field] = dflt
			}
			return
		}
		res.addField(v.errType)
		return
	}

	fields := v.fields
	for _, rule := range v.rules {
		value = rule.Validate(fields, value, object, input, res)
	}
	object[field] = value
}

func (v *BoolValidator) addField(field string) InputValidator {
	field, fields := expandFields(field, v.fields, v.field)

	rules := make([]BoolRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.fields(fields)
	}

	return &BoolValidator{
		field:       field,
		fields:      fields,
		dflt:        v.dflt,
		required:    v.required,
		rules:       rules,
		errType:     invalidField(fields, InvalidBoolType, nil),
		errRequired: invalidField(fields, Required, nil),
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

func (v *BoolValidator) Func(fn func(fields []string, value bool, object typed.Typed, input typed.Typed, res *Result) bool) *BoolValidator {
	v.rules = append(v.rules, BoolFunc{fn})
	return v
}

type BoolFunc struct {
	fn func([]string, bool, typed.Typed, typed.Typed, *Result) bool
}

func (v BoolFunc) Validate(fields []string, value bool, object typed.Typed, input typed.Typed, res *Result) bool {
	return v.fn(fields, value, object, input, res)
}

func (r BoolFunc) fields(fields []string) BoolRule {
	return r
}
