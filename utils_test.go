package utils

import (
	"bytes"
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_EncodeRequestId(t *testing.T) {
	r1 := EncodeRequestId(1, 1)
	assert.Equal(t, len(r1), 8)

	seen := make(map[string]struct{}, 1000)
	for i := uint32(0); i < 1000; i++ {
		r1 := EncodeRequestId(i, 1)
		r2 := EncodeRequestId(i, 2) // different instance id
		if _, ok := seen[r1]; ok {
			assert.Fail(t, "duplicate request ids")
		}
		if _, ok := seen[r2]; ok {
			assert.Fail(t, "duplicate request ids")
		}
		seen[r1] = struct{}{}
		seen[r2] = struct{}{}
	}
}

func Test_B2S(t *testing.T) {
	assert.Equal(t, B2S([]byte("")), "")
	assert.Equal(t, B2S([]byte("abc")), "abc")

	assert.True(t, bytes.Equal(S2B(""), []byte{}))
	assert.True(t, bytes.Equal(S2B("129"), []byte("129")))
}
