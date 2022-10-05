package log

import (
	"math"
	"strconv"
	"testing"

	"src.goblgobl.com/tests/assert"
)

func Test_Field_Int(t *testing.T) {
	f := NewField().Int("over", 9000).Finalize()
	assert.Equal(t, f.fields["over"].(int), 9000)
	assert.Equal(t, string(f.kv), "over=9000")

	f = NewField().Int("o", math.MaxInt).Finalize()
	assert.Equal(t, len(f.fields), 1)
	assert.Equal(t, f.fields["o"].(int), math.MaxInt)
	assert.Equal(t, string(f.kv), "o="+strconv.Itoa(math.MaxInt))

	f = NewField().Int("o", math.MinInt).Finalize()
	assert.Equal(t, len(f.fields), 1)
	assert.Equal(t, f.fields["o"].(int), math.MinInt)
	assert.Equal(t, string(f.kv), "o="+strconv.Itoa(math.MinInt))
}

func Test_Field_String(t *testing.T) {
	f := NewField().String("leto", "atreides").Finalize()
	assert.Equal(t, len(f.fields), 1)
	assert.Equal(t, f.fields["leto"].(string), "atreides")
	assert.Equal(t, string(f.kv), "leto=atreides")

	f = NewField().String("name", "ghanima atreides").Finalize()
	assert.Equal(t, len(f.fields), 1)
	assert.Equal(t, f.fields["name"].(string), "ghanima atreides")
	assert.Equal(t, string(f.kv), "name=\"ghanima atreides\"")
}

func Test_Field_Multiple(t *testing.T) {
	f := NewField().
		String("leto", "atreides II").
		String("type", "worm").
		Int("age", 3000).
		Finalize()
	assert.Equal(t, len(f.fields), 3)
	assert.Equal(t, f.fields["leto"].(string), "atreides II")
	assert.Equal(t, f.fields["type"].(string), "worm")
	assert.Equal(t, f.fields["age"].(int), 3000)

	kv := KvParse(string(f.kv))
	assert.Equal(t, len(kv), 3)
	assert.Equal(t, kv["leto"], `"atreides II"`)
	assert.Equal(t, kv["type"], "worm")
	assert.Equal(t, kv["age"], "3000")
}
