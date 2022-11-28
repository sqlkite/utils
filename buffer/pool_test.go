package buffer

import (
	"testing"

	"src.sqlkite.com/tests/assert"
)

func Test_Pool_Checkout_and_Release(t *testing.T) {
	b := Checkout(22)
	b.Write([]byte("abc"))
	assert.Equal(t, b.max, 22)
	assert.Equal(t, len(b.data), 65536)

	Release(b)

	s, err := b.String()
	assert.Nil(t, err)
	assert.Equal(t, s, "")
}

func Test_Pool_R(t *testing.T) {
	b := Checkout(22)
	assert.Equal(t, b.max, 22)
	assert.Equal(t, len(b.data), 65536)
}
