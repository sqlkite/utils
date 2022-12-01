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
	o := Object().
		Field("name", f1).Field("name_clone", f1).
		Field("code", f2).Field("code_clone", f2)

	_, res := testInput(o)
	assert.Validation(t, res).
		FieldsHaveNoErrors("name", "name_clone").
		Field("code", Required).
		Field("code_clone", Required)

	_, res = testInput(o, "code", "1", "code_clone", "1")
	assert.Validation(t, res).
		FieldsHaveNoErrors("code", "name", "code_clone", "name_clone")
}

func Test_String_Default(t *testing.T) {
	f1 := String().Default("leto")
	f2 := String().Required().Default("leto")
	o := Object().
		Field("a", f1).Field("a_clone", f1).
		Field("b", f2).Field("b_clone", f2)

	// default doesn't really make sense with required, required
	// takes precedence
	data, res := testInput(o)
	assert.Equal(t, data.String("a"), "leto")
	assert.Equal(t, data.String("a_clone"), "leto")
	assert.Validation(t, res).
		Field("b", Required).
		Field("b_clone", Required)
}

func Test_String_Type(t *testing.T) {
	o := Object().
		Field("name", String())

	_, res := testInput(o, "name", 3)
	assert.Validation(t, res).
		Field("name", InvalidStringType)
}

func Test_String_Length(t *testing.T) {
	f1 := String().Length(0, 3)
	f2 := String().Length(2, 0)
	f3 := String().Length(2, 4)
	o := Object().
		Field("f1", f1).Field("f1_clone", f1).
		Field("f2", f2).Field("f2_clone", f2).
		Field("f3", f3).Field("f3_clone", f3)

	_, res := testInput(o, "f1", "1234", "f2", "1", "f3", "1", "f1_clone", "1234", "f2_clone", "1", "f3_clone", "1")
	assert.Validation(t, res).
		Field("f1", InvalidStringLength, map[string]any{"min": 0, "max": 3}).
		Field("f2", InvalidStringLength, map[string]any{"min": 2, "max": 0}).
		Field("f3", InvalidStringLength, map[string]any{"min": 2, "max": 4}).
		Field("f1_clone", InvalidStringLength, map[string]any{"min": 0, "max": 3}).
		Field("f2_clone", InvalidStringLength, map[string]any{"min": 2, "max": 0}).
		Field("f3_clone", InvalidStringLength, map[string]any{"min": 2, "max": 4})

	_, res = testInput(o, "f1", "123", "f2", "12", "f3", "12345", "f1_clone", "123", "f2_clone", "12", "f3_clone", "12345")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2", "f1_clone", "f2_clone").
		Field("f3", InvalidStringLength, map[string]any{"min": 2, "max": 4}).
		Field("f3_clone", InvalidStringLength, map[string]any{"min": 2, "max": 4})

	_, res = testInput(o, "f1", "1", "f2", "123456677", "f3", "12", "f1_clone", "1", "f2_clone", "123456677", "f3_clone", "12")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2", "f3", "f1_clone", "f2_clone", "f3_clone")

	_, res = testInput(o, "f3", "1234", "f3_clone", "1234")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f3", "f3_clone")

	_, res = testInput(o, "f3", "123", "f3_clone", "123")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f3", "f3_clone")
}

func Test_String_Choice(t *testing.T) {
	f1 := String().Choice("c1", "c2")
	o1 := Object().
		Field("f", f1).Field("f_clone", f1)

	_, res := testInput(o1, "f", "c1", "f_clone", "c2")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f", "f_clone")

	_, res = testInput(o1, "f", "nope", "f_clone", "C2") // case sensitive
	assert.Validation(t, res).
		Field("f", InvalidStringChoice, map[string]any{"valid": []string{"c1", "c2"}}).
		Field("f_clone", InvalidStringChoice, map[string]any{"valid": []string{"c1", "c2"}})
}

