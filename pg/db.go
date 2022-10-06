package pg

import (
	"context"

	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func New(url string) (DB, error) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return DB{}, log.Err(utils.ERR_PG_INIT, err).String("url", url)
	}
	return DB{pool}, nil
}

func Scalar[T any](db DB, sql string, args ...any) (T, error) {
	row := db.Pool.QueryRow(context.Background(), sql, args...)

	var value T
	err := row.Scan(&value)
	if err == pgx.ErrNoRows {
		return value, utils.ErrNoRows
	}
	return value, err
}

func (db DB) TableExists(tableName string) (bool, error) {
	sql := `
		select exists (
			select from pg_tables
			where schemaname = 'public' and tablename = $1
		)
	`
	exists, err := Scalar[bool](db, sql, tableName)
	if err == utils.ErrNoRows {
		return false, nil
	}
	return exists, err
}

func (db DB) Transaction(fn func(tx pgx.Tx) error) error {
	bg := context.Background()
	tx, err := db.Pool.Begin(bg)
	if err != nil {
		return err
	}

	defer tx.Rollback(bg)
	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(bg)
}
