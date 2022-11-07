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

func Test_Paging(t *testing.T) {
	perpage, offset := Paging(0, 0, 10)
	assert.Equal(t, offset, 0)
	assert.Equal(t, perpage, 10)

	perpage, offset = Paging(-1, -1, 11)
	assert.Equal(t, offset, 0)
	assert.Equal(t, perpage, 11)

	perpage, offset = Paging(1, 0, 10)
	assert.Equal(t, offset, 0)
	assert.Equal(t, perpage, 1)

	perpage, offset = Paging(10, 1, 20)
	assert.Equal(t, offset, 0)
	assert.Equal(t, perpage, 10)

	perpage, offset = Paging(10, 2, 20)
	assert.Equal(t, offset, 10)
	assert.Equal(t, perpage, 10)
}
