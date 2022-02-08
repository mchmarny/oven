package oven

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/mchmarny/oven/examples/book"
	"github.com/mchmarny/oven/pkg/id"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
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

func BenchmarkSetup(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}
	projectID := getEnv("PROJECT_ID", "")
	ctx := context.Background()
	service := New(ctx, projectID)

	// save books
	for i := 0; i < benchmarkSize; i++ {
		k := &book.Book{
			BookID:    id.NewID(),
			Name:      fmt.Sprintf("Benchmark %d", i+1),
			Author:    benchmarkAuthor,
			Published: time.Now().UTC(),
			Pages:     100,
			HardCover: false,
		}
		if err := service.Save(ctx, book.CollectionName, k.BookID, k); err != nil {
			b.Fatalf("failed to save: %v", err)
		}
	}
}

func benchmarkQueryNative(b *testing.B, limit int) {
	if testing.Short() {
		b.Skip()
	}
	projectID := getEnv("PROJECT_ID", "")
	ctx := context.Background()
	service := New(ctx, projectID)

	list := make([]*book.Book, 0)
	col, err := service.GetCollection(ctx, book.CollectionName)
	if err != nil {
		b.Fatalf("error getting collection %s: %v", book.CollectionName, err)
	}

	it := col.Where("author", "==", benchmarkAuthor).OrderBy("published", firestore.Desc).Limit(limit).Documents(ctx)

	for {
		d, e := it.Next()
		if e == iterator.Done {
			break
		}
		if e != nil {
			b.Fatal(e)
		}

		item := &book.Book{}
		if e := d.DataTo(item); e != nil {
			b.Fatal(e)
		}

		list = append(list, item)
	}

	if len(list) != limit {
		b.Fatalf("expected %d, got %d", limit, len(list))
	}
}

func benchmarkQueryMeta(b *testing.B, limit int) {
	if testing.Short() {
		b.Skip()
	}
	projectID := getEnv("PROJECT_ID", "")
	ctx := context.Background()
	service := New(ctx, projectID)

	// get all books by the author
	var list []book.Book
	q := &Query{
		Collection: book.CollectionName,
		Criteria: &Criterion{
			Path:      "author", // `firestore:"author"`
			Operation: OperationTypeEqual,
			Value:     benchmarkAuthor,
		},
		OrderBy: "published", // `firestore:"published"`
		Desc:    true,
		Limit:   limit,
	}

	if err := service.Query(ctx, q, &list); err != nil {
		b.Fatalf("failed to query: %v", err)
	}

	if len(list) != limit {
		b.Fatalf("expected %d, got %d", limit, len(list))
	}
}

func BenchmarkQueryNative3(b *testing.B) { benchmarkQueryNative(b, 3) }
func BenchmarkQueryMeta3(b *testing.B)   { benchmarkQueryMeta(b, 3) }

func BenchmarkQueryNative10(b *testing.B) { benchmarkQueryNative(b, 10) }
func BenchmarkQueryMeta10(b *testing.B)   { benchmarkQueryMeta(b, 10) }
func BenchmarkQueryNative50(b *testing.B) { benchmarkQueryNative(b, 50) }
func BenchmarkQueryMeta50(b *testing.B)   { benchmarkQueryMeta(b, 50) }
