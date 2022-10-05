package http

/*
Generic errors that don't chnange. Things like internal server error.
*/

import (
	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/json"
	"src.goblgobl.com/utils/log"

	"github.com/valyala/fasthttp"
)

var (
	GenericServerError       = StaticServerError(utils.RES_SERVER_ERROR)
	GenericFailedToSerialize = StaticServerError(utils.RES_SERIALIZE_STATIC)
	InvalidJSON              = StaticError(400, utils.RES_INVALID_JSON_PAYLOAD, "invalid json payload")
)

// We know the status/body/logData upfront (lets us optimize
// EnhanceLog)
type StaticResponse struct {
	status  int
	body    []byte
	logData log.Field
}

func (r StaticResponse) Write(conn *fasthttp.RequestCtx) {
	conn.SetStatusCode(r.status)
	conn.SetBody(r.body)
}

func (r StaticResponse) EnhanceLog(logger log.Logger) log.Logger {
	logger.Field(r.logData)
	return logger
}

func StaticError(status int, code int, error string) StaticResponse {
	data := struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{
		Code:  code,
		Error: error,
	}
	body, _ := json.Marshal(data)

	logData := log.NewField().
		Int("code", code).
		Int("status", status).
		Int("res", len(body)).
		Finalize()

	return StaticResponse{
		body:    body,
		status:  status,
		logData: logData,
	}
}

func StaticNotFound(code int) StaticResponse {
	return StaticError(404, code, "not found")
}

func StaticServerError(code int) StaticResponse {
	return StaticError(500, code, "internal server error")
}

func StaticUnavailableError(code int) StaticResponse {
	return StaticError(503, code, "service unavailable")
}