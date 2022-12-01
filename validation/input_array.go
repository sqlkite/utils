package validation

import (
	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
)

type ArrayRule interface {
	fields(fields []string) ArrayRule
	Validate(fields []string, value []typed.Typed, object typed.Typed, input typed.Typed, res *Result) []typed.Typed
}

func Array() *ArrayValidator {
	return &ArrayValidator{}
}

type ArrayValidator struct {
	field       string
	fields      []string
	required    bool
	dflt        []typed.Typed
	rules       []ArrayRule
	validator   *ObjectValidator
	errType     InvalidField
	errRequired InvalidField
}

func (v *ArrayValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	panic("ArrayValidator.argstoType not supported")
}

func (v *ArrayValidator) validate(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
	values, exists := object.ObjectsIf(v.field)

	if !exists {
		if _, exists := object[field]; !exists {
			if v.required {
				res.addField(v.errRequired)
			} else if dflt := v.dflt; dflt != nil {
				object[field] = dflt
			}
			return
		}
		res.addField(v.errType)
		return
	}

	// first we apply validation on the array itself (e.g. min length)
	fields := v.fields
	for _, rule := range v.rules {
		values = rule.Validate(fields, values, object, input, res)
	}

	// next we apply validation on every item within the array
	res.beginArray()
	validator := v.validator
	for i, value := range values {
		res.arrayIndex(i)
		validator.Validate(value, res)
	}
	res.endArray()
}

func (v *ArrayValidator) addField(field string) InputValidator {
	validator := v.validator.addField("#").addField(field).(*ObjectValidator)

	_, fields := expandFields(field, v.fields, v.field)

	rules := make([]ArrayRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.fields(fields)
	}

	return &ArrayValidator{
		field:       field,
		fields:      fields,
		required:    v.required,
		dflt:        v.dflt,
		rules:       rules,
		validator:   validator,
		errType:     invalidField(fields, InvalidArrayType, nil),
		errRequired: invalidField(fields, Required, nil),
	}
}

func (v *ArrayValidator) Required() *ArrayValidator {
	v.required = true
	return v
}

func (v *ArrayValidator) Default(value []typed.Typed) *ArrayValidator {
	v.dflt = value
	return v
}

func (v *ArrayValidator) Validator(validator *ObjectValidator) *ArrayValidator {
	v.validator = validator
	return v
}

func (v *ArrayValidator) Min(min int) *ArrayValidator {
	v.rules = append(v.rules, ArrayMin{min: min})
	return v
}

func (v *ArrayValidator) Max(max int) *ArrayValidator {
	v.rules = append(v.rules, ArrayMax{max: max})
	return v
}

func (v *ArrayValidator) Range(min int, max int) *ArrayValidator {
	v.rules = append(v.rules, ArrayRange{
		min: min,
		max: max,
	})
	return v
}

type ArrayMin struct {
	min int
	err InvalidField
}

func (r ArrayMin) Validate(fields []string, values []typed.Typed, object typed.Typed, input typed.Typed, res *Result) []typed.Typed {
	if len(values) < r.min {
		res.addField(r.err)
	}
	return values
}

func (r ArrayMin) fields(fields []string) ArrayRule {
	min := r.min
	return ArrayMin{
		min: min,
		err: invalidField(fields, InvalidArrayMinLength, Min(min), min),
	}
}

type ArrayMax struct {
	max int
	err InvalidField
}

func (r ArrayMax) Validate(fields []string, values []typed.Typed, object typed.Typed, input typed.Typed, res *Result) []typed.Typed {
	if len(values) > r.max {
		res.addField(r.err)
	}
	return values
}

func (r ArrayMax) fields(fields []string) ArrayRule {
	max := r.max
	return ArrayMax{
		max: max,
		err: invalidField(fields, InvalidArrayMaxLength, Max(max), max),
	}
}

type ArrayRange struct {
	min int
	max int
	err InvalidField
}

func (r ArrayRange) Validate(fields []string, values []typed.Typed, object typed.Typed, input typed.Typed, res *Result) []typed.Typed {
	if len(values) < r.min || len(values) > r.max {
		res.addField(r.err)
	}
	return values
}

func (r ArrayRange) fields(fields []string) ArrayRule {
	min := r.min
	max := r.max
	return ArrayRange{
		min: min,
		max: max,
		err: invalidField(fields, InvalidArrayRangeLength, Range(min, max), min, max),
	}
}
