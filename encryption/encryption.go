package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"src.goblgobl.com/utils"
)

type Value struct {
	Nonce []byte `json:"nonce"`
	Data  []byte `json:"data"`
}

func Encrypt(key []byte, plainText string) (Value, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return Value{}, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return Value{}, err
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return Value{}, err
	}

	encrypted := gcm.Seal(nil, nonce, utils.S2B(plainText), nil)
	return Value{Nonce: nonce, Data: encrypted}, nil
}

func Decrypt(key []byte, value Value) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, value.Nonce, value.Data, nil)
}
