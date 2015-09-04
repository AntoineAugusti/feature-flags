package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint32ToBytes(t *testing.T) {
	assert.Equal(t, Uint32ToBytes(uint32(1)), []byte{0x31})
	assert.Equal(t, Uint32ToBytes(uint32(2)), []byte{0x32})
	assert.Equal(t, Uint32ToBytes(uint32(10)), []byte{0x31, 0x30})
}

func TestIntInSlice(t *testing.T) {
	assert.True(t, IntInSlice(1, []uint32{1, 2, 3}))
	assert.True(t, IntInSlice(2, []uint32{1, 2, 3}))
	assert.True(t, IntInSlice(3, []uint32{1, 2, 3}))

	assert.False(t, IntInSlice(4, []uint32{1, 2, 3}))
}

func TestStringInSlice(t *testing.T) {
	assert.True(t, StringInSlice("foo", []string{"foo", "bar"}))
	assert.True(t, StringInSlice("bar", []string{"foo", "bar"}))

	assert.False(t, StringInSlice("baz", []string{"foo", "bar"}))
}
