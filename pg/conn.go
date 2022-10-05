package pg

import (
	"context"

	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Conn struct {
	*pgxpool.Pool
}

func New(url string) (Conn, error) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return Conn{}, log.Err(utils.ERR_PG_INIT, err).String("url", url)
	}
	return Conn{pool}, nil
}
