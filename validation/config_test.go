package validation

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Configure_Defaults(t *testing.T) {
	err := Configure(Config{})
	assert.Nil(t, err)
	assert.Equal(t, len(globalPool.list), 100)

	r := globalPool.Checkout()
	defer r.Release()
	assert.Equal(t, len(r.errors), 20)
	assert.Equal(t, r.pool, globalPool)
}

func Test_Configure_Custom(t *testing.T) {
	err := Configure(Config{
		PoolSize:  3,
		MaxErrors: 7,
	})

	assert.Nil(t, err)
	assert.Equal(t, len(globalPool.list), 3)

	r := Checkout()
	defer r.Release()
	assert.Equal(t, len(r.errors), 7)
}
