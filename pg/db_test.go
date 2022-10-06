package pg

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"src.goblgobl.com/tests"
	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils"
)

var db DB

func init() {
	var err error
	db, err = New(tests.PG())
	if err != nil {
		panic(err)
	}
}

func Test_New_Invalid(t *testing.T) {
	_, err := New("nope")
	assert.Equal(t, err.Error(), "code: 3003 - cannot parse `nope`: failed to parse as DSN (invalid dsn)")
}

func Test_Scalar(t *testing.T) {
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

func Test_DB_TableExist(t *testing.T) {
	exists, err := db.TableExists("test_migrations")
	assert.Nil(t, err)
	assert.False(t, exists)

	db.Exec(context.Background(), "create table if not exists test_migrations (id int)")
	exists, err = db.TableExists("test_migrations")
	db.Exec(context.Background(), "drop table test_migrations")
	assert.Nil(t, err)
	assert.True(t, exists)
}

func Test_DB_Transaction_Rollback(t *testing.T) {
	forcedErr := errors.New("forced error")
	err := db.Transaction(func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), "create table if not exists test_migrations (id int)")
		assert.Nil(t, err)
		return forcedErr
	})
	assert.True(t, errors.Is(err, forcedErr))
	exists, err := db.TableExists("test_migrations")
	assert.Nil(t, err)
	assert.False(t, exists)
}

func Test_DB_Transaction_Commit(t *testing.T) {
	err := db.Transaction(func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), "create table if not exists test_migrations (id int)")
		assert.Nil(t, err)
		return nil
	})
	assert.Nil(t, err)
	exists, _ := db.TableExists("test_migrations")
	db.Exec(context.Background(), "drop table test_migrations")
	assert.True(t, exists)
}

func testDB(fn func(DB)) {

	defer db.Close()
	fn(db)
}
