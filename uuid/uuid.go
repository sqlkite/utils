package uuid

import (
	"github.com/google/uuid"
)

func init() {
	uuid.EnableRandPool()
}

func IsValid(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

func String() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func FromBytes(bytes []byte) (string, error) {
	id, err := uuid.FromBytes(bytes)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
