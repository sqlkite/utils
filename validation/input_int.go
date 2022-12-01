package validation

import (
	"strconv"

	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils"
	"src.sqlkite.com/utils/typed"
)

type IntRule interface {
	fields(fields []string) IntRule
	Validate(fields []string, value int, object typed.Typed, input typed.Typed, res *Result) int
}

func Int() *IntValidator {
	return &IntValidator{}
}

type IntValidator struct {
	field       string
	fields      []string
	dflt        int
	required    bool
	rules       []IntRule
	errType     InvalidField
	errRequired InvalidField
}

func (v *IntValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	field := v.field
	if value := args.Peek(field); value != nil {
		if n, err := strconv.ParseInt(utils.B2S(value), 10, 0); err == nil {
			t[field] = n
		} else {
			t[field] = value
		}
	}
}

func (v *IntValidator) validate(object typed.Typed, input typed.Typed, res *Result) {
	field := v.field
	value, exists := object.IntIf(field)

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

func (v *IntValidator) addField(field string) InputValidator {
	field, fields := expandFields(field, v.fields, v.field)

	rules := make([]IntRule, len(v.rules))
	for i, rule := range v.rules {
		rules[i] = rule.fields(fields)
	}

	return &IntValidator{
		field:       field,
		fields:      fields,
		dflt:        v.dflt,
		required:    v.required,
		rules:       rules,
		errType:     invalidField(fields, InvalidIntType, nil),
		errRequired: invalidField(fields, Required, nil),
	}
}

func (v *IntValidator) Required() *IntValidator {
	v.required = true
	return v
}

func (v *IntValidator) Default(value int) *IntValidator {
	v.dflt = value
	return v
}

func (v *IntValidator) Min(min int) *IntValidator {
	v.rules = append(v.rules, IntMin{min: min})
	return v
}

func (v *IntValidator) Max(max int) *IntValidator {
	v.rules = append(v.rules, IntMax{max: max})
	return v
}

func (v *IntValidator) Range(min int, max int) *IntValidator {
	v.rules = append(v.rules, IntRange{
		min: min,
		max: max,
	})
	return v
}

func (v *IntValidator) Func(fn func(fields []string, value int, object typed.Typed, input typed.Typed, res *Result) int) *IntValidator {
	v.rules = append(v.rules, IntFunc{fn: fn})
	return v
}

type IntMin struct {
	min int
	err InvalidField
}

func (r IntMin) Validate(fields []string, value int, object typed.Typed, input typed.Typed, res *Result) int {
	if value < r.min {
		res.addField(r.err)
	}
	return value
}

func (r IntMin) fields(fields []string) IntRule {
	min := r.min
	return IntMin{
		min: min,
		err: invalidField(fields, InvalidIntMin, Min(min), min),
	}
}

type IntMax struct {
	max int
	err InvalidField
}

func (r IntMax) Validate(fields []string, value int, object typed.Typed, input typed.Typed, res *Result) int {
	if value > r.max {
		res.addField(r.err)
	}
	return value
}

func (r IntMax) fields(fields []string) IntRule {
	max := r.max
	return IntMax{
		max: max,
		err: invalidField(fields, InvalidIntMax, Max(max), max),
	}
}

type IntRange struct {
	min int
	max int
	err InvalidField
}

func (r IntRange) Validate(fields []string, value int, object typed.Typed, input typed.Typed, res *Result) int {
	if value < r.min || value > r.max {
		res.addField(r.err)
	}
	return value
}

func (r IntRange) fields(fields []string) IntRule {
	min := r.min
	max := r.max
	return IntRange{
		min: min,
		max: max,
		err: invalidField(fields, InvalidIntRange, Range(min, max), min, max),
	}
}

type IntFunc struct {
	fn func([]string, int, typed.Typed, typed.Typed, *Result) int
}

func (v IntFunc) Validate(fields []string, value int, object typed.Typed, input typed.Typed, res *Result) int {
	return v.fn(fields, value, object, input, res)
}

func (r IntFunc) fields(fields []string) IntRule {
	return r
}
