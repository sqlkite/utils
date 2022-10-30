package http

import (
	"testing"

	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils/log"
	"src.goblgobl.com/utils/typed"
	"src.goblgobl.com/utils/validation"

	"github.com/valyala/fasthttp"
)

type TestResponse struct {
	status int
	body   string
	json   typed.Typed
	log    map[string]string
}

func Test_Ok_NoBody(t *testing.T) {
	res := read(Ok(nil))
	assert.Equal(t, res.status, 200)
	assert.Equal(t, len(res.body), 0)
	assert.Equal(t, res.log["res"], "0")
	assert.Equal(t, res.log["status"], "200")
}

func Test_Ok_Body(t *testing.T) {
	res := read(Ok(map[string]any{"over": 9000}))
	assert.Equal(t, res.status, 200)
	assert.Equal(t, res.body, `{"over":9000}`)
	assert.Equal(t, res.log["res"], "13")
	assert.Equal(t, res.log["status"], "200")
}

func Test_Ok_InvalidBody(t *testing.T) {
	res := read(Ok(make(chan bool)))

	errorId := res.log["eid"]
	assert.Equal(t, len(errorId), 36)
	assert.Equal(t, res.log["res"], "95")
	assert.Equal(t, res.log["code"], "2002")
	assert.Equal(t, res.log["status"], "500")

	assert.Equal(t, res.status, 500)
	assert.Equal(t, res.json.Int("code"), 2002)
	assert.Equal(t, res.json.String("error"), "internal server error")
	assert.Equal(t, res.json.String("error_id"), errorId)

}

func Test_StaticNotFound(t *testing.T) {
	res := read(StaticNotFound(1023))
	assert.Equal(t, res.status, 404)
	assert.Equal(t, res.body, `{"code":1023,"error":"not found"}`)
	assert.Equal(t, res.log["res"], "33")
	assert.Equal(t, res.log["code"], "1023")
	assert.Equal(t, res.log["status"], "404")
}

func Test_StaticError(t *testing.T) {
	res := read(StaticError(511, 1002, "oops"))
	assert.Equal(t, res.status, 511)
	assert.Equal(t, res.body, `{"code":1002,"error":"oops"}`)
	assert.Equal(t, res.log["res"], "28")
	assert.Equal(t, res.log["code"], "1002")
	assert.Equal(t, res.log["status"], "511")
}

func Test_ServerError(t *testing.T) {
	res := read(ServerError())
	assert.Equal(t, res.status, 500)
	assert.Equal(t, res.json.Int("code"), 2001)
	assert.Equal(t, res.json.String("error"), "internal server error")

	errorId := res.json.String("error_id")
	assert.Equal(t, len(errorId), 36)

	assert.Equal(t, res.log["res"], "95")
	assert.Equal(t, res.log["code"], "2001")
	assert.Equal(t, res.log["status"], "500")
	assert.Equal(t, res.log["eid"], errorId)
}

func Test_Validation(t *testing.T) {
	result := validation.NewResult(5)
	result.InvalidField("field1", validation.Required, nil)
	result.InvalidField("field2", validation.Required, 331)
	result.Invalid(validation.InvalidStringType, nil)
	result.Invalid(validation.InvalidStringType, map[string]any{"over": 9000})

	res := read(Validation(result))
	assert.Equal(t, res.status, 400)
	assert.Equal(t, res.json.Int("code"), 2004)
	assert.Equal(t, res.json.String("error"), "invalid data")

	invalid := res.json.Objects("invalid")
	assert.Equal(t, len(invalid), 4)
	assert.Equal(t, invalid[0].Int("code"), 1001)
	assert.Equal(t, invalid[0].String("field"), "field1")
	assert.Equal(t, invalid[0].String("error"), "required")
	assert.Nil(t, invalid[0].Object("data"))

	assert.Equal(t, invalid[1].Int("code"), 1001)
	assert.Equal(t, invalid[1].String("field"), "field2")
	assert.Equal(t, invalid[1].String("error"), "required")
	assert.Equal(t, invalid[1].Int("data"), 331)

	assert.Equal(t, invalid[2].Int("code"), 1002)
	assert.Equal(t, invalid[2].String("field"), "")
	assert.Equal(t, invalid[2].String("error"), "must be a string")
	assert.Nil(t, invalid[2].Object("data"))

	assert.Equal(t, invalid[3].Int("code"), 1002)
	assert.Equal(t, invalid[3].String("field"), "")
	assert.Equal(t, invalid[3].String("error"), "must be a string")
	assert.Equal(t, invalid[3].Object("data").Int("over"), 9000)

	assert.Equal(t, res.log["res"], "262")
	assert.Equal(t, res.log["code"], "2004")
	assert.Equal(t, res.log["status"], "400")
}

func read(res Response) TestResponse {
	conn := &fasthttp.RequestCtx{}
	res.Write(conn)

	body := conn.Response.Body()
	var json typed.Typed
	if len(body) > 0 {
		json = typed.Must(body)
	}

	logger := log.Info("test")
	defer logger.Release()
	res.EnhanceLog(logger)

	return TestResponse{
		json:   json,
		body:   string(body),
		status: conn.Response.StatusCode(),
		log:    log.KvParse(string(logger.Bytes())),
	}
}
