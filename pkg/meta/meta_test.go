package meta

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestDoc struct {
	DocID       string    `json:"id" firestore:"id"`
	StringValue string    `json:"s1" firestore:"s1"`
	TimeValue   time.Time `json:"t1" firestore:"t1"`
	Int64Value  int64     `json:"i1" firestore:"i1"`
}

func TestMeta(t *testing.T) {
	t.Run("invalid destination", func(t *testing.T) {
		_, err := New(nil)
		assert.Error(t, err)
	})

	t.Run("invalid destination type", func(t *testing.T) {
		var item TestDoc
		_, err := New(item)
		assert.Error(t, err)
		_, err = New(&item)
		assert.Error(t, err)
	})
	t.Run("struct pointer array", func(t *testing.T) {
		var list []*TestDoc
		m, err := New(&list)
		assert.NoError(t, err)
		assert.NotNil(t, m)
		assert.NotNil(t, m.itemType)
		assert.NotNil(t, m.list)
		assert.True(t, m.listByPtr)
		assert.Equal(t, reflect.TypeOf(TestDoc{}), m.itemType)
		assert.Equal(t, reflect.TypeOf([]*TestDoc{}), m.list.Type())
		// this will panic but it validates the full round reflection  trip
		m.Append(m.Item())
	})
	t.Run("struct array", func(t *testing.T) {
		var list []TestDoc
		m, err := New(&list)
		assert.NoError(t, err)
		assert.NotNil(t, m)
		assert.NotNil(t, m.itemType)
		assert.NotNil(t, m.list)
		assert.False(t, m.listByPtr)
		assert.Equal(t, reflect.TypeOf(TestDoc{}), m.itemType)
		assert.Equal(t, reflect.TypeOf([]TestDoc{}), m.list.Type())
		m.Append(m.Item())
	})
	t.Run("struct count", func(t *testing.T) {
		var list []TestDoc
		m, err := New(&list)
		assert.NoError(t, err)
		assert.NotNil(t, m)
		m.Append(m.Item())
		m.Append(m.Item())
		m.Append(m.Item())
		assert.Equal(t, 3, m.list.Len())
	})
}

func BenchmarkMeta10(b *testing.B)    { benchmarkMeta(b, 10) }
func BenchmarkMeta100(b *testing.B)   { benchmarkMeta(b, 100) }
func BenchmarkMeta1000(b *testing.B)  { benchmarkMeta(b, 1000) }
func BenchmarkMeta10000(b *testing.B) { benchmarkMeta(b, 10000) }

func benchmarkMeta(b *testing.B, limit int) {
	if testing.Short() {
		b.Skip()
	}
	var list []TestDoc
	m, err := New(&list)
	if err != nil {
		b.Fatalf("failed to create meta: %v", err)
	}

	for i := 0; i < limit; i++ {
		m.Append(m.Item())
	}

	if m.list.Len() != limit {
		b.Fatalf("expected %d items, got %d", limit, m.list.Len())
	}
}