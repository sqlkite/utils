package sqlite

import (
	"testing"

	"src.sqlkite.com/tests/assert"
)

func Test_MigrateAll_NormalRun(t *testing.T) {
	migrateTest := func(conn Conn) {
		err := MigrateAll(conn, []Migration{
			Migration{1, MigrateOne},
			Migration{2, MigrateTwo},
		})
		assert.Nil(t, err)

		value, err := Scalar[int](conn, "select * from test_migrations")
		assert.Nil(t, err)
		assert.Equal(t, value, 9001)

		var version int
		rows := conn.Rows("select version from sqlkite_migrations order by version")
		defer rows.Close()

		rows.Next()
		rows.Scan(&version)
		assert.Equal(t, version, 1)

		rows.Next()
		rows.Scan(&version)
		assert.Equal(t, version, 2)

		assert.False(t, rows.Next())
	}

	testConn(func(conn Conn) {
		migrateTest(conn)
		migrateTest(conn) // this should be a noop
	})
}

func Test_MigrateAll_Error(t *testing.T) {
	migrateTest := func(conn Conn) {
		err := MigrateAll(conn, []Migration{
			Migration{1, MigrateOne},
			Migration{2, MigrateTwo},
			Migration{3, MigrateErr},
		})
		assert.StringContains(t, err.Error(), "Failed to run sqlite migration #3")

		value, err := Scalar[int](conn, "select * from test_migrations")
		assert.Nil(t, err)
		assert.Equal(t, value, 9001)

		var version int
		rows := conn.Rows("select version from sqlkite_migrations order by version")
		defer rows.Close()

		rows.Next()
		rows.Scan(&version)
		assert.Equal(t, version, 1)

		rows.Next()
		rows.Scan(&version)
		assert.Equal(t, version, 2)

		assert.False(t, rows.Next())
	}

	testConn(func(conn Conn) {
		migrateTest(conn)

		// this should be a noop
		migrateTest(conn)
	})
}

func MigrateOne(conn Conn) error {
	return conn.Exec("create table test_migrations (id integer not null)")
}

func MigrateTwo(conn Conn) error {
	return conn.Exec("insert into test_migrations(id) values (9001)")
}

func MigrateErr(conn Conn) error {
	return conn.Exec("fail")
}
