package sqlite

import (
	"errors"

	driver "src.goblgobl.com/sqlite"
	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/log"
)

var ErrNoRows = errors.New("no rows in result set")

type Conn struct {
	driver.Conn
}

func New(filePath string, create bool) (Conn, error) {
	conn, err := driver.Open(filePath, create)
	if err != nil {
		return Conn{}, log.Err(utils.ERR_SQLITE_INIT, err).String("path", filePath)
	}
	return Conn{conn}, nil
}

func Scalar[T any](conn Conn, sql string, args ...any) (T, error) {
	row := conn.Conn.Row(sql, args...)

	var value T
	exists, err := row.Scan(&value)
	if err != nil {
		return value, err
	}
	if !exists {
		return value, utils.ErrNoRows
	}
	return value, err
}

func (c Conn) TableExists(tableName string) (bool, error) {
	sql := `
		select exists (
			select 1 from sqlite_master
			where type = 'table' and name = ?1
		)
	`
	exists, err := Scalar[bool](c, sql, tableName)
	if err == utils.ErrNoRows {
		return false, nil
	}
	return exists, err
}
