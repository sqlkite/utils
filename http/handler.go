package http

import (
	"time"

	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils/log"
)

type Env interface {
	Release()
	RequestId() string
	Info(string) log.Logger
	Error(string) log.Logger
}

func Handler[T Env](routeName string, loadEnv func(ctx *fasthttp.RequestCtx) (T, Response, error), next func(ctx *fasthttp.RequestCtx, env T) (Response, error)) func(ctx *fasthttp.RequestCtx) {
	return func(conn *fasthttp.RequestCtx) {
		start := time.Now()

		var haveEnv bool
		var logger log.Logger
		env, res, err := loadEnv(conn)

		if res == nil && err == nil {
			haveEnv = true
			defer env.Release()

			header := &conn.Response.Header
			header.SetContentTypeBytes([]byte("application/json"))
			header.SetBytesK([]byte("RequestId"), env.RequestId())

			res, err = next(conn, env)
		}

		if err == nil {
			if haveEnv {
				logger = env.Info("req")
			} else {
				logger = log.Error("req")
			}
		} else {
			if haveEnv {
				logger = env.Error("handler").Err(err)
			} else {
				logger = log.Error("handler").Err(err)
			}
			res = ServerError()
		}

		res.Write(conn)

		res.EnhanceLog(logger).
			String("route", routeName).
			Int64("ms", time.Now().Sub(start).Milliseconds()).
			Log()
	}
}