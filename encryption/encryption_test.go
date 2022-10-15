package encryption

import (
	"crypto/rand"
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Encrypt_Decrypt(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)

	value, err := Encrypt(key, "it's over 9000!!")
	assert.Nil(t, err)

	plain, err := Decrypt(key, value)
	assert.Nil(t, err)
	assert.Equal(t, string(plain), "it's over 9000!!")

}
