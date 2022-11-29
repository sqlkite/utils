package validation

import (
	"encoding/hex"
	"testing"

	"github.com/valyala/fasthttp"
	"src.sqlkite.com/tests/assert"
	"src.sqlkite.com/utils/typed"
)

func Test_String_Required(t *testing.T) {
	f1 := String()
	f2 := String().Required()
	i := Input().
		Field("name", f1).Field("name_clone", f1).
		Field("code", f2).Field("code_clone", f2)

	_, res := testInput(i)
	assert.Validation(t, res).
		FieldsHaveNoErrors("name", "name_clone").
		Field("code", Required).
		Field("code_clone", Required)

	_, res = testInput(i, "code", "1", "code_clone", "1")
	assert.Validation(t, res).
		FieldsHaveNoErrors("code", "name", "code_clone", "name_clone")
}

func Test_String_Default(t *testing.T) {
	f1 := String().Default("leto")
	f2 := String().Required().Default("leto")
	i := Input().
		Field("a", f1).Field("a_clone", f1).
		Field("b", f2).Field("b_clone", f2)

	// default doesn't really make sense with required, required
	// takes precedence
	data, res := testInput(i)
	assert.Equal(t, data.String("a"), "leto")
	assert.Equal(t, data.String("a_clone"), "leto")
	assert.Validation(t, res).
		Field("b", Required).
		Field("b_clone", Required)
}

func Test_String_Type(t *testing.T) {
	i := Input().
		Field("name", String())

	_, res := testInput(i, "name", 3)
	assert.Validation(t, res).
		Field("name", InvalidStringType)
}

func Test_String_Length(t *testing.T) {
	f1 := String().Length(0, 3)
	f2 := String().Length(2, 0)
	f3 := String().Length(2, 4)
	i := Input().
		Field("f1", f1).Field("f1_clone", f1).
		Field("f2", f2).Field("f2_clone", f2).
		Field("f3", f3).Field("f3_clone", f3)

	_, res := testInput(i, "f1", "1234", "f2", "1", "f3", "1", "f1_clone", "1234", "f2_clone", "1", "f3_clone", "1")
	assert.Validation(t, res).
		Field("f1", InvalidStringLength, map[string]any{"min": 0, "max": 3}).
		Field("f2", InvalidStringLength, map[string]any{"min": 2, "max": 0}).
		Field("f3", InvalidStringLength, map[string]any{"min": 2, "max": 4}).
		Field("f1_clone", InvalidStringLength, map[string]any{"min": 0, "max": 3}).
		Field("f2_clone", InvalidStringLength, map[string]any{"min": 2, "max": 0}).
		Field("f3_clone", InvalidStringLength, map[string]any{"min": 2, "max": 4})

	_, res = testInput(i, "f1", "123", "f2", "12", "f3", "12345", "f1_clone", "123", "f2_clone", "12", "f3_clone", "12345")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2", "f1_clone", "f2_clone").
		Field("f3", InvalidStringLength, map[string]any{"min": 2, "max": 4}).
		Field("f3_clone", InvalidStringLength, map[string]any{"min": 2, "max": 4})

	_, res = testInput(i, "f1", "1", "f2", "123456677", "f3", "12", "f1_clone", "1", "f2_clone", "123456677", "f3_clone", "12")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2", "f3", "f1_clone", "f2_clone", "f3_clone")

	_, res = testInput(i, "f3", "1234", "f3_clone", "1234")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f3", "f3_clone")

	_, res = testInput(i, "f3", "123", "f3_clone", "123")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f3", "f3_clone")
}

func Test_String_Pattern(t *testing.T) {
	f1 := String().Pattern("\\d.")
	i := Input().
		Field("f", f1).Field("f_clone", f1)

	_, res := testInput(i, "f", "1d", "f_clone", "1d")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f", "f_clone")

	_, res = testInput(i, "f", "1", "f_clone", "1")
	assert.Validation(t, res).
		Field("f", InvalidStringPattern, nil).
		Field("f_clone", InvalidStringPattern, nil)
}

