package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mchmarny/oven"
	"github.com/mchmarny/oven/examples/book"
	"github.com/mchmarny/oven/pkg/id"
)

func main() {
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatal("PROJECT_ID not set")
	}

	ctx := context.Background()

	// create a new oven service instance
	service := oven.New(ctx, projectID)

	bookCollectionSize := 10
	bookAuthor := "Douglas Adams"

	// save books
	for i := 0; i < bookCollectionSize; i++ {
		b := &book.Book{
			BookID:    id.NewID(),
			Name:      fmt.Sprintf("Galaxy, volume %d", i),
			Author:    bookAuthor,
			Published: time.Now().AddDate(-1, -2, -i).UTC(),
			Pages:     100 + i,
			HardCover: false,
		}
		if err := service.Save(ctx, book.CollectionName, b.BookID, b); err != nil {
			log.Fatalf("failed to save: %v", err)
		}
	}

	// get all books by the author
	var list []book.Book
	q := &oven.Query{
		Collection: book.CollectionName,
		Criteria: &oven.Criterion{
			Path:      "author", // `firestore:"author"`
			Operation: oven.OperationTypeEqual,
			Value:     bookAuthor,
		},
		OrderBy: "published", // `firestore:"published"`
		Desc:    true,
		Limit:   bookCollectionSize,
	}

	if err := service.Query(ctx, q, &list); err != nil {
		log.Fatalf("failed to query: %v", err)
	}

	for i, b := range list {
		fmt.Printf("book[%d]: %+v\n", i, b)
	}
}
