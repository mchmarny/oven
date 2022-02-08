package oven

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/mchmarny/oven/pkg/id"
	"github.com/stretchr/testify/assert"
)

const (
	numOfTestQueryDocs = 10
)

func TestOvenQuery(t *testing.T) {
	t.Run("invalid destination", func(t *testing.T) {
		_, err := getDestinationMeta(nil)
		assert.Error(t, err)
	})

	t.Run("invalid destination type", func(t *testing.T) {
		var item TestDoc
		_, err := getDestinationMeta(item)
		assert.Error(t, err)
		_, err = getDestinationMeta(&item)
		assert.Error(t, err)
	})
	t.Run("struct pointer array", func(t *testing.T) {
		var list []*TestDoc
		m, err := getDestinationMeta(&list)
		assert.NoError(t, err)
		assert.NotNil(t, m)
		assert.NotNil(t, m.itemType)
		assert.NotNil(t, m.list)
		assert.True(t, m.listByPtr)
		assert.Equal(t, reflect.TypeOf(TestDoc{}), m.itemType)
		assert.Equal(t, reflect.TypeOf([]*TestDoc{}), m.list.Type())
		// this will panic but it validates the full round reflection  trip
		m.append(m.new())
	})
	t.Run("struct array", func(t *testing.T) {
		var list []TestDoc
		m, err := getDestinationMeta(&list)
		assert.NoError(t, err)
		assert.NotNil(t, m)
		assert.NotNil(t, m.itemType)
		assert.NotNil(t, m.list)
		assert.False(t, m.listByPtr)
		assert.Equal(t, reflect.TypeOf(TestDoc{}), m.itemType)
		assert.Equal(t, reflect.TypeOf([]TestDoc{}), m.list.Type())
		m.append(m.new())
	})
	t.Run("struct count", func(t *testing.T) {
		var list []TestDoc
		m, err := getDestinationMeta(&list)
		assert.NoError(t, err)
		assert.NotNil(t, m)
		m.append(m.new())
		m.append(m.new())
		m.append(m.new())
		assert.Equal(t, 3, m.list.Len())
	})
}

func TestOvenQueryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	projectID := getEnv("PROJECT_ID", "")
	assert.NotEmpty(t, projectID)

	ctx := context.Background()

	t.Run("query", func(t *testing.T) {
		s := New(ctx, projectID)
		assert.NotNil(t, s)

		col := fmt.Sprintf("testcol%d", time.Now().Nanosecond())
		val := fmt.Sprintf("test-%d", time.Now().Nanosecond())

		// save docs
		for i := 0; i < numOfTestQueryDocs; i++ {
			d := &TestDoc{
				DocID:       id.NewID(),
				StringValue: val,
				TimeValue:   time.Now().UTC(),
				Int64Value:  time.Now().UTC().Unix(),
			}
			err := s.Save(ctx, col, d.DocID, d)
			assert.NoError(t, err)
		}

		var list []*TestDoc
		q := &Query{
			Collection: col,
			Criteria: &Criterion{
				Path:      "s1",
				Operation: OperationTypeEqual,
				Value:     val,
			},
			OrderBy: "s1",
			Limit:   numOfTestQueryDocs,
		}

		err := s.Query(ctx, q, &list)
		assert.NoError(t, err)
		assert.Len(t, list, numOfTestQueryDocs)
	})
}
