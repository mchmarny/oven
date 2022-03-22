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
)

func TestOvenQueryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	client := getTestClient(ctx, t)

	t.Run("query", func(t *testing.T) {
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
			err := Save(ctx, client, col, d.DocID, d)
			assert.NoError(t, err)
		}

		var list []*TestDoc
		q := &Criteria{
			Collection: col,
			Criterions: []*Criterion{
				{
					Path:      "s1",
					Operation: OperationTypeEqual,
					Value:     val,
				},
			},
			OrderBy: "s1",
			Limit:   numOfTestQueryDocs,
		}

		list, err := Query[TestDoc](ctx, client, q)
		assert.NoError(t, err)
		assert.Len(t, list, numOfTestQueryDocs)
	})
}
