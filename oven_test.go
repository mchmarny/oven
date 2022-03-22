package oven

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/mchmarny/oven/pkg/id"
	"github.com/stretchr/testify/assert"
)

func TestOven(t *testing.T) {
	ctx := context.Background()
	c := &firestore.Client{}

	type S struct {
		ID string `firestore:"id"`
	}

	t.Run("new with no args", func(t *testing.T) {
		_, err := GetCollection(c, "")
		assert.Error(t, err)
		err = Save[S](ctx, c, "", "", nil)
		assert.Error(t, err)
		err = Update(ctx, c, "", "", nil)
		assert.Error(t, err)
		s, err := Get[S](ctx, c, "", "")
		assert.Error(t, err)
		assert.Nil(t, s)
		err = Delete(ctx, c, "", "")
		assert.Error(t, err)
	})
	t.Run("data not found error", func(t *testing.T) {
		assert.False(t, isDataNotFoundError(errors.New("data not found")))
	})
}

type TestDoc struct {
	DocID       string    `json:"id" firestore:"id"`
	StringValue string    `json:"s1" firestore:"s1"`
	TimeValue   time.Time `json:"t1" firestore:"t1"`
	Int64Value  int64     `json:"i1" firestore:"i1"`
}

func TestOvenIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	client := getTestClient(ctx, t)

	t.Run("round trip", func(t *testing.T) {
		col := fmt.Sprintf("testcol%d", time.Now().Nanosecond())
		docID := id.NewID()

		d := &TestDoc{
			DocID:       docID,
			StringValue: fmt.Sprintf("test-%d", time.Now().Nanosecond()),
			TimeValue:   time.Now().UTC(),
			Int64Value:  time.Now().UTC().Unix(),
		}

		// save
		err := Save(ctx, client, col, d.DocID, d)
		assert.NoError(t, err)

		// get
		d2, err := Get[TestDoc](ctx, client, col, docID)
		assert.NoError(t, err)
		assert.Equal(t, d.StringValue, d2.StringValue)
		assert.Equal(t, d.Int64Value, d2.Int64Value)
		// emulator seems to drop a few nanoseconds of precision
		assert.Equal(t, d.TimeValue.Format(time.RFC3339), d2.TimeValue.Format(time.RFC3339))

		// update
		updatedValue := "updated"
		m1 := map[string]interface{}{"s1": updatedValue}
		err = Update(ctx, client, col, docID, m1)
		assert.NoError(t, err)

		d3, err := Get[TestDoc](ctx, client, col, docID)
		assert.NoError(t, err)
		assert.Equal(t, updatedValue, d3.StringValue)
		assert.Equal(t, d2.TimeValue.Format(time.RFC3339), d3.TimeValue.Format(time.RFC3339))
		assert.Equal(t, d2.Int64Value, d3.Int64Value)

		// delete
		err = Delete(ctx, client, col, docID)
		assert.NoError(t, err)

		// no data found error after delete
		_, err = Get[TestDoc](ctx, client, col, docID)
		assert.Error(t, err)
		assert.Equal(t, ErrDataNotFound, err)
	})
}

func getEnv(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}

func getTestClient(ctx context.Context, t *testing.T) *firestore.Client {
	projectID := getEnv("PROJECT_ID", "")
	assert.NotEmpty(t, projectID)
	c, err := firestore.NewClient(ctx, projectID)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	return c
}