func Test_String_Pattern(t *testing.T) {
	f1 := String().Pattern("\\d.")
	o1 := Object().
		Field("f", f1).Field("f_clone", f1)

	_, res := testInput(o1, "f", "1d", "f_clone", "1d")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f", "f_clone")

	_, res = testInput(o1, "f", "1", "f_clone", "1")
	assert.Validation(t, res).
		Field("f", InvalidStringPattern, nil).
		FieldMessage("f", "is not valid"). // default/generic error
		Field("f_clone", InvalidStringPattern, nil).
		FieldMessage("f_clone", "is not valid") // default/generic error

	// explicit error message
	f2 := String().Pattern("^\\d$", "must be a number")
	o2 := Object().
		Field("f", f2).Field("f_clone", f2)

	_, res = testInput(o2, "f", "1d", "f_clone", "1d")
	assert.Validation(t, res).
		Field("f", InvalidStringPattern, nil).
		FieldMessage("f", "must be a number").
		Field("f_clone", InvalidStringPattern, nil).
		FieldMessage("f_clone", "must be a number")
}

func Test_String_Func(t *testing.T) {
	f1 := String().Func(func(path []string, value string, object typed.Typed, input typed.Typed, res *Result) string {
		if value == "a" {
			return "a1"
		}
		res.InvalidField(path, InvalidStringPattern, nil)
		return value
	})

	o := Object().Field("f", f1).Field("f_clone", f1)

	data, res := testInput(o, "f", "a", "f_clone", "a")
	assert.Equal(t, data.String("f"), "a1")
	assert.Equal(t, data.String("f_clone"), "a1")
	assert.Validation(t, res).
		FieldsHaveNoErrors("f", "f_clone")

	data, res = testInput(o, "f", "b", "f_clone", "b")
	assert.Equal(t, data.String("f"), "b")
	assert.Equal(t, data.String("f_clone"), "b")
	assert.Validation(t, res).
		Field("f", InvalidStringPattern, nil).
		Field("f_clone", InvalidStringPattern, nil)
}

func Test_String_Converter(t *testing.T) {
	f1 := String().Convert(func(path []string, value string, object typed.Typed, input typed.Typed, res *Result) any {
		b, err := hex.DecodeString(value)
		if err == nil {
			return b
		}
		res.InvalidField(path, InvalidStringPattern, nil)
		return nil
	})

	o := Object().Field("f", f1).Field("f_clone", f1)

	data, res := testInput(o, "f", "FFFe", "f_clone", "FFFe")
	assert.Bytes(t, data.Bytes("f"), []byte{255, 254})
	assert.Bytes(t, data.Bytes("f_clone"), []byte{255, 254})
	assert.Validation(t, res).
		FieldsHaveNoErrors("f", "f_clone")

	data, res = testInput(o, "f", "z", "f_clone", "z")
	assert.True(t, data.Bytes("f") == nil)
	assert.True(t, data.Bytes("f_clone") == nil)
	assert.Validation(t, res).
		Field("f", InvalidStringPattern, nil).
		Field("f_clone", InvalidStringPattern, nil)
}

func Test_String_Args(t *testing.T) {
	o := Object().
		Field("name", String().Required().Length(4, 4))

	_, res := testArgs(o, "name", "leto")
	assert.Validation(t, res).FieldsHaveNoErrors("name")
}

func Test_Int_Required(t *testing.T) {
	o := Object().
		Field("name", Int()).
		Field("code", Int().Required())

	_, res := testInput(o)
	assert.Validation(t, res).
		FieldsHaveNoErrors("name").
		Field("code", Required)

	_, res = testInput(o, "code", 1)
	assert.Validation(t, res).
		FieldsHaveNoErrors("code", "name")
}

func Test_Int_Type(t *testing.T) {
	o := Object().
		Field("a", Int())

	_, res := testInput(o, "a", "leto")
	assert.Validation(t, res).
		Field("a", InvalidIntType)

	data, res := testInput(o, "a", "-3292")
	assert.Validation(t, res).
		FieldsHaveNoErrors("a")
	assert.Equal(t, data.Int("a"), -3292)
}

