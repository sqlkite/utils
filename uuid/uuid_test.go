package uuid

import (
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_String(t *testing.T) {
	uuid1 := String()
	uuid2 := String()
	assert.Equal(t, len(uuid1), 36)
	assert.Equal(t, len(uuid2), 36)
	assert.True(t, uuid1 != uuid2)
}

func Test_FromBytes(t *testing.T) {
	s, err := FromBytes([]byte{204, 193, 82, 169, 150, 64, 52, 71, 92, 228, 173, 248, 223, 220, 70, 252})
	assert.Nil(t, err)
	assert.Equal(t, s, "ccc152a9-9640-3447-5ce4-adf8dfdc46fc")

	s, err = FromBytes([]byte{0})
	assert.NotNil(t, err)
	assert.Equal(t, s, "")
}
