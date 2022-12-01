package validation

import (
	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
	"src.sqlkite.com/utils/uuid"
)

func UUID() *UUIDValidator {
	return &UUIDValidator{}
}

type UUIDValidator struct {
	field       string
	fields      []string
	dflt        string
	required    bool
	errType     InvalidField
	errRequired InvalidField
}

func (v *UUIDValidator) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	field := v.field
	if value := args.Peek(field); value != nil {
		t[field] = string(value)
	}
}

func (v *UUIDValidator) validate(object typed.Typed, input typed.Typed, res *Result) {
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

	if !uuid.IsValid(value) {
		res.addField(v.errType)
	}
}

func (v *UUIDValidator) addField(field string) InputValidator {
	field, fields := expandFields(field, v.fields, v.field)
	return &UUIDValidator{
		field:       field,
		fields:      fields,
		dflt:        v.dflt,
		required:    v.required,
		errType:     invalidField(fields, InvalidUUIDType, nil),
		errRequired: invalidField(fields, Required, nil),
	}
}

func (v *UUIDValidator) Required() *UUIDValidator {
	v.required = true
	return v
}

func (v *UUIDValidator) Default(value string) *UUIDValidator {
	v.dflt = value
	return v
}
