package validation

import (
	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
	"src.sqlkite.com/utils/uuid"
)

func UUID(field string) *InputUUID {
	return &InputUUID{
		field:       field,
		errType:     inputError(field, InvalidUUIDType, nil),
		errRequired: inputError(field, Required, nil),
	}
}

type InputUUID struct {
	field       string
	dflt        string
	required    bool
	errType     InvalidField
	errRequired InvalidField
}

func (i *InputUUID) argsToTyped(args *fasthttp.Args, t typed.Typed) {
	field := i.field
	if value := args.Peek(field); value != nil {
		t[field] = string(value)
	}
}

func (i *InputUUID) Clone(field string) *InputUUID {
	return &InputUUID{
		field:       field,
		dflt:        i.dflt,
		required:    i.required,
		errType:     inputError(field, InvalidUUIDType, nil),
		errRequired: inputError(field, Required, nil),
	}
}

func (i *InputUUID) validate(input typed.Typed, res *Result) {
	field := i.field
	value, exists := input.StringIf(field)

	if !exists {
		if _, exists = input[field]; !exists && i.required {
			res.add(i.errRequired)
		} else if exists {
			res.add(i.errType)
		}
		if dflt := i.dflt; dflt != "" {
			input[field] = dflt
		}
		return
	}

	if !uuid.IsValid(value) {
		res.add(i.errType)
	}
}

func (i *InputUUID) Required() *InputUUID {
	i.required = true
	return i
}

func (i *InputUUID) Default(value string) *InputUUID {
	i.dflt = value
	return i
}
