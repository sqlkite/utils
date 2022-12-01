package validation

/*
When we first create a validation rule, we have no field name. For
example:

	countValidation := Int().Required().Min(0)

Only once we add it to an object does it get a field name:

	Object().Field("count", countValidation)

This means that `countValidation` can be used with different field
names (either in the same object or in a different object):

	Object().
		Field("count", countValidation).
		Field("confirm_count", countValidation)

Adding a field to an object (e.g. Field(NAME, validationRule)) is
relatively expensive. We deep-clone the rule and, in most cases,
pre-create an InvalidField error.

This logic also applies for nested objects:

	Object().Field("user", Object().
		Field("name", nameValidation).
		Field("age", ageValidation)
	)

After this is called, the code that represents the name validation
will have a field of ["user", "name"]. This is done in two steps,
first, the call to `Field("name", nameValidation)` will give it a
field of ["name"]. Then when the "user" field is added all children
as passed the new parent, thus the field becomes ["user", "name"].

This approach is inneficient, as we create a temporary name validator
for the ["name"] field, and then replace it with the final validator
for the ["user", "name"] field. However, this only happens on startup.
*/

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/typed"
)

type InputValidator interface {
	addField(name string) InputValidator
	validate(object typed.Typed, input typed.Typed, res *Result)
	argsToTyped(args *fasthttp.Args, dest typed.Typed)
}

func invalid(meta Meta, data any, args ...any) Invalid {
	err := meta.Error
	if len(args) > 0 {
		err = fmt.Sprintf(err, args...)
	}
	return Invalid{
		Data:  data,
		Code:  meta.Code,
		Error: err,
	}
}

func invalidField(fields []string, meta Meta, data any, args ...any) InvalidField {
	err := meta.Error
	if len(args) > 0 {
		err = fmt.Sprintf(err, args...)
	}
	return InvalidField{
		Fields: fields,
		Invalid: Invalid{
			Data:  data,
			Code:  meta.Code,
			Error: err,
		},
	}
}

// Validators are built as:
//
//	Object().Field("user", Object().Field("name", String()))
//
// So the first call to addField is "name" and then "user".
// For this reason the field being added is prepended to our list, e.g.:
//  1. ["name"]
//  2. ["user", "name"]
//
// Furthermore, the first addField that we see becomes the key that
// we use to lookup the value (e.g. in the above example "name") is
// the actual field for the validator.
//
// In other words:
//  1. field gets prepended to fields
//  2. if existing hasn't been set yet, field becomes the new name
//  3. if existing has been set, it remains the field name
func expandFields(field string, fields []string, existing string) (string, []string) {
	if fields == nil {
		fields = []string{field}
	} else {
		fields = append([]string{field}, fields...)
	}

	if existing == "" {
		return field, fields
	}
	return existing, fields
}
