package http

import (
	"time"

	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils/log"
	"src.goblgobl.com/utils/uuid"
)

type Env interface {
	Release()
	RequestId() string
	Info(string) log.Logger
	Error(string) log.Logger
}

func Handler[T Env](routeName string, loadEnv func(ctx *fasthttp.RequestCtx) (T, bool), next func(ctx *fasthttp.RequestCtx, env T) (Response, error)) func(ctx *fasthttp.RequestCtx) {
	return func(conn *fasthttp.RequestCtx) {
		start := time.Now()

		env, ok := loadEnv(conn)
		if !ok {
			// loadEnv is responsible for writing the respoinse
			return
		}
		defer env.Release()

		header := &conn.Response.Header
		header.SetContentTypeBytes([]byte("application/json"))
		header.SetBytesK([]byte("RequestId"), env.RequestId())

		r, err := next(conn, env)

		var logger log.Logger
		if err == nil {
			logger = env.Info("req")
		} else {
			// We could log the error directly in the req log
			// but this could contain sensitive or private information
			// (e.g. it could come from a 3rd party library, say, a failed
			// smtp request, and for all we know, it includes the email + smtp password)
			// So we log a separate env-free error (no pid or rid)
			// and tie this error to the req log via the errorId.
			errorId := uuid.String()
			header.SetBytesK([]byte("Error-Id"), errorId)
			logger = env.Error("env_handler_err").String("eid", errorId).Err(err)
			r = GenericServerError
		}

		r.Write(conn)

		r.EnhanceLog(logger).
			String("route", routeName).
			Int64("ms", time.Now().Sub(start).Milliseconds()).
			Log()
	}
}