func Test_String_Func(t *testing.T) {
	f1 := String().Func(func(field string, value string, input typed.Typed, res *Result) string {
		if value == "a" {
			return "a1"
		}
		res.InvalidField(field, InvalidStringPattern, nil)
		return value
	})

	i := Input().Field("f", f1).Field("f_clone", f1)

	data, res := testInput(i, "f", "a", "f_clone", "a")
	assert.Equal(t, data.String("f"), "a1")
	assert.Equal(t, data.String("f_clone"), "a1")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f", "f_clone")

	data, res = testInput(i, "f", "b", "f_clone", "b")
	assert.Equal(t, data.String("f"), "b")
	assert.Equal(t, data.String("f_clone"), "b")
	assert.Validation(t, res).
		Field("f", InvalidStringPattern, nil).
		Field("f_clone", InvalidStringPattern, nil)
}

func Test_String_Converter(t *testing.T) {
	f1 := String().Convert(func(field string, value string, input typed.Typed, res *Result) any {
		b, err := hex.DecodeString(value)
		if err == nil {
			return b
		}
		res.InvalidField(field, InvalidStringPattern, nil)
		return nil
	})

	i := Input().Field("f", f1).Field("f_clone", f1)

	data, res := testInput(i, "f", "FFFe", "f_clone", "FFFe")
	assert.Bytes(t, data.Bytes("f"), []byte{255, 254})
	assert.Bytes(t, data.Bytes("f_clone"), []byte{255, 254})
	assert.Validation(t, res).
		FieldsHaveNoErrors("f", "f_clone")

	data, res = testInput(i, "f", "z", "f_clone", "z")
	assert.True(t, data.Bytes("f") == nil)
	assert.True(t, data.Bytes("f_clone") == nil)
	assert.Validation(t, res).
		Field("f", InvalidStringPattern, nil).
		Field("f_clone", InvalidStringPattern, nil)
}

func Test_String_Args(t *testing.T) {
	i := Input().
		Field("name", String().Required().Length(4, 4))

	_, res := testArgs(i, "name", "leto")
	assert.Validation(t, res).FieldsHaveNoErrors("name")
}

func Test_Int_Required(t *testing.T) {
	i := Input().
		Field("name", Int()).
		Field("code", Int().Required())

	_, res := testInput(i)
	assert.Validation(t, res).
		FieldsHaveNoErrors("name").
		Field("code", Required)

	_, res = testInput(i, "code", 1)
	assert.Validation(t, res).
		FieldsHaveNoErrors("code", "name")
}

func Test_Int_Type(t *testing.T) {
	i := Input().
		Field("a", Int())

	_, res := testInput(i, "a", "leto")
	assert.Validation(t, res).
		Field("a", InvalidIntType)

	data, res := testInput(i, "a", "-3292")
	assert.Validation(t, res).
		FieldsHaveNoErrors("a")
	assert.Equal(t, data.Int("a"), -3292)
}

func Test_Int_Default(t *testing.T) {
	i := Input().
		Field("a", Int().Default(99)).
		Field("b", Int().Required().Default(88))

	// default doesn't really make sense with required, required
	// takes precedence
	data, res := testInput(i)
	assert.Equal(t, data.Int("a"), 99)
	assert.Validation(t, res).
		Field("b", Required)
}

