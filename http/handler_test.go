package http

import (
	"errors"
	"strings"
	"testing"

	"github.com/valyala/fasthttp"
	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils/log"
	"src.goblgobl.com/utils/typed"
)

func Test_Handler_EnvLoaderFail(t *testing.T) {
	testLoader := func(conn *fasthttp.RequestCtx) (*TestEnv, bool) {
		conn.SetStatusCode(1)
		return nil, false
	}

	conn := &fasthttp.RequestCtx{}
	Handler("", testLoader, func(conn *fasthttp.RequestCtx, env *TestEnv) (Response, error) {
		assert.Fail(t, "next should not be called")
		return nil, nil
	})(conn)

	assert.Equal(t, conn.Response.StatusCode(), 1)
}

func Test_Handler_CallsHandlerWithEnv(t *testing.T) {
	testLoader := func(conn *fasthttp.RequestCtx) (*TestEnv, bool) {
		return testEnv(200), true
	}

	conn := &fasthttp.RequestCtx{}
	Handler("", testLoader, func(conn *fasthttp.RequestCtx, env *TestEnv) (Response, error) {
		assert.Equal(t, env.id, 200)
		return StaticError(2, 2, ""), nil
	})(conn)
	assert.Equal(t, conn.Response.StatusCode(), 2)
}

func Test_Handler_LogsResponse(t *testing.T) {
	out := &strings.Builder{}
	testLoader := func(conn *fasthttp.RequestCtx) (*TestEnv, bool) {
		e := testEnv(201)
		e.logger = log.NewKvLogger(1024, out, nil)
		return e, true
	}

	conn := &fasthttp.RequestCtx{}
	Handler("test-route", testLoader, func(conn *fasthttp.RequestCtx, env *TestEnv) (Response, error) {
		return StaticNotFound(9001), nil
	})(conn)

	reqLog := log.KvParse(out.String())
	assert.Equal(t, reqLog["l"], "info")
	assert.Equal(t, reqLog["status"], "404")
	assert.Equal(t, reqLog["route"], "test-route")
	assert.Equal(t, reqLog["res"], "33")
	assert.Equal(t, reqLog["code"], "9001")
	assert.Equal(t, reqLog["c"], "req")
}

func Test_Server_Handler_LogsError(t *testing.T) {
	out := &strings.Builder{}
	testLoader := func(conn *fasthttp.RequestCtx) (*TestEnv, bool) {
		e := testEnv(202)
		e.logger = log.NewKvLogger(1024, out, nil)
		return e, true
	}

	conn := &fasthttp.RequestCtx{}
	Handler("test2", testLoader, func(conn *fasthttp.RequestCtx, env *TestEnv) (Response, error) {
		return nil, errors.New("Not Over 9000!")
	})(conn)

	res := conn.Response
	assert.Equal(t, res.StatusCode(), 500)
	reqLog := log.KvParse(out.String())
	assert.Equal(t, reqLog["l"], "error")
	assert.Equal(t, reqLog["status"], "500")
	assert.Equal(t, reqLog["route"], "test2")
	assert.Equal(t, reqLog["res"], "45")
	assert.Equal(t, reqLog["code"], "2001")
	assert.Equal(t, reqLog["c"], "env_handler_err")
	assert.Equal(t, reqLog["err"], `"Not Over 9000!"`)
	assert.Equal(t, reqLog["eid"], string(res.Header.Peek("Error-Id")))
}

type TestEnv struct {
	id       int
	released bool
	logger   log.Logger
}

func testEnv(id int) *TestEnv {
	return &TestEnv{
		id:     id,
		logger: log.Noop{},
	}
}

func (e *TestEnv) Release() {
	e.released = true
	e.logger.Release()

}

func (e TestEnv) RequestId() string {
	return ""
}

func (e TestEnv) Info(ctx string) log.Logger {
	return e.logger.Info(ctx)
}

func (e TestEnv) Error(ctx string) log.Logger {
	return e.logger.Error(ctx)
}

func assertCode(t *testing.T, conn *fasthttp.RequestCtx, expected int) {
	t.Helper()
	res := conn.Response
	body := res.Body()
	json, _ := typed.Json(body)
	assert.Equal(t, json.Int("code"), expected)
}
