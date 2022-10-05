package concurrent

import (
	"errors"
	"strconv"
	"testing"

	"src.goblgobl.com/tests/assert"
)

type TestItem struct {
	id string
}

func Test_Get_Loader(t *testing.T) {
	m := NewMap[*TestItem](func(id string) (*TestItem, error) {
		if id == "nope" {
			return nil, nil
		}
		if id == "error" {
			return nil, errors.New("invalid id")
		}
		return &TestItem{id: id}, nil
	})

	i, err := m.Get("nope")
	assert.Nil(t, err)
	assert.Nil(t, i)

	i, err = m.Get("error")
	assert.Equal(t, err.Error(), "invalid id")
	assert.Nil(t, i)

	i1, err := m.Get("valid")
	assert.Nil(t, err)
	assert.Equal(t, i1.id, "valid")

	i2, err := m.Get("valid")
	assert.Nil(t, err)
	assert.Equal(t, i2, i1)
}

// silly test, but I want to get the sharding under some type of explicit test
func Test_Get_Fuzz(t *testing.T) {
	m := NewMap[*TestItem](func(id string) (*TestItem, error) {
		return &TestItem{id: id}, nil
	})

	for i := 0; i < 200; i++ {
		id := strconv.Itoa(i)
		i, err := m.Get(id)
		assert.Nil(t, err)
		assert.Equal(t, i.id, id)
	}
}
