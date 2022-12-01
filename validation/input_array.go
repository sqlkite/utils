package validation

import (
	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
)

func Array() *ArrayValidator {
	return &ArrayValidator{}
}

type ArrayValidator struct {
	field       string
	fields      []string
	required    bool
	dflt        []typed.Typed
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

	return &ArrayValidator{
		field:       field,
		fields:      fields,
		required:    v.required,
		dflt:        v.dflt,
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
