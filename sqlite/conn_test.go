package sqlite

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"src.goblgobl.com/sqlite"
	"src.goblgobl.com/tests/assert"
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
	err = conn.Row("select 1").Scan(&value)
	assert.Nil(t, err)
	assert.Equal(t, value, 1)
}

func Test_Scalar(t *testing.T) {
	testConn(func(conn Conn) {
		value, err := Scalar[int](conn, "select 123 from doesnotexist")
		assert.False(t, conn.IsNotFound(err))
		assert.True(t, strings.Contains(err.Error(), "sqlite: no such table: doesnotexist (code: 1)"))
		assert.Equal(t, value, 0)

		value, err = Scalar[int](conn, "select 566")
		assert.Nil(t, err)
		assert.Equal(t, value, 566)

		value, err = Scalar[int](conn, "select 566 where $1", false)
		assert.True(t, conn.IsNotFound(err))
		assert.True(t, errors.Is(err, sqlite.ErrNoRows))
		assert.Equal(t, value, 0)

		str, err := Scalar[string](conn, "select 'hello'")
		assert.Nil(t, err)
		assert.Equal(t, str, "hello")
	})
}

func Test_Conn_TableExist(t *testing.T) {
	testConn(func(conn Conn) {
		exists, err := conn.TableExists("migrations")
		assert.Nil(t, err)
		assert.False(t, exists)

		conn.Exec("create table migrations (id int)")
		exists, err = conn.TableExists("migrations")
		assert.Nil(t, err)
		assert.True(t, exists)
	})
}

func Test_Conn_Placeholder(t *testing.T) {
	conn := Conn{}
	for i := 0; i < 50; i++ {
		assert.Equal(t, conn.Placeholder(i), fmt.Sprintf("?%d", i+1))
	}
}

func Test_Conn_RowToMap(t *testing.T) {
	testConn(func(conn Conn) {
		m, err := conn.RowToMap("select 1 where false")
		assert.True(t, errors.Is(err, ErrNoRows))
		assert.Equal(t, len(m), 0)

		m, err = conn.RowToMap("select 1 as a, 'b' as b")
		assert.Nil(t, err)
		assert.Equal(t, len(m), 2)
		assert.Equal(t, m.Int("a"), 1)
		assert.Equal(t, m.String("b"), "b")
	})
}

func testConn(fn func(Conn)) {
	conn, err := New(":memory:", false)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fn(conn)
}
