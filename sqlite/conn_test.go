package sqlite

import (
	"errors"
	"strings"
	"testing"

	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils"
)

func Test_New_InvalidPath(t *testing.T) {
	_, err := New("/tmp/hopefully/does/not/exist", false)
	assert.Equal(t, err.Error(), "code: 3004 - file does not exist")
}

func Test_New_Success(t *testing.T) {
	conn, err := New(":memory:", false)
	assert.Nil(t, err)
	defer conn.Close()

	var value int
	exists, err := conn.Row("select 1").Scan(&value)
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, value, 1)
}

func Test_Scalar(t *testing.T) {
	conn, err := New(":memory:", false)
	assert.Nil(t, err)
	defer conn.Close()

	value, err := Scalar[int](conn, "select 123 from doesnotexist")
	assert.True(t, strings.Contains(err.Error(), "sqlite: no such table: doesnotexist (code: 1)"))
	assert.Equal(t, value, 0)

	value, err = Scalar[int](conn, "select 566")
	assert.Nil(t, err)
	assert.Equal(t, value, 566)

	value, err = Scalar[int](conn, "select 566 where $1", false)
	assert.True(t, errors.Is(err, utils.ErrNoRows))
	assert.Equal(t, value, 0)

	str, err := Scalar[string](conn, "select 'hello'")
	assert.Nil(t, err)
	assert.Equal(t, str, "hello")
}
