package http

import (
	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils"
	"src.sqlkite.com/utils/json"
	"src.sqlkite.com/utils/log"
	"src.sqlkite.com/utils/uuid"
)

var (
	serverErrorLogData = log.NewField().
				Int("code", utils.RES_SERVER_ERROR).
				Int("status", 500).
				Finalize()

	serializationErrorLogData = log.NewField().
					Int("code", utils.RES_SERIALIZATION_ERROR).
					Int("status", 500).
					Finalize()
)

type ErrorIdResponse struct {
	errorId string
	body    []byte
	logData log.Field
}

func (r ErrorIdResponse) Write(conn *fasthttp.RequestCtx) {
	conn.SetStatusCode(500)
	conn.Response.Header.SetBytesK([]byte("Error-Id"), r.errorId)
	conn.SetBody(r.body)
}

func (r ErrorIdResponse) EnhanceLog(logger log.Logger) log.Logger {
	logger.Field(r.logData).String("eid", r.errorId).Int("res", len(r.body))
	return logger
}

func ServerError() Response {
	errorId := uuid.String()

	data := struct {
		Code    int    `json:"code"`
		Error   string `json:"error"`
		ErrorId string `json:"error_id"`
	}{
		ErrorId: errorId,
		Code:    utils.RES_SERVER_ERROR,
		Error:   "internal server error",
	}
	body, _ := json.Marshal(data)

	return ErrorIdResponse{
		body:    body,
		errorId: errorId,
		logData: serverErrorLogData,
	}
}

func SerializationError() Response {
	errorId := uuid.String()

	data := struct {
		Code    int    `json:"code"`
		Error   string `json:"error"`
		ErrorId string `json:"error_id"`
	}{
		ErrorId: errorId,
		Code:    utils.RES_SERIALIZATION_ERROR,
		Error:   "internal server error",
	}
	body, _ := json.Marshal(data)

	return ErrorIdResponse{
		body:    body,
		errorId: errorId,
		logData: serializationErrorLogData,
	}
}
