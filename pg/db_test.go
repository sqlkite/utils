package pg

import (
	"errors"
	"testing"

	"src.goblgobl.com/tests"
	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils"
)

func Test_New_Invalid(t *testing.T) {
	_, err := New("nope")
	assert.Equal(t, err.Error(), "code: 3003 - cannot parse `nope`: failed to parse as DSN (invalid dsn)")
}

func Test_Scalar(t *testing.T) {
	db, err := New(tests.PG())
	assert.Nil(t, err)
	defer db.Close()

	value, err := Scalar[int](db, "select 566 from doesnotexist")
	assert.Equal(t, err.Error(), `ERROR: relation "doesnotexist" does not exist (SQLSTATE 42P01)`)
	assert.Equal(t, value, 0)

	value, err = Scalar[int](db, "select 566")
	assert.Nil(t, err)
	assert.Equal(t, value, 566)

	value, err = Scalar[int](db, "select 566 where $1", false)
	assert.True(t, errors.Is(err, utils.ErrNoRows))
	assert.Equal(t, value, 0)

	str, err := Scalar[string](db, "select 'hello'")
	assert.Nil(t, err)
	assert.Equal(t, str, "hello")
}
