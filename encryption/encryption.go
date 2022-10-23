package encryption

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
	"src.goblgobl.com/utils"
)

func Encrypt(key [32]byte, plainText string) ([]byte, error) {
	var nonce [24]byte
	nonceSlice := nonce[:]
	if _, err := io.ReadFull(rand.Reader, nonceSlice); err != nil {
		return nil, err
	}
	encrypted := secretbox.Seal(nonceSlice, utils.S2B(plainText), &nonce, &key)
	return encrypted, nil
}

func Decrypt(key [32]byte, encrypted []byte) ([]byte, bool) {
	nonce := (*[24]byte)(encrypted)
	return secretbox.Open(nil, encrypted[24:], nonce, &key)
}
