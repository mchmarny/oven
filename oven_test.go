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
	t.Run("new with no args", func(t *testing.T) {
		ctx := context.Background()
		s := New(ctx, "test")
		assert.NotNil(t, s)
		assert.NoError(t, s.Close())
		_, err := s.GetCollection(ctx, "")
		assert.Error(t, err)
		err = s.Save(ctx, "", "", nil)
		assert.Error(t, err)
		err = s.Update(ctx, "", "", nil)
		assert.Error(t, err)
		err = s.Get(ctx, "", "", nil)
		assert.Error(t, err)
		err = s.Delete(ctx, "", "")
		assert.Error(t, err)
	})
	t.Run("new with firestore client", func(t *testing.T) {
		c := &firestore.Client{}
		s := NewWithClient(c)
		assert.NotNil(t, s)
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

	projectID := getEnv("PROJECT_ID", "")
	assert.NotEmpty(t, projectID)

	ctx := context.Background()

	t.Run("new with valid args", func(t *testing.T) {
		s := New(ctx, projectID)
		assert.NotNil(t, s)
	})

	t.Run("round trip", func(t *testing.T) {
		s := New(ctx, projectID)
		assert.NotNil(t, s)

		col := fmt.Sprintf("testcol%d", time.Now().Nanosecond())
		docID := id.NewID()

		d := &TestDoc{
			DocID:       docID,
			StringValue: fmt.Sprintf("test-%d", time.Now().Nanosecond()),
			TimeValue:   time.Now().UTC(),
			Int64Value:  time.Now().UTC().Unix(),
		}

		// save
		err := s.Save(ctx, col, d.DocID, d)
		assert.NoError(t, err)

		// get
		d2 := &TestDoc{}
		err = s.Get(ctx, col, docID, d2)
		assert.NoError(t, err)
		assert.Equal(t, d.StringValue, d2.StringValue)
		assert.Equal(t, d.Int64Value, d2.Int64Value)
		// emulator seems to drop a few nanoseconds of precision
		assert.Equal(t, d.TimeValue.Format(time.RFC3339), d2.TimeValue.Format(time.RFC3339))

		// update
		updatedValue := "updated"
		m1 := map[string]interface{}{"s1": updatedValue}
		err = s.Update(ctx, col, docID, m1)
		assert.NoError(t, err)

		d3 := &TestDoc{}
		err = s.Get(ctx, col, docID, d3)
		assert.NoError(t, err)
		assert.Equal(t, updatedValue, d3.StringValue)
		assert.Equal(t, d2.TimeValue.Format(time.RFC3339), d3.TimeValue.Format(time.RFC3339))
		assert.Equal(t, d2.Int64Value, d3.Int64Value)

		// delete
		err = s.Delete(ctx, col, docID)
		assert.NoError(t, err)

		// no data found error after delete
		d4 := &TestDoc{}
		err = s.Get(ctx, col, docID, d4)
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
