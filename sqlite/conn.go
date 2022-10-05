package sqlite

import (
	driver "src.goblgobl.com/sqlite"
	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/log"
)

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
