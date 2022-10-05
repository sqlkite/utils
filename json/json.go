// Wraps the JSON library that we're using so that we
// can [easily??] change it

package json

import (
	stdlib "encoding/json"

	json "github.com/goccy/go-json"
)

type Number stdlib.Number

func Marshal(data any) ([]byte, error) {
	return json.Marshal(data)
}

func Unmarshal(data []byte, into any) error {
	return json.Unmarshal(data, into)
}
