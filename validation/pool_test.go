package validation

import (
	"testing"

	"src.sqlkite.com/tests/assert"
)

func Test_Pool_Depleted(t *testing.T) {
	p := NewPool(2, 1)
	assert.Equal(t, p.Len(), 2)
	assert.Equal(t, p.Depleted(), 0)

	l1 := p.Checkout()
	assert.Equal(t, p.Len(), 1)
	assert.Equal(t, p.Depleted(), 0)

	l2 := p.Checkout()
	assert.Equal(t, p.Len(), 0)
	assert.Equal(t, p.Depleted(), 0)

	l3 := p.Checkout()
	assert.Equal(t, p.Len(), 0)
	assert.Equal(t, p.Depleted(), 1)
	assert.Equal(t, p.Depleted(), 0) // calling Delpeted resets it

	assert.NotEqual(t, l1, l2)
	assert.NotEqual(t, l1, l3)
	assert.NotEqual(t, l2, l3)
}

func Test_Pool_DynamicCreationWontReleaseToPool(t *testing.T) {
	p := NewPool(1, 3)

	l1 := p.Checkout()
	l2 := p.Checkout()
	assert.NotEqual(t, l1, l2)

	l1.Release()
	l2.Release()

	assert.Equal(t, p.Len(), 1)
}
