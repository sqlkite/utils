package json

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Marshal(t *testing.T) {
	data, err := Marshal(map[string]any{"over": 9000})
	assert.Nil(t, err)
	assert.Equal(t, string(data), `{"over":9000}`)

	data, err = Marshal(make(chan int))
	assert.Equal(t, len(data), 0)
	assert.Equal(t, err.Error(), "json: unsupported type: chan int")
}

func Test_Unmarshal(t *testing.T) {
	var into map[string]any
	err := Unmarshal([]byte(`{"leto":"ghanima"}`), &into)
	assert.Nil(t, err)
	assert.Equal(t, into["leto"].(string), "ghanima")

	err = Unmarshal([]byte(`{`), &into)
	assert.True(t, err != nil)
}
