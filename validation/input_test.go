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

	_, res := testInput(i)
	assert.Validation(t, res).
		FieldsHaveNoErrors("name").
		Field("code", Required)

	_, res = testInput(i, "code", "1")
	assert.Validation(t, res).
		FieldsHaveNoErrors("code", "name")
}

func Test_String_Default(t *testing.T) {
	i := Input().
		Field(String("a", false).Default("leto")).
		Field(String("b", true).Default("leto"))

	// default doesn't really make sense with required, required
	// takes precedence
	data, res := testInput(i)
	assert.Equal(t, data.String("a"), "leto")
	assert.Validation(t, res).
		Field("b", Required)
}

func Test_String_Type(t *testing.T) {
	i := Input().
		Field(String("name", false))

	_, res := testInput(i, "name", 3)
	assert.Validation(t, res).
		Field("name", InvalidStringType)
}

func Test_String_Length(t *testing.T) {
	i := Input().
		Field(String("f1", false).Length(0, 3)).
		Field(String("f2", false).Length(2, 0)).
		Field(String("f3", false).Length(2, 4))

	_, res := testInput(i, "f1", "1234", "f2", "1", "f3", "1")
	assert.Validation(t, res).
		Field("f1", InvalidStringLength, map[string]any{"min": 0, "max": 3}).
		Field("f2", InvalidStringLength, map[string]any{"min": 2, "max": 0}).
		Field("f3", InvalidStringLength, map[string]any{"min": 2, "max": 4})

	_, res = testInput(i, "f1", "123", "f2", "12", "f3", "12345")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2").
		Field("f3", InvalidStringLength, map[string]any{"min": 2, "max": 4})

	_, res = testInput(i, "f1", "1", "f2", "123456677", "f3", "12")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2", "f3")

	_, res = testInput(i, "f3", "1234")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f3")

	_, res = testInput(i, "f3", "123")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f3")
}

func Test_String_Pattern(t *testing.T) {
	i := Input().
		Field(String("f", false).Pattern("\\d."))

	_, res := testInput(i, "f", "1d")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f")

	_, res = testInput(i, "f", "1")
	assert.Validation(t, res).
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

	data, res := testInput(i, "f", "a")
	assert.Equal(t, data.String("f"), "a1")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f")

	data, res = testInput(i, "f", "b")
	assert.Equal(t, data.String("f"), "b")
	assert.Validation(t, res).
		Field("f", InvalidStringPattern, nil)
}

func Test_Int_Required(t *testing.T) {
	i := Input().
		Field(Int("name", false)).
		Field(Int("code", true))

	_, res := testInput(i)
	assert.Validation(t, res).
		FieldsHaveNoErrors("name").
		Field("code", Required)

	_, res = testInput(i, "code", 1)
	assert.Validation(t, res).
		FieldsHaveNoErrors("code", "name")
}

func Test_Int_Default(t *testing.T) {
	i := Input().
		Field(Int("a", false).Default(99)).
		Field(Int("b", true).Default(88))

	// default doesn't really make sense with required, required
	// takes precedence
	data, res := testInput(i)
	assert.Equal(t, data.Int("a"), 99)
	assert.Validation(t, res).
		Field("b", Required)
}

func Test_Int_MinMax(t *testing.T) {
	i := Input().
		Field(Int("f1", false).Min(10)).
		Field(Int("f2", false).Max(10))

	_, res := testInput(i, "f1", 9, "f2", 11)
	assert.Validation(t, res).
		Field("f1", InvalidIntMin, map[string]any{"min": 10}).
		Field("f2", InvalidIntMax, map[string]any{"max": 10})

	_, res = testInput(i, "f1", 10, "f2", 10)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2")

	_, res = testInput(i, "f1", 11, "f2", 9)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2")
}

func Test_Int_Range(t *testing.T) {
	i := Input().
		Field(Int("f1", false).Range(10, 20))

	for _, value := range []int{9, 21, 0, 30} {
		_, res := testInput(i, "f1", value)
		assert.Validation(t, res).
			Field("f1", InvalidIntRange, map[string]any{"min": 10, "max": 20})
	}

	for _, value := range []int{10, 11, 19, 20} {
		_, res := testInput(i, "f1", value)
		assert.Validation(t, res).
			FieldsHaveNoErrors("f1")
	}

	_, res := testInput(i, "f1", 21)
	assert.Validation(t, res).
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

	data, res := testInput(i, "f", 9001)
	assert.Equal(t, data.Int("f"), 9002)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f")

	data, res = testInput(i, "f", 8000)
	assert.Equal(t, data.Int("f"), 8000)
	assert.Validation(t, res).
		Field("f", InvalidIntMax, nil)
}

func testInput(i *input, args ...any) (typed.Typed, *Result) {
	m := make(typed.Typed, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		m[args[i].(string)] = args[i+1]
	}

	res := NewResult(5)
	i.Validate(m, res)
	return m, res
}
