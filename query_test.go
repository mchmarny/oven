package oven

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mchmarny/oven/pkg/id"
	"github.com/stretchr/testify/assert"
)

const (
	numOfTestQueryDocs = 10
	benchmarkSize      = 100
	benchmarkAuthor    = "Douglas Adams"
)

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