func Test_Int_MinMax(t *testing.T) {
	i := Input().
		Field("f1", Int().Min(10)).
		Field("f2", Int().Max(10))

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
		Field("f1", Int().Range(10, 20))

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
		Field("f", Int().Func(func(field string, value int, input typed.Typed, res *Result) int {
			if value == 9001 {
				return 9002
			}
			res.InvalidField(field, InvalidIntMax, nil)
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

func Test_Int_Args(t *testing.T) {
	i := Input().Field("id", Int().Required().Range(4, 4))
	input, res := testArgs(i, "id", "4")
	assert.Validation(t, res).FieldsHaveNoErrors("id")
	assert.Equal(t, input.Int("id"), 4)

	input, res = testArgs(i, "id", "nope")
	assert.Validation(t, res).Field("id", InvalidIntType)
	assert.Equal(t, input.IntOr("id", -1), -1)
}

func Test_Bool_Required(t *testing.T) {
	i := Input().
		Field("required", Bool()).
		Field("age", Bool().Required())

	_, res := testInput(i)
	assert.Validation(t, res).
		FieldsHaveNoErrors("required").
		Field("agree", Required)

	_, res = testInput(i, "agree", true)
	assert.Validation(t, res).
		FieldsHaveNoErrors("required", "agree")
}

func Test_Bool_Type(t *testing.T) {
	i := Input().
		Field("a", Bool())

	_, res := testInput(i, "a", "leto")
	assert.Validation(t, res).
		Field("a", InvalidBoolType)

	data, res := testInput(i, "a", "true")
	assert.Validation(t, res).
		FieldsHaveNoErrors("a")
	assert.Equal(t, data.Bool("a"), true)
}

func Test_Bool_Default(t *testing.T) {
	i := Input().
		Field("a", Bool().Default(true)).
		Field("b", Bool().Required().Default(true))

	// default doesn't really make sense with required, required
	// takes precedence
	data, res := testInput(i)
	assert.Equal(t, data.Bool("a"), true)
	assert.Validation(t, res).
		Field("b", Required)
}

func Test_Bool_Func(t *testing.T) {
	i := Input().
		Field("f", Bool().Func(func(field string, value bool, input typed.Typed, res *Result) bool {
			if value == false {
				return true
			}
			res.InvalidField(field, InvalidBoolType, nil)
			return value
		}))

	data, res := testInput(i, "f", false)
	assert.Equal(t, data.Bool("f"), true)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f")

	data, res = testInput(i, "f", true)
	assert.Equal(t, data.Bool("f"), true)
	assert.Validation(t, res).
		Field("f", InvalidBoolType, nil)
}

func Test_Bool_Args(t *testing.T) {
	i := Input().Field("agree", Bool().Required())
	for _, value := range []string{"true", "TRUE", "True"} {
		input, res := testArgs(i, "agree", value)
		assert.Validation(t, res).FieldsHaveNoErrors("agree")
		assert.True(t, input.Bool("agree"))
	}

	for _, value := range []string{"false", "FALSE", "False"} {
		input, res := testArgs(i, "agree", value)
		assert.Validation(t, res).FieldsHaveNoErrors("agree")
		assert.False(t, input.Bool("agree"))
	}

	input, res := testArgs(i, "agree", "other")
	assert.Validation(t, res).Field("agree", InvalidBoolType)
	_, isBool := input.BoolIf("agree")
	assert.False(t, isBool)
}

func Test_UUID_Required(t *testing.T) {
	f1 := UUID()
	f2 := UUID().Required()
	i := Input().
		Field("id", f1).Field("id_clone", f1).
		Field("parent_id", f2).Field("parent_id_clone", f2)

	_, res := testInput(i)
	assert.Validation(t, res).
		FieldsHaveNoErrors("id", "id_clone").
		Field("parent_id", Required).
		Field("parent_id_clone", Required)

	_, res = testInput(i, "parent_id", "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF", "parent_id_clone", "00000000-0000-0000-0000-000000000000")
	assert.Validation(t, res).
		FieldsHaveNoErrors("parent_id", "id", "parent_id_clone", "id_clone")
}

func Test_UUID_Type(t *testing.T) {
	i := Input().Field("id", UUID())

	_, res := testInput(i, "id", 3)
	assert.Validation(t, res).
		Field("id", InvalidUUIDType)

	_, res = testInput(i, "id", "Z0000000-0000-0000-0000-00000000000Z")
	assert.Validation(t, res).
		Field("id", InvalidUUIDType)
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

func testArgs(i *input, args ...string) (typed.Typed, *Result) {
	m := new(fasthttp.Args)
	for i := 0; i < len(args); i += 2 {
		m.Add(args[i], args[i+1])
	}

	res := NewResult(5)
	input, _ := i.ValidateArgs(m, res)
	return input, res
}
