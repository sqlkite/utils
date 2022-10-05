package utils

import (
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
