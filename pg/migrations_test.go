package pg

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"src.goblgobl.com/tests/assert"
)

func Test_MigrateAll_NormalRun(t *testing.T) {
	bg := context.Background()
	db.Exec(bg, "drop table if exists test_migrations")
	db.Exec(bg, "drop table if exists gobl_migrations")
	migrateTest := func() {
		err := MigrateAll(db, []Migration{
			Migration{1, MigrateOne},
			Migration{2, MigrateTwo},
		})
		assert.Nil(t, err)

		value, err := Scalar[int](db, "select * from test_migrations")
		assert.Nil(t, err)
		assert.Equal(t, value, 9001)

		var version int
		rows, _ := db.Query(bg, "select version from gobl_migrations order by version")
		defer rows.Close()

		rows.Next()
		rows.Scan(&version)
		assert.Equal(t, version, 1)

		rows.Next()
		rows.Scan(&version)
		assert.Equal(t, version, 2)

		assert.False(t, rows.Next())
	}

	migrateTest()
	migrateTest() // this should be a noop
}

func Test_MigrateAll_Error(t *testing.T) {
	bg := context.Background()
	db.Exec(bg, "drop table if exists test_migrations")
	db.Exec(bg, "drop table if exists gobl_migrations")
	migrateTest := func() {
		err := MigrateAll(db, []Migration{
			Migration{1, MigrateOne},
			Migration{2, MigrateTwo},
			Migration{3, MigrateErr},
		})
		assert.StringContains(t, err.Error(), "Failed to run pg migration #3")

		value, err := Scalar[int](db, "select * from test_migrations")
		assert.Nil(t, err)
		assert.Equal(t, value, 9001)

		var version int
		rows, _ := db.Query(bg, "select version from gobl_migrations order by version")
		defer rows.Close()

		rows.Next()
		rows.Scan(&version)
		assert.Equal(t, version, 1)

		rows.Next()
		rows.Scan(&version)
		assert.Equal(t, version, 2)

		assert.False(t, rows.Next())
	}

	migrateTest()
	migrateTest() // this should be a noop
}

func MigrateOne(tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), "create table test_migrations (id integer not null)")
	return err
}

func MigrateTwo(tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), "insert into test_migrations(id) values (9001)")
	return err
}

func MigrateErr(tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), "fail")
	return err
}
