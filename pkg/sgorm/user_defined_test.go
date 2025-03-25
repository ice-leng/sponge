package sgorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	v1 := false
	assert.Equal(t, NewBool(v1).String(), "false")
	assert.Equal(t, NewBool(&v1).String(), "false")

	v2 := true
	assert.Equal(t, NewBool(v2).String(), "true")
	assert.Equal(t, NewBool(&v2).String(), "true")

	assert.Nil(t, NewBool(nil))
	assert.Equal(t, NewBool(nil).String(), "false")

	assert.NoError(t, NewBool(v2).Scan(nil))
	assert.NoError(t, NewBool(v2).Scan([]byte{0}))
	assert.NoError(t, NewBool(v2).Scan([]byte{1}))
	assert.NoError(t, NewBool(v2).Scan("true"))
	assert.NoError(t, NewBool(v2).Scan("1"))
	assert.NoError(t, NewBool(v2).Scan("t"))
	assert.Error(t, NewBool(v2).Scan(3.14))

	_, err := NewBool(v1).Value()
	assert.NoError(t, err)
	_, err = NewBool(v2).Value()
	assert.NoError(t, err)

	SetDriver("postgres")
	value, err := NewBool(v2).Value()
	assert.NoError(t, err)
	assert.Equal(t, true, value)
}
