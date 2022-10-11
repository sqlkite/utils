package validation

import (
	"testing"

	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils/typed"
)

func Test_String_Required(t *testing.T) {
	i := Input().
		Field(String("name", false)).
		Field(String("code", true))

	assert.Validation(t, testInput(i)).
		FieldsHaveNoErrors("name").
		Field("code", Required)

	assert.Validation(t, testInput(i, "code", "1")).
		FieldsHaveNoErrors("code", "name")
}

func Test_String_Type(t *testing.T) {
	i := Input().
		Field(String("name", false))

	assert.Validation(t, testInput(i, "name", 3)).
		Field("name", InvalidStringType)
}

func Test_String_Length(t *testing.T) {
	i := Input().
		Field(String("f1", false).Length(0, 3)).
		Field(String("f2", false).Length(2, 0)).
		Field(String("f3", false).Length(2, 4))

	assert.Validation(t, testInput(i, "f1", "1234", "f2", "1", "f3", "1")).
		Field("f1", InvalidStringLength, map[string]any{"min": 0, "max": 3}).
		Field("f2", InvalidStringLength, map[string]any{"min": 2, "max": 0}).
		Field("f3", InvalidStringLength, map[string]any{"min": 2, "max": 4})

	assert.Validation(t, testInput(i, "f1", "123", "f2", "12", "f3", "12345")).
		FieldsHaveNoErrors("f1", "f2").
		Field("f3", InvalidStringLength, map[string]any{"min": 2, "max": 4})

	assert.Validation(t, testInput(i, "f1", "1", "f2", "123456677", "f3", "12")).
		FieldsHaveNoErrors("f1", "f2", "f3")

	assert.Validation(t, testInput(i, "f3", "1234")).
		FieldsHaveNoErrors("f3")

	assert.Validation(t, testInput(i, "f3", "123")).
		FieldsHaveNoErrors("f3")
}

func Test_String_Pattern(t *testing.T) {
	i := Input().
		Field(String("f", false).Pattern("\\d."))

	assert.Validation(t, testInput(i, "f", "1d")).
		FieldsHaveNoErrors("f")

	assert.Validation(t, testInput(i, "f", "1")).
		Field("f", InvalidStringPattern, nil)
}

func Test_String_Func(t *testing.T) {
	i := Input().
		Field(String("f", false).Func(func(field string, value string, input typed.Typed, res *Result) string {
			if value == "a" {
				return "a1"
			}
			res.add(inputError(field, InvalidStringPattern, nil))
			return value
		}))

	assert.Validation(t, testInput(i, "f", "a")).
		FieldsHaveNoErrors("f")

	assert.Validation(t, testInput(i, "f", "b")).
		Field("f", InvalidStringPattern, nil)
}

func Test_Int_Required(t *testing.T) {
	i := Input().
		Field(Int("name", false)).
		Field(Int("code", true))

	assert.Validation(t, testInput(i)).
		FieldsHaveNoErrors("name").
		Field("code", Required)

	assert.Validation(t, testInput(i, "code", 1)).
		FieldsHaveNoErrors("code", "name")
}

func Test_Int_MinMax(t *testing.T) {
	i := Input().
		Field(Int("f1", false).Min(10)).
		Field(Int("f2", false).Max(10))

	assert.Validation(t, testInput(i, "f1", 9, "f2", 11)).
		Field("f1", InvalidIntMin, map[string]any{"min": 10}).
		Field("f2", InvalidIntMax, map[string]any{"max": 10})

	assert.Validation(t, testInput(i, "f1", 10, "f2", 10)).
		FieldsHaveNoErrors("f1", "f2")

	assert.Validation(t, testInput(i, "f1", 11, "f2", 9)).
		FieldsHaveNoErrors("f1", "f2")
}

func Test_Int_Range(t *testing.T) {
	i := Input().
		Field(Int("f1", false).Range(10, 20))

	for _, value := range []int{9, 21, 0, 30} {
		assert.Validation(t, testInput(i, "f1", value)).
			Field("f1", InvalidIntRange, map[string]any{"min": 10, "max": 20})
	}

	for _, value := range []int{10, 11, 19, 20} {
		assert.Validation(t, testInput(i, "f1", value)).
			FieldsHaveNoErrors("f1")
	}

	assert.Validation(t, testInput(i, "f1", 21)).
		Field("f1", InvalidIntRange, map[string]any{"min": 10, "max": 20})
}

func Test_Int_Func(t *testing.T) {
	i := Input().
		Field(Int("f", false).Func(func(field string, value int, input typed.Typed, res *Result) int {
			if value == 9001 {
				return 9002
			}
			res.add(inputError(field, InvalidIntMax, nil))
			return value
		}))

	assert.Validation(t, testInput(i, "f", 9001)).
		FieldsHaveNoErrors("f")

	assert.Validation(t, testInput(i, "f", 8000)).
		Field("f", InvalidIntMax, nil)
}

func testInput(i *input, args ...any) *Result {
	m := make(typed.Typed, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		m[args[i].(string)] = args[i+1]
	}

	res := NewResult(5)
	i.Validate(m, res)
	return res
}
