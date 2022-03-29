package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
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

	// create firestore client
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	bookCollectionSize := 10
	bookAuthor := "Douglas Adams"

	// save books
	for i := 0; i < bookCollectionSize; i++ {
		b := &book.Book{
			BookID:    id.NewID(),
			Name:      fmt.Sprintf("Galaxy, volume %d", i),
			Author:    bookAuthor,
			Published: time.Now().AddDate(-1, -2, -i).UTC(),
			Pages:     book.SampleBookNumOfPages,
			HardCover: false,
		}
		if err = oven.Save(ctx, client, book.CollectionName, b.BookID, b); err != nil {
			log.Fatalf("failed to save: %v", err)
		}
	}

	// get all books by the author
	q := &oven.Criteria{
		Collection: book.CollectionName,
		Criterions: []*oven.Criterion{
			{
				Path:      "author", // `firestore:"author"`
				Operation: oven.OperationTypeEqual,
				Value:     bookAuthor,
			},
		},
		OrderBy: "published", // `firestore:"published"`
		Desc:    true,
		Limit:   bookCollectionSize,
	}

	list, err := oven.Query[book.Book](ctx, client, q)
	if err != nil {
		log.Fatalf("failed to query: %v", err)
	}

	for i, b := range list {
		fmt.Printf("book[%d]: %+v\n", i, b)
	}
}
