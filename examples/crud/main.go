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

	// create a book
	book1 := &book.Book{
		BookID:    id.NewID(),
		Name:      "The Hitchhiker's Guide to the Galaxy",
		Author:    "Douglas Adams",
		Published: time.Now().AddDate(-1, -2, -3).UTC(),
		Pages:     100,
		HardCover: true,
	}

	fmt.Printf("unsaved book: %+v\n", book1)

	// save book
	if err = oven.Save(ctx, client, book.CollectionName, book1.BookID, book1); err != nil {
		log.Fatalf("failed to save: %v", err)
	}

	// get the previously saved book
	book2, err := oven.Get[book.Book](ctx, client, book.CollectionName, book1.BookID)
	if err != nil {
		log.Fatalf("failed to get: %v", err)
	}
	fmt.Printf("saved book: %+v\n", book2)

	// check book's content
	if fmt.Sprintf("%v", book1) != fmt.Sprintf("%v", book2) {
		log.Fatalf("books are not the same, want %v, got %v", book1, book2)
	}

	// update
	updatedPages := 200
	updates := map[string]interface{}{"pages": updatedPages}
	if err = oven.Update(ctx, client, book.CollectionName, book2.BookID, updates); err != nil {
		log.Fatalf("failed to update: %v", err)
	}

	book3, err := oven.Get[book.Book](ctx, client, book.CollectionName, book2.BookID)
	if err != nil {
		log.Fatalf("failed to get: %v", err)
	}
	fmt.Printf("updated book: %+v\n", book3)

	// check book's content
	if book3.Pages != updatedPages {
		log.Fatalf("book not updated, wanted %d, got %d", updatedPages, book3.Pages)
	}

	// delete
	if err = oven.Delete(ctx, client, book.CollectionName, book3.BookID); err != nil {
		log.Fatalf("failed to delete: %v", err)
	}

	// no data found error after delete
	_, err = oven.Get[book.Book](ctx, client, book.CollectionName, book3.BookID)
	if err != nil {
		if err != oven.ErrDataNotFound {
			log.Fatalf("expected ErrDataNotFound error, got: %v", err)
		}
	}
	fmt.Print("book deleted\n")
}
