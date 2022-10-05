package log

import (
	"errors"
	"strings"
	"testing"

	"src.goblgobl.com/tests/assert"
)

// stupid test
func Test_Noop_Is_Noop(t *testing.T) {
	out := &strings.Builder{}
	l := Noop{}

	l.Info("x").LogTo(out)
	assert.Equal(t, out.String(), "")

	l.Warn("x").LogTo(out)
	assert.Equal(t, out.String(), "")

	l.Error("x").LogTo(out)
	assert.Equal(t, out.String(), "")

	l.Fatal("x").LogTo(out)
	assert.Equal(t, out.String(), "")

	l.Fatal("x").
		String("s", "s").
		Int("i", 1).
		Int64("i64", 99).
		LogTo(out)
	assert.Equal(t, out.String(), "")

	l.Field(NewField().Int("power", 9001).Finalize()).Fixed()
	l.String("s", "s1").MultiUse()
	assert.True(t, l.Err(errors.New("nope")).Bytes() == nil)
	l.Reset()
}