func Test_Int_Default(t *testing.T) {
	o := Object().
		Field("a", Int().Default(99)).
		Field("b", Int().Required().Default(88))

	// default doesn't really make sense with required, required
	// takes precedence
	data, res := testInput(o)
	assert.Equal(t, data.Int("a"), 99)
	assert.Validation(t, res).
		Field("b", Required)
}

func Test_Int_MinMax(t *testing.T) {
	o := Object().
		Field("f1", Int().Min(10)).
		Field("f2", Int().Max(10))

	_, res := testInput(o, "f1", 9, "f2", 11)
	assert.Validation(t, res).
		Field("f1", InvalidIntMin, map[string]any{"min": 10}).
		Field("f2", InvalidIntMax, map[string]any{"max": 10})

	_, res = testInput(o, "f1", 10, "f2", 10)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2")

	_, res = testInput(o, "f1", 11, "f2", 9)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f1", "f2")
}

func Test_Int_Range(t *testing.T) {
	o := Object().
		Field("f1", Int().Range(10, 20))

	for _, value := range []int{9, 21, 0, 30} {
		_, res := testInput(o, "f1", value)
		assert.Validation(t, res).
			Field("f1", InvalidIntRange, map[string]any{"min": 10, "max": 20})
	}

	for _, value := range []int{10, 11, 19, 20} {
		_, res := testInput(o, "f1", value)
		assert.Validation(t, res).
			FieldsHaveNoErrors("f1")
	}

	_, res := testInput(o, "f1", 21)
	assert.Validation(t, res).
		Field("f1", InvalidIntRange, map[string]any{"min": 10, "max": 20})
}

func Test_Int_Func(t *testing.T) {
	o := Object().
		Field("f", Int().Func(func(path []string, value int, object typed.Typed, input typed.Typed, res *Result) int {
			if value == 9001 {
				return 9002
			}
			res.InvalidField(path, InvalidIntMax, nil)
			return value
		}))

	data, res := testInput(o, "f", 9001)
	assert.Equal(t, data.Int("f"), 9002)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f")

	data, res = testInput(o, "f", 8000)
	assert.Equal(t, data.Int("f"), 8000)
	assert.Validation(t, res).
		Field("f", InvalidIntMax, nil)
}

func Test_Int_Args(t *testing.T) {
	o := Object().Field("id", Int().Required().Range(4, 4))
	input, res := testArgs(o, "id", "4")
	assert.Validation(t, res).FieldsHaveNoErrors("id")
	assert.Equal(t, input.Int("id"), 4)

	input, res = testArgs(o, "id", "nope")
	assert.Validation(t, res).Field("id", InvalidIntType)
	assert.Equal(t, input.IntOr("id", -1), -1)
}

func Test_Bool_Required(t *testing.T) {
	o := Object().
		Field("required", Bool()).
		Field("agree", Bool().Required())

	_, res := testInput(o)
	assert.Validation(t, res).
		FieldsHaveNoErrors("required").
		Field("agree", Required)

	_, res = testInput(o, "agree", true)
	assert.Validation(t, res).
		FieldsHaveNoErrors("required", "agree")
}

func Test_Bool_Type(t *testing.T) {
	o := Object().
		Field("a", Bool())

	_, res := testInput(o, "a", "leto")
	assert.Validation(t, res).
		Field("a", InvalidBoolType)

	data, res := testInput(o, "a", "true")
	assert.Validation(t, res).
		FieldsHaveNoErrors("a")
	assert.Equal(t, data.Bool("a"), true)
}

func Test_Bool_Default(t *testing.T) {
	o := Object().
		Field("a", Bool().Default(true)).
		Field("b", Bool().Required().Default(true))

	// default doesn't really make sense with required, required
	// takes precedence
	data, res := testInput(o)
	assert.Equal(t, data.Bool("a"), true)
	assert.Validation(t, res).
		Field("b", Required)
}

