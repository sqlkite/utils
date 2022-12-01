package validation

import (
	"testing"

	"src.sqlkite.com/tests/assert"
	"src.sqlkite.com/utils/json"
	"src.sqlkite.com/utils/typed"
)

func Test_Result_InvalidField_NoData(t *testing.T) {
	r := NewResult(4)
	assert.True(t, r.IsValid())

	r.InvalidField([]string{"field1"}, Required, nil)
	assert.False(t, r.IsValid())

	invalid := r.Errors()[0].(InvalidField)
	assert.Nil(t, invalid.Data)
	assert.List(t, invalid.Fields, []string{"field1"})
	assert.Equal(t, invalid.Code, 1001)
	assert.Equal(t, invalid.Error, "required")
}

func Test_Result_InvalidField_Data(t *testing.T) {
	r := NewResult(4)
	r.InvalidField([]string{"field1"}, InvalidStringType, 33)
	invalid := r.Errors()[0].(InvalidField)
	assert.Equal(t, invalid.Data.(int), 33)
	assert.Equal(t, invalid.Code, 1002)
	assert.List(t, invalid.Fields, []string{"field1"})
	assert.Equal(t, invalid.Error, "must be a string")
}

func Test_Result_Invalid_NoData(t *testing.T) {
	r := NewResult(4)
	r.Invalid(InvalidStringPattern, nil)
	invalid := r.Errors()[0].(Invalid)
	assert.Nil(t, invalid.Data)
	assert.Equal(t, invalid.Code, 1004)
	assert.Equal(t, invalid.Error, "is not valid")
}

func Test_Result_Invalid_Range(t *testing.T) {
	r := NewResult(4)
	r.Invalid(InvalidStringLength, Range(5, 10))
	invalid := r.Errors()[0].(Invalid)
	assert.Equal(t, invalid.Code, 1003)
	assert.Equal(t, invalid.Data.(DataRange).Min.(int), 5)
	assert.Equal(t, invalid.Data.(DataRange).Max.(int), 10)
	assert.Equal(t, invalid.Error, "must be between %d and %d characters")
}

func serialize(r *Result) typed.Typed {
	data, err := json.Marshal(r.Errors()[0])
	if err != nil {
		panic(err)
	}
	return typed.Must(data)
}
