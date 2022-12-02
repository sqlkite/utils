package validation

import (
	"strconv"

	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils"
	"src.sqlkite.com/utils/typed"
)

type FloatRule interface {
	fields(fields []string) FloatRule
	Validate(fields []string, value float64, object typed.Typed, input typed.Typed, res *Result) float64
}

func Float() *FloatValidator {
	return &FloatValidator{}
}

type FloatValidator struct {
	field       string
	fields      []string
	dflt        float64
	required    bool
	rules       []FloatRule
	errType     InvalidField
	errRequired InvalidField
}

func (v *FloatValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	field := v.field
	if value := args.Peek(field); value != nil {
		if n, err := strconv.ParseFloat(utils.B2S(value), 64); err == nil {
			t[field] = n
		} else {
			t[field] = value
		}
	}
}

func (v *FloatValidator) validate(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
	value, exists := object.FloatIf(field)

	if !exists {
		if _, exists := object[field]; !exists {
			if v.required {
				res.addField(v.errRequired)
			} else if dflt := v.dflt; dflt != 0 {
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

func (v *FloatValidator) addField(field string) InputValidator {
	field, fields := expandFields(field, v.fields, v.field)

	rules := make([]FloatRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.fields(fields)
	}

	return &FloatValidator{
		field:       field,
		fields:      fields,
		dflt:        v.dflt,
		required:    v.required,
		rules:       rules,
		errType:     invalidField(fields, InvalidFloatType, nil),
		errRequired: invalidField(fields, Required, nil),
	}
}

func (v *FloatValidator) Required() *FloatValidator {
	v.required = true
	return v
}

func (v *FloatValidator) Default(value float64) *FloatValidator {
	v.dflt = value
	return v
}

func (v *FloatValidator) Min(min float64) *FloatValidator {
	v.rules = append(v.rules, FloatMin{min: min})
	return v
}

func (v *FloatValidator) Max(max float64) *FloatValidator {
	v.rules = append(v.rules, FloatMax{max: max})
	return v
}

func (v *FloatValidator) Range(min float64, max float64) *FloatValidator {
	v.rules = append(v.rules, FloatRange{
		min: min,
		max: max,
	})
	return v
}

func (v *FloatValidator) Func(fn func(fields []string, value float64, object typed.Typed, input typed.Typed, res *Result) float64) *FloatValidator {
	v.rules = append(v.rules, FloatFunc{fn: fn})
	return v
}

type FloatMin struct {
	min float64
	err InvalidField
}

func (r FloatMin) Validate(fields []string, value float64, object typed.Typed, input typed.Typed, res *Result) float64 {
	if value < r.min {
		res.addField(r.err)
	}
	return value
}

func (r FloatMin) fields(fields []string) FloatRule {
	min := r.min
	return FloatMin{
		min: min,
		err: invalidField(fields, InvalidFloatMin, Min(min), min),
	}
}

type FloatMax struct {
	max float64
	err InvalidField
}

func (r FloatMax) Validate(fields []string, value float64, object typed.Typed, input typed.Typed, res *Result) float64 {
	if value > r.max {
		res.addField(r.err)
	}
	return value
}

func (r FloatMax) fields(fields []string) FloatRule {
	max := r.max
	return FloatMax{
		max: max,
		err: invalidField(fields, InvalidFloatMax, Max(max), max),
	}
}

type FloatRange struct {
	min float64
	max float64
	err InvalidField
}

func (r FloatRange) Validate(fields []string, value float64, object typed.Typed, input typed.Typed, res *Result) float64 {
	if value < r.min || value > r.max {
		res.addField(r.err)
	}
	return value
}

func (r FloatRange) fields(fields []string) FloatRule {
	min := r.min
	max := r.max
	return FloatRange{
		min: min,
		max: max,
		err: invalidField(fields, InvalidFloatRange, Range(min, max), min, max),
	}
}

type FloatFunc struct {
	fn func([]string, float64, typed.Typed, typed.Typed, *Result) float64
}

func (v FloatFunc) Validate(fields []string, value float64, object typed.Typed, input typed.Typed, res *Result) float64 {
	return v.fn(fields, value, object, input, res)
}

func (r FloatFunc) fields(fields []string) FloatRule {
	return r
}