func Test_Bool_Func(t *testing.T) {
	o := Object().
		Field("f", Bool().Func(func(path []string, value bool, object typed.Typed, input typed.Typed, res *Result) bool {
			if value == false {
				return true
			}
			res.InvalidField(path, InvalidBoolType, nil)
			return value
		}))

	data, res := testInput(o, "f", false)
	assert.Equal(t, data.Bool("f"), true)
	assert.Validation(t, res).
		FieldsHaveNoErrors("f")

	data, res = testInput(o, "f", true)
	assert.Equal(t, data.Bool("f"), true)
	assert.Validation(t, res).
		Field("f", InvalidBoolType, nil)
}

func Test_Bool_Args(t *testing.T) {
	o := Object().Field("agree", Bool().Required())
	for _, value := range []string{"true", "TRUE", "True"} {
		input, res := testArgs(o, "agree", value)
		assert.Validation(t, res).FieldsHaveNoErrors("agree")
		assert.True(t, input.Bool("agree"))
	}

	for _, value := range []string{"false", "FALSE", "False"} {
		input, res := testArgs(o, "agree", value)
		assert.Validation(t, res).FieldsHaveNoErrors("agree")
		assert.False(t, input.Bool("agree"))
	}

	input, res := testArgs(o, "agree", "other")
	assert.Validation(t, res).Field("agree", InvalidBoolType)
	_, isBool := input.BoolIf("agree")
	assert.False(t, isBool)
}

func Test_UUID_Required(t *testing.T) {
	f1 := UUID()
	f2 := UUID().Required()
	o := Object().
		Field("id", f1).Field("id_clone", f1).
		Field("parent_id", f2).Field("parent_id_clone", f2)

	_, res := testInput(o)
	assert.Validation(t, res).
		FieldsHaveNoErrors("id", "id_clone").
		Field("parent_id", Required).
		Field("parent_id_clone", Required)

	_, res = testInput(o, "parent_id", "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF", "parent_id_clone", "00000000-0000-0000-0000-000000000000")
	assert.Validation(t, res).
		FieldsHaveNoErrors("parent_id", "id", "parent_id_clone", "id_clone")
}

func Test_UUID_Type(t *testing.T) {
	o := Object().Field("id", UUID())

	_, res := testInput(o, "id", 3)
	assert.Validation(t, res).
		Field("id", InvalidUUIDType)

	_, res = testInput(o, "id", "Z0000000-0000-0000-0000-00000000000Z")
	assert.Validation(t, res).
		Field("id", InvalidUUIDType)
}

func Test_Nested_Object(t *testing.T) {
	child := Object().
		Field("age", Int().Required()).
		Field("name", String().Required())

	o1 := Object().Field("user", child)
	_, res := testInput(o1, "id", 3)

	assert.Validation(t, res).
		Field("user.age", Required).
		Field("user.name", Required)

	o2 := Object().Field("entry", o1)
	_, res = testInput(o2, "id", 3)
	assert.Validation(t, res).
		Field("entry.user.age", Required).
		Field("entry.user.name", Required)

	_, res = testInput(o2, "entry", typed.Typed{"user": typed.Typed{"age": 3000, "name": "Leto"}})
	assert.Validation(t, res).FieldsHaveNoErrors("entry.user.age", "entry.user.name")
}

func Test_Array_Object(t *testing.T) {
	child := Object().Field("name", String().Required())
	o1 := Object().
		Field("users", Array().Required().Validator(child))

	_, res := testInput(o1)
	assert.Validation(t, res).Field("users", Required)

	_, res = testInput(o1, "users", 1)
	assert.Validation(t, res).Field("users", InvalidArrayType)

	_, res = testInput(o1, "users", []typed.Typed{typed.Typed{}})
	assert.Validation(t, res).Field("users.0.name", Required)

	_, res = testInput(o1, "users", []typed.Typed{typed.Typed{"name": "leto"}})
	assert.Validation(t, res).FieldsHaveNoErrors("users.0.name")

	_, res = testInput(o1, "users", []typed.Typed{
		typed.Typed{"name": "leto"},
		typed.Typed{"name": 3},
	})
	assert.Validation(t, res).
		Field("users.1.name", InvalidStringType).
		FieldsHaveNoErrors("users.0.name")
}

