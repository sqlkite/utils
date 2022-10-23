package encryption

import (
	"crypto/rand"
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Encrypt_Decrypt(t *testing.T) {
	var key [32]byte
	rand.Read(key[:])

	value, err := Encrypt(key, "it's over 9000!!")
	assert.Nil(t, err)
	plain, ok := Decrypt(key, value)
	assert.True(t, ok)
	assert.Equal(t, string(plain), "it's over 9000!!")
}
