package id

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	t.Run("new ID has prefix", func(t *testing.T) {
		s := NewID()
		assert.NotEmpty(t, s)
		assert.True(t, strings.HasPrefix(s, recordIDPrefix))
	})
	t.Run("new ID is of consistent length", func(t *testing.T) {
		want := len(NewID())
		for i := 0; i < 10; i++ {
			assert.Equal(t, want, len(NewID()))
		}
	})
	t.Run("to ID is of consistent length", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := fmt.Sprintf("id-%d-%d", i, time.Now().Nanosecond())
			id1 := ToID(s)
			id2 := ToID(s)
			assert.Equal(t, id1, id2)
		}
	})
}