func Test_Array_MinAndMax(t *testing.T) {
	createItem := func() typed.Typed {
		return typed.Typed{"name": "n"}
	}

	child := Object().Field("name", String())
	o1 := Object().Field("users", Array().Min(2).Max(3).Required().Validator(child))

	_, res := testInput(o1, "users", []typed.Typed{createItem()})
	assert.Validation(t, res).Field("users", InvalidArrayMinLength, map[string]any{"min": 2})

	// 4 items, too many
	_, res = testInput(o1, "users", []typed.Typed{
		createItem(), createItem(), createItem(), createItem(),
	})
	assert.Validation(t, res).Field("users", InvalidArrayMaxLength, map[string]any{"max": 3})

	// 2 items, good
	_, res = testInput(o1, "users", []typed.Typed{
		createItem(), createItem(),
	})
	assert.Validation(t, res).FieldsHaveNoErrors("users")

	// 3 items, good
	_, res = testInput(o1, "users", []typed.Typed{
		createItem(), createItem(), createItem(),
	})
	assert.Validation(t, res).FieldsHaveNoErrors("users")
}

func Test_Array_Range(t *testing.T) {
	createItem := func() typed.Typed {
		return typed.Typed{"name": "n"}
	}

	child := Object().Field("name", String())
	o1 := Object().Field("users", Array().Range(2, 3).Required().Validator(child))

	_, res := testInput(o1, "users", []typed.Typed{createItem()})
	assert.Validation(t, res).Field("users", InvalidArrayRangeLength, map[string]any{"min": 2, "max": 3})

	// 4 items, too many
	_, res = testInput(o1, "users", []typed.Typed{
		createItem(), createItem(), createItem(), createItem(),
	})
	assert.Validation(t, res).Field("users", InvalidArrayRangeLength, map[string]any{"min": 2, "max": 3})

	// 2 items, good
	_, res = testInput(o1, "users", []typed.Typed{
		createItem(), createItem(),
	})
	assert.Validation(t, res).FieldsHaveNoErrors("users")

	// 3 items, good
	_, res = testInput(o1, "users", []typed.Typed{
		createItem(), createItem(), createItem(),
	})
	assert.Validation(t, res).FieldsHaveNoErrors("users")

}

func Test_Any_Required(t *testing.T) {
	o := Object().
		Field("name", Any()).
		Field("code", Any().Required())

	_, res := testInput(o)
	assert.Validation(t, res).
		FieldsHaveNoErrors("name").
		Field("code", Required)

	_, res = testInput(o, "code", 1)
	assert.Validation(t, res).
		FieldsHaveNoErrors("code", "name")
}

func Test_Any_Default(t *testing.T) {
	o := Object().Field("name", Any().Default(32))

	data, _ := testInput(o)
	assert.Equal(t, data.Int("name"), 32)
}

func Test_Any_Func(t *testing.T) {
	o := Object().Field("name", Any().Func(func(path []string, value any, object typed.Typed, input typed.Typed, res *Result) any {
		assert.Equal(t, value.(string), "one-one")
		return 11
	}))

	data, _ := testInput(o, "name", "one-one")
	assert.Equal(t, data.Int("name"), 11)
}

func testInput(o *ObjectValidator, args ...any) (typed.Typed, *Result) {
	m := make(typed.Typed, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		m[args[i].(string)] = args[i+1]
	}

	res := NewResult(10)
	o.Validate(m, res)
	return m, res
}

func testArgs(o *ObjectValidator, args ...string) (typed.Typed, *Result) {
	m := new(fasthttp.Args)
	for i := 0; i < len(args); i += 2 {
		m.Add(args[i], args[i+1])
	}

	res := NewResult(5)
	input, _ := o.ValidateArgs(m, res)
	return input, res
}
