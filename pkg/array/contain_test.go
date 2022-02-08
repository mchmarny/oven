package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayContains(t *testing.T) {
	list := []int64{1, 2, 3, 4, 5}
	assert.False(t, Contains(nil, int64(1)))
	assert.True(t, Contains(list, int64(1)))
	assert.False(t, Contains(list, int64(0)))
}
